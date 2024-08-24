package locale

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/snivilised/li18ngo"
)

// 🔊 snippets help:
// - to invoke the snippet, type the prefix, eg p18e (plain i18n error)
// then type the values of the placeholders, with a <TAB> to go to the next
// when finished, hit <ENTER>
// - to update snippets, hit <SHIFT>-<CMD>-P, Configure User Snippets

// ❌ ErrLackPoolFunc

// LackPoolFuncErrorTemplData will be returned when invokers don't provide function for pool.
type LackPoolFuncErrorTemplData struct {
	pantsTemplData
}

// Message
func (td LackPoolFuncErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "lack-pool-func.error",
		Description: "error when the function for the PoolFunc is missing and needs to be defined",
		Other:       "must provide function for pool func",
	}
}

type LackPoolFuncError struct {
	li18ngo.LocalisableError
}

var ErrLackPoolFunc = LackPoolFuncError{
	LocalisableError: li18ngo.LocalisableError{
		Data: LackPoolFuncErrorTemplData{},
	},
}

// ❌ ErrInvalidPoolExpiry

// InvalidPoolExpiryErrorTemplData will be returned when setting a negative number as the
// periodic duration to purge goroutines.
type InvalidPoolExpiryErrorTemplData struct {
	pantsTemplData
}

// Message
func (td InvalidPoolExpiryErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-pool-expiry.error",
		Description: "error when negative number set as the periodic duration to purge goroutines",
		Other:       "invalid expiry for pool",
	}
}

type InvalidPoolExpiryError struct {
	li18ngo.LocalisableError
}

var ErrInvalidPoolExpiry = InvalidPoolExpiryError{
	LocalisableError: li18ngo.LocalisableError{
		Data: InvalidPoolExpiryErrorTemplData{},
	},
}

// ❌ ErrPoolClosed

// PoolClosedErrorTemplData will be returned when submitting task to a closed pool.
type PoolClosedErrorTemplData struct {
	pantsTemplData
}

// Message
func (td PoolClosedErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "pool-closed.error",
		Description: "error created when submitting task to a closed pool.",
		Other:       "this pool has been closed",
	}
}

type PoolClosedError struct {
	li18ngo.LocalisableError
}

var ErrPoolClosed = PoolClosedError{
	LocalisableError: li18ngo.LocalisableError{
		Data: PoolClosedErrorTemplData{},
	},
}

// ❌ ErrPoolOverload

// PoolOverloadErrorTemplData will be returned when the pool is full and no
// workers available.
type PoolOverloadErrorTemplData struct {
	pantsTemplData
}

// Message
func (td PoolOverloadErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "pool-overload.error",
		Description: "will be returned when the pool is full and no workers available.",
		Other:       "too many goroutines blocked on submit or Nonblocking is set",
	}
}

type PoolOverloadError struct {
	li18ngo.LocalisableError
}

var ErrPoolOverload = PoolOverloadError{
	LocalisableError: li18ngo.LocalisableError{
		Data: PoolOverloadErrorTemplData{},
	},
}

// ❌ ErrInvalidPreAllocSize

// InvalidPreAllocSizeErrorTemplData will be returned when trying to set up a
// negative capacity under PreAlloc mode.
type InvalidPreAllocSizeErrorTemplData struct {
	pantsTemplData
}

// Message
func (td InvalidPreAllocSizeErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-prealloc-size.error",
		Description: "error created when trying to set up a negative capacity under PreAlloc mode",
		Other:       "can not set up a negative capacity under PreAlloc mode",
	}
}

type InvalidPreAllocSizeError struct {
	li18ngo.LocalisableError
}

var ErrInvalidPreAllocSize = InvalidPreAllocSizeError{
	LocalisableError: li18ngo.LocalisableError{
		Data: InvalidPreAllocSizeErrorTemplData{},
	},
}

// ❌ ErrTimeout

// TimeoutErrorTemplData will be returned if an operation timed out.
type TimeoutErrorTemplData struct {
	pantsTemplData
}

// Message
func (td TimeoutErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "timeout.error",
		Description: "error created if an operation timed out.",
		Other:       "operation timed out",
	}
}

type TimeoutError struct {
	li18ngo.LocalisableError
}

var ErrTimeout = TimeoutError{
	LocalisableError: li18ngo.LocalisableError{
		Data: TimeoutErrorTemplData{},
	},
}

// ❌ QueueIsFullTemplData

// QueueIsFullTemplData error created when the worker queue is full.
type QueueIsFullErrorTemplData struct {
	pantsTemplData
}

// Message
func (td QueueIsFullErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "queue-is-full.error",
		Description: "error created when the worker queue is full.",
		Other:       "the queue is full",
	}
}

type QueueIsFullError struct {
	li18ngo.LocalisableError
}

var ErrQueueIsFull = QueueIsFullError{
	LocalisableError: li18ngo.LocalisableError{
		Data: QueueIsFullErrorTemplData{},
	},
}

// ❌ QueueIsReleasedTemplData

// QueueIsReleasedTemplData will be returned when trying to insert item to a
// released worker queue
type QueueIsReleasedErrorTemplData struct {
	pantsTemplData
}

// Message
func (td QueueIsReleasedErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "queue-is-released.error",
		Description: "error created when trying to insert item to a released worker queue.",
		Other:       "the queue length is zero",
	}
}

type QueueIsReleasedError struct {
	li18ngo.LocalisableError
}

var ErrQueueIsReleased = QueueIsReleasedError{
	LocalisableError: li18ngo.LocalisableError{
		Data: QueueIsReleasedErrorTemplData{},
	},
}

// ❌❌ FooBar

// FooBarTemplData - TODO: this is a none existent error that should be
// replaced by the client. Its just defined here to illustrate the pattern
// that should be used to implement i18n with extendio. Also note,
// that this message has been removed from the translation files, so
// it is not useable at run time.
type FooBarTemplData struct {
	pantsTemplData
	Path   string
	Reason error
}

// the ID should use spp/library specific code, so replace astrolib with the
// name of the library implementing this template project.
func (td FooBarTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "foo-bar.astrolib.nav",
		Description: "Foo Bar description",
		Other:       "foo bar failure '{{.Path}}' (reason: {{.Reason}})",
	}
}

// FooBarErrorBehaviourQuery used to query if an error is:
// "Failed to read directory contents from the path specified"
type FooBarErrorBehaviourQuery interface {
	FooBar() bool
}

type FooBarError struct {
	li18ngo.LocalisableError
}

// FooBar enables the client to check if error is FooBarError
// via FooBarErrorBehaviourQuery
func (e FooBarError) FooBar() bool {
	return true
}

// NewFooBarError creates a FooBarError
func NewFooBarError(path string, reason error) FooBarError {
	return FooBarError{
		LocalisableError: li18ngo.LocalisableError{
			Data: FooBarTemplData{
				Path:   path,
				Reason: reason,
			},
		},
	}
}
