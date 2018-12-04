package testutils

import (
	. "github.com/jucardi/go-db"
)

// All functions that return IQuery to easily initialize
var (
	_ IRepository = (*RepositoryMock)(nil)
)

// RepositoryMock is a mock implementation of IRepository
type RepositoryMock struct {
	*mockBase
}

// MockRepo returns a new instance of IRepository for mocking purposes
func MockRepo() *RepositoryMock {
	ret := &RepositoryMock{
		mockBase: newMock(),
	}

	return ret
}

func (r *RepositoryMock) Insert(docs ...interface{}) error {
	return r.ReturnError("Insert", docs...)
}

func (r *RepositoryMock) Drop() error {
	return r.ReturnError("Drop")
}

func (r *RepositoryMock) Where(condition interface{}, args ...interface{}) IQuery {
	return r.returnQuery("Where", condition, args)
}

func (r *RepositoryMock) Not(condition interface{}, args ...interface{}) IQuery {
	return r.returnQuery("Not", condition, args)
}

func (r *RepositoryMock) AddIndex(indexName string, fields ...string) error {
	return r.ReturnError("AddIndex", indexName, fields)
}

func (r *RepositoryMock) DropIndex(indexName string) error {
	return r.ReturnError("DropIndex", indexName)
}

func (r *RepositoryMock) AddUniqueIndex(indexName string, fields ...string) error {
	return r.ReturnError("AddUniqueIndex", indexName, fields)
}
