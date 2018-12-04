package dbx

import (
	"gopkg.in/jucardi/go-logger-lib.v1/log"
)

// IDatabase ...
type IDatabase interface {
	// Clone clones a new db connection without search conditions
	Clone() IDatabase

	// Close close current db connection.
	Close()

	// Callbacks returns the callbacks container to be able to add callbacks on Create, Update, Delete or Query.
	Callbacks() ICallbacksManager

	// SetLogger replaces default logger
	SetLogger(log log.ILogger)

	// R returns an instance of a repository (table if SQL, collection if Mongo). Alias 'Repo'
	R(name string) IRepository

	// Repo returns an instance of a repository (table if SQL, collection if Mongo). Alias 'R'
	Repo(name string) IRepository

	// Raw executes the provided script (sql script for SQL, javascript for MongoDB) and attempts to unmarshal the result.
	// Aliases: Exec, Run
	Raw(script string, result interface{}) error

	// Exec executes the provided script (sql script for SQL, javascript for MongoDB) and attempts to unmarshal the result.
	// Aliases: Raw, Run
	Exec(script string, result interface{}) error

	// Run executes the provided script (sql script for SQL, javascript for MongoDB) and attempts to unmarshal the result.
	// Aliases: Raw, Exec
	Run(script string, result interface{}) error

	// HasRepo check has table or not
	HasRepo(name string) bool

	// CreateRepo creates a repository in the database by the given name (table if SQL, collection if Mongo).
	// Uses the reference object to create the schema (SQL).
	CreateRepo(name string, ref ...interface{}) error

	// Migrate starts a migration process using the scripts located in the 'dataDir'
	Migrate(dataDir string, failOnOrderMismatch ...bool) error
}
