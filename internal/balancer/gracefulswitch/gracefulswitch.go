package gracefulswitch

import (
	"errors"
	"fmt"
	"sync"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/balancer"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/connectivity"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/balancer/base"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/resolver"
)

var errBalancerClosed = errors.New("gracefulSwitchBalancer is closed")

// NewBalancer returns a graceful switch Balancer.
func NewBalancer(cc balancer.ClientConn, opts balancer.BuildOptions) *Balancer {
	return &Balancer{
		cc:    cc,
		bOpts: opts,
	}
}

// Balancer is a utility to gracefully switch from one balancer to
// a new balancer. It implements the gracefulswitch.Balancer interface.
type Balancer struct {
	bOpts balancer.BuildOptions
	cc    balancer.ClientConn

	// mu protects the following fields and all fields within balancerCurrent
	// and balancerPending. mu does not need to be held when calling into the
	// child balancers, as all calls into these children happen only as a direct
	// result of a call into the gracefulSwitchBalancer, which are also
	// guaranteed to be synchronous. There is one exception: an UpdateState call
	// from a child balancer when current and pending are populated can lead to
	// calling Close() on the current. To prevent that racing with an
	// UpdateSubConnState from the channel, we hold currentMu during Close and
	// UpdateSubConnState calls.
	mu              sync.Mutex
	balancerCurrent *balancerWrapper
	balancerPending *balancerWrapper
	closed          bool // set to true when this balancer is closed

	// currentMu must be locked before mu. This mutex guards against this
	// sequence of events: UpdateSubConnState() called, finds the
	// balancerCurrent, gives up lock, updateState comes in, causes Close() on
	// balancerCurrent before the UpdateSubConnState is called on the
	// balancerCurrent.
	currentMu sync.Mutex
}

// swap swaps out the current lb with the pending lb and updates the ClientConn.
// The caller must hold gsb.mu.
func (gsb *Balancer) swap() {
	gsb.cc.UpdateState(gsb.balancerPending.lastState)
	cur := gsb.balancerCurrent
	gsb.balancerCurrent = gsb.balancerPending
	gsb.balancerPending = nil
	go func() {
		gsb.currentMu.Lock()
		defer gsb.currentMu.Unlock()
		cur.Close()
	}()
}

// Helper function that checks if the balancer passed in is current or pending.
// The caller must hold gsb.mu.
func (gsb *Balancer) balancerCurrentOrPending(bw *balancerWrapper) bool {
	return bw == gsb.balancerCurrent || bw == gsb.balancerPending
}

// SwitchTo initializes the graceful switch process, which completes based on
// connectivity state changes on the current/pending balancer. Thus, the switch
// process is not complete when this method returns. This method must be called
// synchronously alongside the rest of the balancer.Balancer methods this
// Graceful Switch Balancer implements.
func (gsb *Balancer) SwitchTo(builder balancer.Builder) error {
	gsb.mu.Lock()
	if gsb.closed {
		gsb.mu.Unlock()
		return errBalancerClosed
	}
	bw := &balancerWrapper{
		gsb: gsb,
		lastState: balancer.State{
			ConnectivityState: connectivity.Connecting,
			Picker:            base.NewErrPicker(balancer.ErrNoSubConnAvailable),
		},
		subconns: make(map[balancer.SubConn]bool),
	}
	balToClose := gsb.balancerPending // nil if there is no pending balancer
	if gsb.balancerCurrent == nil {
		gsb.balancerCurrent = bw
	} else {
		gsb.balancerPending = bw
	}
	gsb.mu.Unlock()
	balToClose.Close()
	// This function takes a builder instead of a balancer because builder.Build
	// can call back inline, and this utility needs to handle the callbacks.
	newBalancer := builder.Build(bw, gsb.bOpts)
	if newBalancer == nil {
		// This is illegal and should never happen; we clear the balancerWrapper
		// we were constructing if it happens to avoid a potential panic.
		gsb.mu.Lock()
		if gsb.balancerPending != nil {
			gsb.balancerPending = nil
		} else {
			gsb.balancerCurrent = nil
		}
		gsb.mu.Unlock()
		return balancer.ErrBadResolverState
	}

	// This write doesn't need to take gsb.mu because this field never gets read
	// or written to on any calls from the current or pending. Calls from WeCross
	// to this balancer are guaranteed to be called synchronously, so this
	// bw.Balancer field will never be forwarded to until this SwitchTo()
	// function returns.
	bw.Balancer = newBalancer
	return nil
}

