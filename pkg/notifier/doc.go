// Package notifier is an interface that must be implemented to be used by the monitor.
//
// A notifier will take a notification message and then use it, in some way, to convey
// differences from the previously known good state and the current state.
package notifier
