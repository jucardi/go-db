package testutils

import (
	. "github.com/jucardi/go-db"
	"github.com/jucardi/gotestx/mock"
	"reflect"
)

type mockable interface {
	WhenReturn(funcName string, retArgs ...interface{})
	Times(funcName string) int
}

type mockBase struct {
	*mock.MockBase
}

func (m *mockBase) returnQuery(name string, args ...interface{}) IQuery {
	if val, ok := m.ReturnSingleArg(name, args...).(IQuery); ok {
		return val
	}
	return nil
}

func (m *mockBase) returnDB(name string, args ...interface{}) IDatabase {
	if val, ok := m.ReturnSingleArg(name, args...).(IDatabase); ok {
		return val
	}
	return nil
}

func (m *mockBase) returnRepository(name string, args ...interface{}) IRepository {
	if val, ok := m.ReturnSingleArg(name, args...).(IRepository); ok {
		return val
	}
	return nil
}

func MockAll() (db *DatabaseMock, repo *RepositoryMock, query *QueryMock) {
	db = MockDB()
	repo = MockRepo()
	query = MockQuery()
	db.WhenReturn("R", repo)
	db.WhenReturn("Repo", repo)
	repo.WhenReturn("Where", query)
	repo.WhenReturn("Not", query)
	return db, repo, query
}

func newMock() *mockBase {
	return &mockBase{MockBase: mock.New()}
}

func initMock(mock mockable, ref interface{}, returnObj interface{}) {
	// Sets self return values
	qType := reflect.TypeOf(ref).Elem()
	eType := reflect.TypeOf((*error)(nil)).Elem()

	for i := 0; i < qType.NumMethod(); i++ {
		m := qType.Method(i)

		if m.Type.NumOut() != 1 {
			continue
		}
		outT := m.Type.Out(0)
		if outT.Implements(qType) {
			mock.WhenReturn(m.Name, returnObj)
		}
		if outT.Implements(eType) {
			mock.WhenReturn(m.Name, nil)
		}
	}
}
