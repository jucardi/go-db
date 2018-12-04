package dbx

// TODO: Implement
type ITransactional interface {
	// Begin begins a transaction
	Begin() ITransaction

	// Commit commits a transaction
	Commit() ITransaction

	// Rollback rollback a transaction
	Rollback() ITransaction
}

type ITransaction interface{

}
