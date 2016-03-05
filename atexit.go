// Copyright 2015 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Package atexit implements support for running deferred functions in case of
// abnormal exit (in contrast to a normal exit when the program returns from
// the main function).
//
// Since calling os.Exit, e.g. during a signal handler, does not call deferred
// functions, a complementary mechanism is required when the program acquires
// resources that are not automatically released by the operating system at
// program termination, e.g. SYSV shared memory.
//
// atexit is designed to work with, and complement, Go standard deferred
// mechanism.
// The Exit function provided by this package must be used, in order to run
// registered deferred functions.
//
// The Exit function SHOULD only be called in case of abnormal program
// termination.
package atexit

import (
	"os"
	"sync"
)

type deferred struct {
	f    func()
	m    sync.Mutex
	done uint32
}

// runAtDefer calls the deferred function during standard deferred mechanism.
func (d *deferred) runAtDefer() {
	d.m.Lock()
	defer d.m.Unlock()

	if d.done == 0 {
		defer func() {
			d.done = 1
		}()

		d.f()
	}
}

// runAtExit calls the deferred function during an abnormal program exit.
func (d *deferred) runAtExit() {
	// No need to unlock mutex and update done field.
	d.m.Lock()

	if d.done == 0 {
		d.f()
	}
}

// Do registers the function f to be called in case of abnormal program
// termination, caused by calling the Exit function, and returns a function
// that can be used in a defer statement.
//
// Do SHOULD only be used to release resources that are not automatically
// released by the operating system at program termination.
// The idiomatic use of Do is
//  AcquireResource(...)
//  ...
//  defer atexit.Do(func() {
//      ReleaseResource(...)
//  })()
//
// The function f is called exactly once, either during standard deferred
// mechanism or as part of atexit termination.
//
// If f panics, atexit considers it to have returned; it will not be called again.
func Do(f func()) func() {
	d := &deferred{f: f}
	// TODO(mperillo): protect with a Mutex.
	dl = append(dl, d)

	return d.runAtDefer
}

// Registered deferred functions.
var dl = make([]*deferred, 0, 10)

// Exit runs all registered deferred functions, in FIFO order, and then causes
// the current program to exit with the given status code.
// In case one of the deferred functions panics, the exit status is ignored and
// control passes to Go runtime.
func Exit(code int) {
	func() {
		for _, d := range dl {
			defer d.runAtExit()
		}
	}()

	os.Exit(code)
}
