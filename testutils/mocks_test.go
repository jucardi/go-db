package testutils

import (
	"fmt"
	. "github.com/jucardi/go-db"
	. "github.com/jucardi/go-testx/testx"
	"reflect"
	"testing"
)

func TestMockAll(t *testing.T) {
	db, repo, query := MockAll()

	Convey("Database Mock", t, func() {
		testMock(t, db, (*IDatabase)(nil), (*IDatabase)(nil), "IDatabase", "Clone")
		testMock(t, db, (*IDatabase)(nil), (*IRepository)(nil), "IRepository", "R", "Repo")
	})
	Convey("Repository Mock", t, func() {
		testMock(t, repo, (*IRepository)(nil), (*IQuery)(nil), "IQuery", "Where", "Not")
	})
	Convey("Query Mock", t, func() {
		testMock(t, query, (*IQuery)(nil), (*IQuery)(nil), "IQuery", "Limit", "Not", "Or", "Page", "Select", "Skip", "Sort", "Where")
	})
}

func testMock(t *testing.T, mock mockable, mockTypeRef, returnTypeRef interface{}, testInterfaceName string, listOfMethodsThatReturnRetTypeRef ... string) {

	mockType := reflect.TypeOf(mockTypeRef).Elem()
	retType := reflect.TypeOf(returnTypeRef).Elem()
	var methodList []string

	Convey(fmt.Sprintf("Functions that return %s should return the mock instance", testInterfaceName), t, func() {
		for i := 0; i < mockType.NumMethod(); i++ {
			m := mockType.Method(i)

			if m.Type.NumOut() != 1 {
				continue
			}
			outT := m.Type.Out(0)
			if outT.Implements(retType) {
				var args []reflect.Value
				for j := 0; j < m.Type.NumIn(); j++ {
					if m.Type.IsVariadic() && j == m.Type.NumIn()-1 {
						args = append(args, reflect.Zero(m.Type.In(j).Elem()))
					} else {
						args = append(args, reflect.Zero(m.Type.In(j)))
					}
				}
				methodList = append(methodList, m.Name)
				method := reflect.ValueOf(mock).MethodByName(m.Name)
				rets := method.Call(args)
				ShouldLen(rets, 1)
				ShouldImplement(returnTypeRef, rets[0].Interface())
			}
		}
	})

	Convey(fmt.Sprintf("Verify the methods that should return %s are the ones expected", testInterfaceName), t, func() {
		ShouldMatchElements(listOfMethodsThatReturnRetTypeRef, methodList)
	})

	Convey(fmt.Sprintf("Verify functions that return %s were called once", testInterfaceName), t, func() {
		for i := 0; i < mockType.NumMethod(); i++ {
			m := mockType.Method(i)

			if m.Type.NumOut() != 1 {
				continue
			}
			outT := m.Type.Out(0)
			if outT.Implements(retType) {
				ShouldEqual(1, mock.Times(m.Name))
			}
		}
	})
}