// balancerWrapper wraps a balancer.Balancer, and overrides some Balancer
// methods to help cleanup SubConns created by the wrapped balancer.
//
// It implements the balancer.ClientConn interface and is passed down in that
// capacity to the wrapped balancer. It maintains a set of subConns created by
// the wrapped balancer and calls from the latter to create/update/remove
// SubConns update this set before being forwarded to the parent ClientConn.
// State updates from the wrapped balancer can result in invocation of the
// graceful switch logic.
type balancerWrapper struct {
	balancer.Balancer
	gsb *Balancer

	lastState balancer.State
	subconns  map[balancer.SubConn]bool // subconns created by this balancer
}

func (bw *balancerWrapper) NewSubConn(addrs []resolver.Address) (balancer.SubConn, error) {
	bw.gsb.mu.Lock()
	if !bw.gsb.balancerCurrentOrPending(bw) {
		bw.gsb.mu.Unlock()
		return nil, fmt.Errorf("%T at address %p that called NewSubConn is deleted", bw, bw)
	}
	bw.gsb.mu.Unlock()

	sc, err := bw.gsb.cc.NewSubConn(addrs)
	if err != nil {
		return nil, err
	}
	bw.gsb.mu.Lock()
	if !bw.gsb.balancerCurrentOrPending(bw) { // balancer was closed during this call
		bw.gsb.cc.RemoveSubConn(sc)
		bw.gsb.mu.Unlock()
		return nil, fmt.Errorf("%T at address %p that called NewSubConn is deleted", bw, bw)
	}
	bw.subconns[sc] = true
	bw.gsb.mu.Unlock()
	return sc, nil
}

func (bw *balancerWrapper) UpdateState(state balancer.State) {
	// Hold the mutex for this entire call to ensure it cannot occur
	// concurrently with other updateState() calls. This causes updates to
	// lastState and calls to cc.UpdateState to happen atomically.
	bw.gsb.mu.Lock()
	defer bw.gsb.mu.Unlock()
	bw.lastState = state

	if !bw.gsb.balancerCurrentOrPending(bw) {
		return
	}

	if bw == bw.gsb.balancerCurrent {
		// In the case that the current balancer exits READY, and there is a pending
		// balancer, you can forward the pending balancer's cached State up to
		// ClientConn and swap the pending into the current. This is because there
		// is no reason to gracefully switch from and keep using the old policy as
		// the ClientConn is not connected to any backends.
		if state.ConnectivityState != connectivity.Ready && bw.gsb.balancerPending != nil {
			bw.gsb.swap()
			return
		}
		// Even if there is a pending balancer waiting to be gracefully switched to,
		// continue to forward current balancer updates to the Client Conn. Ignoring
		// state + picker from the current would cause undefined behavior/cause the
		// system to behave incorrectly from the current LB policies perspective.
		// Also, the current LB is still being used by WeCross to choose SubConns per
		// RPC, and thus should use the most updated form of the current balancer.
		bw.gsb.cc.UpdateState(state)
		return
	}
	// This method is now dealing with a state update from the pending balancer.
	// If the current balancer is currently in a state other than READY, the new
	// policy can be swapped into place immediately. This is because there is no
	// reason to gracefully switch from and keep using the old policy as the
	// ClientConn is not connected to any backends.
	if state.ConnectivityState != connectivity.Connecting || bw.gsb.balancerCurrent.lastState.ConnectivityState != connectivity.Ready {
		bw.gsb.swap()
	}
}

func (bw *balancerWrapper) RemoveSubConn(sc balancer.SubConn) {
	bw.gsb.mu.Lock()
	if !bw.gsb.balancerCurrentOrPending(bw) {
		bw.gsb.mu.Unlock()
		return
	}
	bw.gsb.mu.Unlock()
	bw.gsb.cc.RemoveSubConn(sc)
}

