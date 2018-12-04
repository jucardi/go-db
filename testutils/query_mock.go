package testutils

import (
	. "github.com/jucardi/go-db"
	"github.com/jucardi/go-db/pages"
)

// All functions that return IQuery to easily initialize
var (
	_ IQuery = (*QueryMock)(nil)
)

// QueryMock is a mock implementation of IQuery
type QueryMock struct {
	*mockBase
}

// MockQuery returns a new instance of IQuery for mocking purposes
func MockQuery() *QueryMock {
	ret := &QueryMock{
		mockBase: newMock(),
	}

	initMock(ret, (*IQuery)(nil), ret)
	return ret
}

func (q *QueryMock) Sort(fields ...string) IQuery {
	f := make([]interface{}, len(fields))
	for i, v := range fields {
		f[i] = v
	}
	return q.returnQuery("Sort", f...)
}

func (q *QueryMock) Skip(n int) IQuery {
	return q.returnQuery("Skip", n)
}

func (q *QueryMock) Limit(n int) IQuery {
	return q.returnQuery("Limit", n)
}

func (q *QueryMock) Select(query interface{}, args ...interface{}) IQuery {
	return q.returnQuery("Select", query, args)
}

func (q *QueryMock) Where(condition interface{}, args ...interface{}) IQuery {
	return q.returnQuery("Where", condition, args)
}

func (q *QueryMock) Not(condition interface{}, args ...interface{}) IQuery {
	return q.returnQuery("Not", condition, args)
}

func (q *QueryMock) Or() IQuery {
	return q.returnQuery("Or")
}

func (q *QueryMock) First(result interface{}) error {
	return q.ReturnError("First", result)
}

func (q *QueryMock) One(result interface{}) error {
	return q.ReturnError("One", result)
}

func (q *QueryMock) Last(result interface{}) error {
	return q.ReturnError("Last", result)
}

func (q *QueryMock) Distinct(key string, result interface{}) error {
	return q.ReturnError("Distinct", key, result)
}

func (q *QueryMock) Update(update interface{}, result interface{}) error {
	return q.ReturnError("Update", update, result)
}

func (q *QueryMock) Delete() error {
	return q.ReturnError("Delete")
}

func (q *QueryMock) Remove() error {
	return q.ReturnError("Remove")
}

func (q *QueryMock) Page(page ...*pages.Page) IQuery {
	p := make([]interface{}, len(page))
	for i, v := range page {
		p[i] = v
	}
	return q.returnQuery("Page", p...)
}

func (q *QueryMock) WrapPage(result interface{}, page ...*pages.Page) (*pages.Paginated, error) {
	args := []interface{}{result}
	for _, v := range page {
		args = append(args, v)
	}
	ret, err := q.ReturnSingleArgWithError("WrapPage", args...)

	if ret != nil {
		return ret.(*pages.Paginated), err
	}

	return nil, err
}

func (q *QueryMock) Count() (int, error) {
	ret, err := q.ReturnSingleArgWithError("Count")

	if ret != nil {
		return ret.(int), err
	}

	return 0, err
}

func (q *QueryMock) All(result interface{}) error {
	return q.ReturnError("All", result)
}
