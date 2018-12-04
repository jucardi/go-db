package dbx

// IRepository represents a repository of records (table for SQL, collection for MongoDB)
type IRepository interface {
	// Insert inserts one or more records in the respective repository.
	Insert(docs ...interface{}) error

	// Drop drops the repository (table if SQL, collection if MongoDB)
	Drop() error

	// Where prepares a query script using the provided condition object to match any records that meet the provided condition.
	// Accepts `map`, `struct`. For SQL, it also accepts `string` as conditions, refer http://jinzhu.github.io/gorm/crud.html#query/
	Where(condition interface{}, args ...interface{}) IQuery

	// Not prepares a query script using the provided condition object to match any records that do not meet the provided condition.
	// Accepts `map`, `struct`. For SQL, it also accepts `string` as conditions, refer http://jinzhu.github.io/gorm/crud.html#query
	Not(condition interface{}, args ...interface{}) IQuery

	// AddIndex adds an index for the provided fields (columns for SQL, keys for MongoDB)
	AddIndex(indexName string, fields ...string) error

	// Removes all records that meet the provided query
	Delete(query interface{}, args ...interface{}) error

	// DropIndex remove index with name
	DropIndex(indexName string) error

	// AddUniqueIndex adds a unique index for the provided fields (columns for SQL, keys for MongoDB)
	AddUniqueIndex(indexName string, fields ...string) error
}
