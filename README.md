# atexit [![GoDoc](https://godoc.org/github.com/perillo/atexit?status.svg)](http://godoc.org/github.com/perillo/atexit)

Package `atexit` implements support for running deferred functions in case of
*abnormal exit* (in contrast to a *normal exit* when the program returns from
the main function).

Since calling `os.Exit`, e.g. during a signal handler, does not call deferred
functions, a complementary mechanism is required when the program acquires
resources that are not automatically released by the operating system at
program termination, e.g. `SYSV` shared memory.

`atexit` is designed to work with, and complement, *Go* standard deferred
mechanism.  The `Exit` function provided by this package must be used, in order
to run registered deferred functions.

The `Exit` function **SHOULD** only be called in case of abnormal program
termination.

The idiomatic use of atexit is:
```go
    AcquireResource(...)
    // ...
    defer atexit.Do(func() {
        ReleaseResource(...)
    })()
```

`ReleaseResource` will be called by the runtime deferred mechanism in case of
normal program termination, but will be called by `atexit` when `atexit.Exit`
is called during an abnormal program termination (e.g. after having received a
signal), but only if they have not been already be called by the *Go* runtime.
