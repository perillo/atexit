// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package atexit

import (
	"sync"
	"testing"
	"time"
)

// resource represents a simple resource that can be acquired and released.  It
// is implemented using an integer so a test case can easily check if it was
// correctly released
type resource struct {
	sync.Mutex // protects n
	n          int
}

func (r *resource) acquire() {
	r.Lock()
	r.n++
	r.Unlock()
}

func (r *resource) release() {
	r.Lock()
	r.n--
	r.Unlock()
}

func (r *resource) value() int {
	r.Lock()
	n := r.n
	r.Unlock()

	return n
}

// TestExit tests that an atexit function is called after a normal function
// termination.
func TestExit(t *testing.T) {
	var (
		res resource
		wg  sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		res.acquire()
		defer Do(func() {
			res.release()
		})()
	}()

	wg.Wait()
	if res.value() != 0 {
		t.Error("the resource was not released")
	}
}

// TestAbort tests that an atexit function is called after an abnormal function
// termination.
func TestAbort(t *testing.T) {
	var (
		res resource
		wg  sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		res.acquire()
		defer Do(func() {
			res.release()
		})()

		// Sleep to allow an asynchronous abortion.
		time.Sleep(time.Second)
	}()

	// Simulate an abort.
	time.AfterFunc(500*time.Millisecond, func() {
		exit()
	})

	wg.Wait()
	if res.value() != 0 {
		t.Error("the resource was not released")
	}
}

// operation implements two operations that are not commutative.  The value
// method will return different results depending on whether op1 is called
// first followed by op2, or op2 is called first followed by op1.
type operation struct {
	n int
}

func (op *operation) op1() {
	op.n--
}

func (op *operation) op2() {
	op.n *= 10
}

func (op *operation) value() int {
	return op.n
}

// TestExitOrder tests that atexit functions are called in the correct order
// after a normal function termination.
func TestExitOrder(t *testing.T) {
	var (
		op operation
		wg sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		defer Do(func() {
			op.op1()
		})()
		defer Do(func() {
			op.op2()
		})()
	}()

	wg.Wait()
	if op.value() != (0*10)-1 {
		t.Error("the atexit functions where not called in the correct order")
	}
}

// TestAbortOrder tests that atexit functions are called in the correct order
// after an abnormal function termination.
func TestAbortOrder(t *testing.T) {
	var (
		op operation
		wg sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Change the order of operations compared to TestExitOrder, to make
		// sure the code is actually correct.
		defer Do(func() {
			op.op2()
		})()
		defer Do(func() {
			op.op1()
		})()

		// Sleep to allow an asynchronous abortion.
		time.Sleep(time.Second)
	}()

	// Simulate an abort.
	time.AfterFunc(500*time.Millisecond, func() {
		exit()
	})

	wg.Wait()
	if op.value() != (0-1)*10 {
		t.Error("the atexit functions where not called in the correct order")
	}
}
