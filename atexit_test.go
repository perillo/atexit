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
		t.Errorf("the resource was not released: got %d", res.n)
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
		t.Errorf("the resource was not released: got %d", res.n)
	}
}
