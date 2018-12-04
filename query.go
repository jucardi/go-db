package dbx

import "github.com/jucardi/go-db/pages"

// IQuery is an interface which matches the contract for the `query` struct in `gopkg.in/mgo.v2` package. The function documentation has been narrowed from the original
// in `gopkg.in/mgo.v2`. For additional documentation, please refer to the `mgo.Collection` in the `gopkg.in/mgo.v2` package.
type IQuery interface {
	// Set of extension functions that are not present in the original `mgo` package are defined in the following interface(s):
	IQueryPageExtension

	// Limit restricts the maximum number of records retrieved to n, and also changes the batch size to the same value.
	Limit(n int) IQuery

	// Skip skips over the n initial documents from the query results. Note that this only makes sense with capped collections where documents are naturally ordered by insertion
	// time, or with sorted results.
	Skip(n int) IQuery

	// Sort asks the Database to order returned records according to the provided field names. A field name may be prefixed by - (minus) for it to be sorted in reverse order.
	Sort(fields ...string) IQuery

	// Select specify fields that you want to retrieve from database when querying, by default, will select all fields;
	// When creating/updating, specify fields that you want to save to database
	Select(query interface{}, args ...interface{}) IQuery

	// Where adds an additional condition to match records in the query. It is associated with the previous queries with an AND prepares a query script using the provided condition object to match any records that meet the provided condition.
	// Accepts `map`, `struct`. For SQL, it also accepts `string` as conditions, refer http://jinzhu.github.io/gorm/crud.html#query
	Where(condition interface{}, args ...interface{}) IQuery

	// Not prepares a query script using the provided condition object to match any records that do not meet the provided condition.
	// Accepts `map`, `struct`. For SQL, it also accepts `string` as conditions, refer http://jinzhu.github.io/gorm/crud.html#query
	Not(condition interface{}, args ...interface{}) IQuery

	// Or indicates that any following queries in the chain will be OR'ed with the previous queries
	Or() IQuery

	// Count returns the total number of records in the result set.
	Count() (n int, err error)

	// First executes the query and unmarshals the first obtained record into the result argument. Alias 'One'
	First(result interface{}) error

	// One executes the query and unmarshals the first obtained record into the result argument. Alias 'First'
	One(result interface{}) error

	// Last executes the query and unmarshals the last obtained record into the result argument
	Last(result interface{}) error

	// All works like Iter.All.
	All(result interface{}) error

	// Distinct unmarshals into result the list of distinct values for the given key.
	Distinct(key string, result interface{}) error

	// Update update attributes with callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
	Update(update interface{}) error

	// Delete deletes the records resulting from executing the query. Alias 'Remove'
	Delete() error

	// Remove deletes the records resulting from executing the query. Alias 'Delete'
	Remove() error
}

// IQueryPageExtension encapsulates the new extended functions to the original IQuery
type IQueryPageExtension interface {
	// Page adds to the query the information required to fetch the requested page of objects.
	Page(p ...*pages.Page) IQuery

	// WrapPage attempts to obtain the items in the requested page and wraps the result in *pages.Paginated
	WrapPage(result interface{}, p ...*pages.Page) (*pages.Paginated, error)
}