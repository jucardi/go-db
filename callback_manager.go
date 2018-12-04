package dbx

type CallbackFunc func()

type ICallbacksHandler interface {
	// Before adds a new callback which will be executed before the operation.
	Before(name string, handler CallbackFunc) ICallbacksHandler

	// After adds a new callback which will be executed after the operation.
	After(name string, handler CallbackFunc) ICallbacksHandler

	// OnError
	OnError(name string, handler CallbackFunc) ICallbacksHandler

	// Remove a registered callback
	Remove(name string)
}
