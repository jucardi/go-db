package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/jucardi/go-db"
	"github.com/jucardi/go-db/common"
	"github.com/jucardi/go-db/logger"
	"gopkg.in/jucardi/go-logger-lib.v1/log"
	"gopkg.in/jucardi/go-streams.v1/streams"
)

type IDatabase interface {
	dbx.IDatabase

	// Scopes pass current database connection to arguments `func(*DB) *DB`, which could be used to add conditions dynamically
	//     func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
	//         return db.Where("amount > ?", 1000)
	//     }
	//
	//     func OrderStatus(status []string) func (db *gorm.DB) *gorm.DB {
	//         return func (db *gorm.DB) *gorm.DB {
	//             return db.Scopes(AmountGreaterThan1000).Where("status in (?)", status)
	//         }
	//     }
	//
	//     db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
	// Refer https://jinzhu.github.io/gorm/crud.html#scopes
	Scopes(funcs ...func(IDatabase) IDatabase) IDatabase

	// Unscoped return all record including deleted record, refer Soft Delete https://jinzhu.github.io/gorm/crud.html#soft-delete
	Unscoped() IDatabase

	// Model specify the model you would like to run db operations
	//    // update all users's name to `hello`
	//    db.Model(&User{}).Update("name", "hello")
	//    // if user's primary key is non-blank, will use it as condition, then will only update the user's name to `hello`
	//    db.Model(&user).Update("name", "hello")
	Model(value interface{}) ITable

	// Table specifies the table you would like to run db operations. Alias 'T'
	Table(name string) ITable

	// T specifies the table you would like to run db operations. Alias 'Table'
	T(name string) ITable

	// Debug starts debug mode
	Debug()

	// AddForeignKey Add foreign key to the given scope, e.g:
	//     db.Model(&User{}).AddForeignKey("city_id", "cities(id)", "RESTRICT", "RESTRICT")
	AddForeignKey(field string, dest string, onDelete string, onUpdate string) error

	// RemoveForeignKey Remove foreign key from the given scope, e.g:
	//     db.Model(&User{}).RemoveForeignKey("city_id", "cities(id)")
	RemoveForeignKey(field string, dest string) error

	Db() *gorm.DB
}

type database struct {
	*gorm.DB
	isClone bool
}

func FromDB(db *gorm.DB, isClone bool) IDatabase {
	return &database{
		DB:      db,
		isClone: isClone,
	}
}

func (db *database) Db() *gorm.DB {
	return db.DB
}

func (db *database) Clone() dbx.IDatabase {
	return FromDB(db.DB.New(), true)
}

func (db *database) Close() {
	if db.isClone {
		return
	}

	if err := db.DB.Close(); err != nil {
		logger.Get().Errorf("Error occurred while closing SQL DB connection, %s", err.Error())
	}
}

func (db *database) Callbacks() dbx.ICallbacksManager {
	panic("implement me")
}

func (db *database) SetLogger(l log.ILogger) {
	db.DB.SetLogger(wrapLogger(l))
	if l.GetLevel() == log.DebugLevel {
		db.DB.Debug()
	}
}

func (db *database) R(name string) dbx.IRepository {
	return db.Table(name)
}

func (db *database) Repo(name string) dbx.IRepository {
	return db.Table(name)
}

func (db *database) Raw(script string, result interface{}) error {
	return db.Exec(script, result)
}

func (db *database) Exec(script string, result interface{}) error {
	return db.DB.Raw(script).Scan(result).Error
}

func (db *database) Run(script string) error {
	return db.DB.Exec(script).Error
}

func (db *database) HasRepo(name string) bool {
	return db.DB.HasTable(name)
}

func (db *database) CreateRepo(name string, models ...interface{}) error {
	return db.DB.Table(name).CreateTable(models...).Error
}

func (db *database) Migrate(dataDir string, failOnOrderMismatch ...bool) error {
	fail := true
	if len(failOnOrderMismatch) > 0 {
		fail = failOnOrderMismatch[0]
	}
	if err := common.Migrate(dataDir, db, fail); err != nil {
		return err
	}
	return nil
}

func (db *database) Scopes(funcs ...func(IDatabase) IDatabase) IDatabase {
	fs := streams.From(funcs).Map(func(i interface{}) interface{} {
		f := i.(func(IDatabase) IDatabase)
		return func(x *gorm.DB) *gorm.DB {
			return f(FromDB(x, false)).Db()
		}
	}).ToArray()

	if val, ok := fs.([]func(*gorm.DB) *gorm.DB); ok {
		db.DB.Scopes(val...)
	}
	return db
}

func (db *database) Unscoped() IDatabase {
	return FromDB(db.DB.Unscoped(), false)
}

func (db *database) Model(value interface{}) ITable {
	return newTable(db.DB, value)
}

func (db *database) Table(name string) ITable {
	return newTable(db.DB, name)
}

func (db *database) T(name string) ITable {
	return db.Table(name)
}

func (db *database) Debug() {
	db.DB.Debug()
}

func (db *database) AddForeignKey(field string, dest string, onDelete string, onUpdate string) error {
	return db.DB.AddForeignKey(field, dest, onDelete, onUpdate).Error
}

func (db *database) RemoveForeignKey(field string, dest string) error {
	return db.DB.RemoveForeignKey(field, dest).Error
}
