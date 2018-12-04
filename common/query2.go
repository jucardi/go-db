package common

import (
	"github.com/jucardi/go-db"
)

// AbstractQueryx is a base query helper used for final implementations.
type AbstractQueryx struct {
	Q        dbx.IQuery
	Commands []*Command
}

type Command struct {
	Name string
	Args []interface{}
}

func strToArgs(args []string) (ret []interface{}) {
	for _, v := range args {
		ret = append(ret, v)
	}
	return
}

func (a *AbstractQueryx) addCmd(name string, args ...interface{}) dbx.IQuery {
	a.Commands = append(a.Commands, &Command{Name: name, Args: args})
	return a.Q
}

func (a *AbstractQueryx) Limit(n int) dbx.IQuery {
	return a.addCmd("Limit", n)
}

func (a *AbstractQueryx) Skip(n int) dbx.IQuery {
	return a.addCmd("Skip", n)
}

func (a *AbstractQueryx) Sort(fields ...string) dbx.IQuery {
	return a.addCmd("Sort", strToArgs(fields))
}

func (a *AbstractQueryx) Select(query interface{}, args ...interface{}) dbx.IQuery {
	return a.addCmd("Select", query, args)
}

func (a *AbstractQueryx) Where(query interface{}, args ...interface{}) dbx.IQuery {
	return a.addCmd("Where", query, args)
}

func (a *AbstractQueryx) Not(query interface{}, args ...interface{}) dbx.IQuery {
	return a.addCmd("Not", query, args)
}

func (a *AbstractQueryx) Or() dbx.IQuery {
	return a.addCmd("Or")
}
