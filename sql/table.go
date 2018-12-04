package sql

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/jucardi/go-db"
	"gopkg.in/jucardi/go-strings.v1/stringx"
	"reflect"
)

type ITable interface {
	dbx.IRepository
	// Omit specify fields that you want to ignore when saving to database for creating, updating
	Omit(columns ...string) ITable

	// ScanRows scan `*sql.Rows` to give struct
	ScanRows(rows *sql.Rows, result interface{}) error

	// Pluck used to query single column from a model as a map
	//     var ages []int64
	//     db.Find(&users).Pluck("age", &ages)
	Pluck(column string, value interface{}) ITable

	// Related get related associations
	Related(foreignKeys ...string) ITable

	// Save update value in database, if the value doesn't have primary key, will insert it
	Save(value interface{}) error

	// ModifyColumn modify column to type
	ModifyColumn(column string, typ string) error

	// DropColumn drop a column
	DropColumn(column string) error

	// Association start `Association Mode` to handler relations things easir in that mode, refer: https://jinzhu.github.io/gorm/associations.html#association-mode
	Association(column string) *gorm.Association

	// Preload preload associations with given conditions
	//    db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
	Preload(column string, conditions ...interface{}) ITable

	//// Preloads preloads relations, don`t touch out
	//Preloads(out interface{}) IQuery
}

type table struct {
	*gorm.DB
	name string
}

func newTable(db *gorm.DB, model interface{}) ITable {
	if name, ok := model.(string); ok {
		return &table{
			DB:   db.Table(name),
			name: name,
		}
	}

	return &table{
		DB:   db.Model(model),
		name: stringx.CamelToSnake(reflect.TypeOf(model).Elem().Name()),
	}
}

func (t *table) Insert(docs ...interface{}) error {
	for _, v := range docs {
		t.DB.Create(v)
	}
	return t.Error
}

func (t *table) Drop() error {
	return t.DB.DropTable(t.name).Error
}

func (t *table) Where(condition interface{}, args ...interface{}) dbx.IQuery {
	return newQuery(t.DB, t.name).Where(condition, args...)
}

func (t *table) Not(condition interface{}, args ...interface{}) dbx.IQuery {
	return newQuery(t.DB, t.name).Not(condition, args...)
}

func (t *table) AddIndex(indexName string, fields ...string) error {
	return t.DB.AddIndex(indexName, fields...).Error
}

func (t *table) DropIndex(indexName string) error {
	return t.DB.RemoveIndex(indexName).Error
}

func (t *table) AddUniqueIndex(indexName string, fields ...string) error {
	return t.DB.AddUniqueIndex(indexName, fields...).Error
}

func (t *table) Omit(columns ...string) ITable {
	t.DB.Omit(columns...)
	return t
}

func (t *table) ScanRows(rows *sql.Rows, result interface{}) error {
	return t.DB.ScanRows(rows, result)
}

func (t *table) Pluck(column string, value interface{}) ITable {
	t.DB.Pluck(column, value)
	return t
}

func (t *table) Related(foreignKeys ...string) ITable {
	t.DB.Related(nil, foreignKeys...)
	return t
}

func (t *table) Save(value interface{}) error {
	return t.DB.Save(value).Error
}

func (t *table) ModifyColumn(column string, typ string) error {
	return t.DB.ModifyColumn(column, typ).Error
}

func (t *table) DropColumn(column string) error {
	return t.DB.DropColumn(column).Error
}

func (t *table) Association(column string) *gorm.Association {
	return t.DB.Association(column)
}

func (t *table) Preload(column string, conditions ...interface{}) ITable {
	t.DB.Preload(column, conditions)
	return t
}

func (t *table) Delete(query interface{}, args ...interface{}) error {
	return t.DB.Delete(query, args...).Error
}
