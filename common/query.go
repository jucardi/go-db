package common

import (
	"github.com/jucardi/go-db"
)

// AbstractQuery is a base query helper used for final implementations.
type AbstractQuery struct {
	Q          dbx.IQuery
	LimitVal   *int
	SkipVal    *int
	SortFields []string
	Selects    []*ConditionData
	Queries    [][]*ConditionData
}

type ConditionData struct {
	Negation bool
	Query    interface{}
	Args     []interface{}
}

func newCond(query interface{}, args []interface{}, negated bool) *ConditionData {
	return &ConditionData{Query: query, Args: args, Negation: negated}
}

func (a *AbstractQuery) Limit(n int) dbx.IQuery {
	a.LimitVal = &n
	return a.Q
}

func (a *AbstractQuery) Skip(n int) dbx.IQuery {
	a.SkipVal = &n
	return a.Q
}

func (a *AbstractQuery) Sort(fields ...string) dbx.IQuery {
	for _, f := range fields {
		a.SortFields = append(a.SortFields, f)
	}
	return a.Q
}

func (a *AbstractQuery) Select(query interface{}, args ...interface{}) dbx.IQuery {
	a.Selects = append(a.Selects, newCond(query, args, false))
	return a.Q
}

func (a *AbstractQuery) Where(query interface{}, args ...interface{}) dbx.IQuery {
	a.addToBlock(newCond(query, args, false))
	return a.Q
}

func (a *AbstractQuery) Not(query interface{}, args ...interface{}) dbx.IQuery {
	a.addToBlock(newCond(query, args, true))
	return a.Q
}

func (a *AbstractQuery) Or() dbx.IQuery {
	a.newBlock()
	return a.Q
}

func (a *AbstractQuery) newBlock() {
	a.Queries = append(a.Queries, []*ConditionData{})
}

func (a *AbstractQuery) addToBlock(data *ConditionData) {
	if len(a.Queries) == 0 {
		a.newBlock()
	}

	a.Queries[len(a.Queries)-1] = append(a.Queries[len(a.Queries)-1], data)
}
