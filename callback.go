//

package dbx

// ICallbacksManager ...
type ICallbacksManager interface {
	// Create could be used to register callbacks for creating object
	//     db.Callback().Create().After("gorm:create").Register("plugin:run_after_create", func(*Scope) {
	//       // business logic
	//       ...
	//
	//       // set error if some thing wrong happened, will rollback the creating
	//       scope.Err(errors.New("error"))
	//     })
	Create()
	// Update could be used to register callbacks for updating object, refer `Create` for usage
	Update() ICallbacksHandler
	// Delete could be used to register callbacks for deleting object, refer `Create` for usage
	Delete() ICallbacksHandler
	// Query could be used to register callbacks for querying objects with query methods like `Find`, `First`, `Related`, `Association`...
	// Refer `Create` for usage
	Query() ICallbacksHandler
}
