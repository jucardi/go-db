package testutils

import (
	. "github.com/jucardi/go-db"
	"gopkg.in/jucardi/go-logger-lib.v1/log"
)

// All functions that return IQuery to easily initialize
var (
	_ IDatabase = (*DatabaseMock)(nil)
)

// QueryMock is a mock implementation of IQuery
type DatabaseMock struct {
	*mockBase
}

// MockDB returns a new instance of IDatabase for mocking purposes
func MockDB() *DatabaseMock {
	ret := &DatabaseMock{
		mockBase: newMock(),
	}

	initMock(ret, (*IDatabase)(nil), ret)
	return ret
}

func (db *DatabaseMock) Clone() IDatabase {
	return db.returnDB("Clone")
}

func (db *DatabaseMock) Close() {
	db.Invoke("Close")
}

func (db *DatabaseMock) Callbacks() ICallbacksManager {
	ret := db.Invoke("Callbacks")
	if len(ret) == 0 || ret[0] == nil {
		return nil
	}
	return ret[0].(ICallbacksManager)
}

func (db *DatabaseMock) SetLogger(log log.ILogger) {
	db.Invoke("SetLogger")
}

func (db *DatabaseMock) R(name string) IRepository {
	return db.returnRepository("R", name)
}

func (db *DatabaseMock) Repo(name string) IRepository {
	return db.returnRepository("Repo", name)
}

func (db *DatabaseMock) Exec(script string, result interface{}) error {
	return db.ReturnError("Exec", script, result)
}

func (db *DatabaseMock) Run(script string) error {
	return db.ReturnError("Run", script)
}

func (db *DatabaseMock) HasRepo(name string) bool {
	return db.ReturnBool("HasRepo", name)
}

func (db *DatabaseMock) CreateRepo(name string, ref ...interface{}) error {
	return db.ReturnError("CreateRepo", name, ref)
}

func (db *DatabaseMock) Migrate(dataDir string, failOnOrderMismatch ...bool) error {
	return db.ReturnError("Migrate", dataDir, failOnOrderMismatch)
}

func (db *DatabaseMock) SetScriptExecutor(executor ScriptExecutor) {
	db.Invoke("SetScriptExecutor")
}