func (bw *balancerWrapper) UpdateAddresses(sc balancer.SubConn, addrs []resolver.Address) {
	bw.gsb.mu.Lock()
	if !bw.gsb.balancerCurrentOrPending(bw) {
		bw.gsb.mu.Unlock()
		return
	}
	bw.gsb.mu.Unlock()
	bw.gsb.cc.UpdateAddresses(sc, addrs)
}

// Returns nil if the graceful switch balancer is closed.
func (gsb *Balancer) latestBalancer() *balancerWrapper {
	gsb.mu.Lock()
	defer gsb.mu.Unlock()
	if gsb.balancerPending != nil {
		return gsb.balancerPending
	}
	return gsb.balancerCurrent
}

// UpdateClientConnState forwards the update to the latest balancer created.
func (gsb *Balancer) UpdateClientConnState(state balancer.ClientConnState) error {
	// The resolver data is only relevant to the most recent LB Policy.
	balToUpdate := gsb.latestBalancer()
	if balToUpdate == nil {
		return errBalancerClosed
	}
	// Perform this call without gsb.mu to prevent deadlocks if the child calls
	// back into the channel. The latest balancer can never be closed during a
	// call from the channel, even without gsb.mu held.
	return balToUpdate.UpdateClientConnState(state)
}

// UpdateSubConnState forwards the update to the appropriate child.
func (gsb *Balancer) UpdateSubConnState(sc balancer.SubConn, state balancer.SubConnState) {
	gsb.currentMu.Lock()
	defer gsb.currentMu.Unlock()
	gsb.mu.Lock()
	// Forward update to the appropriate child. Even if there is a pending
	// balancer, the current balancer should continue to get SubConn updates to
	// maintain the proper state while the pending is still connecting.
	var balToUpdate *balancerWrapper
	if gsb.balancerCurrent != nil && gsb.balancerPending.subconns[sc] {
		balToUpdate = gsb.balancerCurrent
	} else if gsb.balancerPending != nil && gsb.balancerPending.subconns[sc] {
		balToUpdate = gsb.balancerPending
	}
	gsb.mu.Unlock()
	if balToUpdate == nil {
		// SubConn belonged to a stale lb policy that has not yet fully closed,
		// or the balancer was already closed.
		return
	}
	balToUpdate.UpdateSubConnState(sc, state)
}

// ResolverError forwards the error to the latest balancer created.
func (gsb *Balancer) ResolverError(err error) {
	// The resolver data is only relevant to the most recent LB Policy.
	balToUpdate := gsb.latestBalancer()
	if balToUpdate == nil {
		return
	}
	// Perform this call without gsb.mu to prevent deadlocks if the child calls
	// back into the channel. The latest balancer can never be closed during a
	// call from the channel, even without gsb.mu held.
	balToUpdate.ResolverError(err)
}

// ExitIdle forwards the call to the latest balancer created.
//
// If the latest balancer does not support ExitIdle, the subConns are
// re-connected to manually.
func (gsb *Balancer) ExitIdle() {
	balToUpdate := gsb.latestBalancer()
	if balToUpdate == nil {
		return
	}
	// There is no need to protect this read with a mutex, as the write to the
	// Balancer field happens in SwitchTo, which completes before this can be
	// called.
	if ei, ok := balToUpdate.Balancer.(balancer.ExitIdler); ok {
		ei.ExitIdle()
		return
	}
	for sc := range balToUpdate.subconns {
		sc.Connect()
	}
}

// Close closes any active child balancers.
func (gsb *Balancer) Close() {
	gsb.mu.Lock()
	gsb.closed = true
	currentBalancerToClose := gsb.balancerCurrent
	gsb.balancerCurrent = nil
	pendingBalancerToClose := gsb.balancerPending
	gsb.balancerPending = nil
	gsb.mu.Unlock()

	currentBalancerToClose.Close()
	pendingBalancerToClose.Close()
}
