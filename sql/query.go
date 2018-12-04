package sql

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/jucardi/go-db"
	"github.com/jucardi/go-db/common"
	"github.com/jucardi/go-db/pages"
	"strings"
)

type IQuery interface {
	dbx.IQuery

	// Group specify the group method on the find
	Group(query string) IQuery

	// Having specify HAVING conditions for GROUP BY
	Having(query interface{}, values ...interface{}) IQuery

	// Joins specify Joins conditions
	//     db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Find(&user)
	Joins(query string, args ...interface{}) IQuery

	// Attrs initialize struct with argument if record not found with `FirstOrInit` https://jinzhu.github.io/gorm/crud.html#firstorinit or `FirstOrCreate` https://jinzhu.github.io/gorm/crud.html#firstorcreate
	Attrs(attrs ...interface{}) IQuery
	// Assign assign result with argument regardless it is found or not with `FirstOrInit` https://jinzhu.github.io/gorm/crud.html#firstorinit or `FirstOrCreate` https://jinzhu.github.io/gorm/crud.html#firstorcreate

	Assign(attrs ...interface{}) IQuery

	// Take return a record that match given conditions, the order will depend on the database implementation
	Take(out interface{}, where ...interface{}) IQuery

	// Scan scan value to a struct
	Scan(dest interface{}) IQuery

	// FirstOrInit find first matched record or initialize a new one with given conditions (only works with struct, map conditions)
	// https://jinzhu.github.io/gorm/crud.html#firstorinit
	FirstOrInit(out interface{}, where ...interface{}) IQuery

	// FirstOrCreate find first matched record or create a new one with given conditions (only works with struct, map conditions)
	// https://jinzhu.github.io/gorm/crud.html#firstorcreate
	FirstOrCreate(out interface{}, where ...interface{}) IQuery

	// UpdateColumn update attributes without callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
	UpdateColumn(attrs ...interface{}) IQuery

	// UpdateColumns update attributes without callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
	UpdateColumns(values interface{}) IQuery

	// Row return `*sql.Row` with given conditions
	Row() *sql.Row

	// Rows return `*sql.Rows` with given conditions
	Rows() (*sql.Rows, error)

	Error() error
}

func newQuery(db *gorm.DB, table string) *query {
	return &query{
		table:   table,
		DB:      db,
		queries: []*gorm.DB{db.Table(table)},
	}
}

type query struct {
	*gorm.DB
	queries []*gorm.DB
	table   string
}

func (q *query) current() *gorm.DB {
	return q.queries[len(q.queries)-1]
}

// Page adds to the query the information required to fetch the requested page of objects.
func (q *query) Page(page ...*pages.Page) dbx.IQuery {
	return common.PageHandler(q, page...)
}

// WrapPage attempts to obtain the items in the requested page and wraps the result in *pages.Paginated
func (q *query) WrapPage(result interface{}, page ...*pages.Page) (*pages.Paginated, error) {
	return common.WrapPageHandler(q, result, page...)
}

func (q *query) Limit(n int) dbx.IQuery {
	q.DB.Limit(n)
	return q
}

func (q *query) Skip(n int) dbx.IQuery {
	q.DB.Offset(n)
	return q
}

func (q *query) Sort(fields ...string) dbx.IQuery {
	for _, field := range fields {
		if strings.HasPrefix(field, "-") {
			q.DB.Order(field[1:] + " desc")
		} else {
			q.DB.Order(field)
		}
	}
	return q
}

func (q *query) Select(query interface{}, args ...interface{}) dbx.IQuery {
	q.current().Select(query, args...)
	return q
}

func (q *query) Where(condition interface{}, args ...interface{}) dbx.IQuery {
	q.current().Where(condition, args...)
	return q
}

func (q *query) Not(condition interface{}, args ...interface{}) dbx.IQuery {
	q.current().Not(condition, args...)
	return q
}

func (q *query) Or() dbx.IQuery {
	q.queries = append(q.queries, q.DB.Table(q.table))
	return q
}

func (q *query) Count() (n int, err error) {
	err = q.prepare().Count(&n).Error
	return
}

func (q *query) First(result interface{}) error {
	return q.prepare().First(result).Error
}

func (q *query) One(result interface{}) error {
	return q.First(result)
}

func (q *query) Last(result interface{}) error {
	return q.prepare().Last(result).Error
}

func (q *query) All(result interface{}) error {
	return q.prepare().First(result).Error
}

func (q *query) Distinct(key string, result interface{}) error {
	return q.prepare().Select("DISTINCT ?", key).Scan(result).Error
}

func (q *query) Update(update interface{}) error {
	return q.prepare().Updates(update).Error
}

func (q *query) Delete() error {
	panic("")
}

func (q *query) Remove() error {
	panic("implement me")
}

func (q *query) Row() *sql.Row {
	return q.prepare().Row()
}

func (q *query) Rows() (*sql.Rows, error) {
	return q.prepare().Rows()
}

func (q *query) prepare() *gorm.DB {
	for _, qry := range q.queries {
		q.DB.Where(qry.SubQuery())
	}
	return q.DB
}