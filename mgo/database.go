package mgo

import (
	"github.com/jucardi/go-db"
	"github.com/jucardi/go-db/common"
	l "github.com/jucardi/go-db/logger"
	"github.com/jucardi/go-logger-lib/log"
	"github.com/jucardi/go-streams/streams"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type IDatabase interface {
	dbx.IDatabase

	// C returns a value representing the named collection. Alias 'Collection'
	C(name string) ICollection

	// C returns a value representing the named collection. Alias 'C'
	Collection(name string) ICollection

	// GridFS returns a GridFS value representing collections in db that follow the standard GridFS specification.
	GridFS(prefix string) *mgo.GridFS

	// Login authenticates with MongoDB using the provided credential. The authentication is valid for the whole Session and will stay valid until Logout is explicitly called for
	// the same Database, or the Session is closed.
	Login(user, pass string) error

	// Logout removes any established authentication credentials for the Database.
	Logout()

	// UpsertUser updates the authentication credentials and the roles for a MongoDB user within the db Database. If the named user doesn't exist it will be created.
	// This method should only be used from MongoDB 2.4 and on. For older MongoDB releases, use the obsolete AddUser method instead.
	UpsertUser(user *mgo.User) error

	// AddUser creates or updates the authentication credentials of user within the db Database.
	// WARNING: This method is obsolete and should only be used with MongoDB 2.2 or earlier. For MongoDB 2.4 and on, use UpsertUser instead.
	AddUser(username, password string, readOnly bool) error

	// RemoveUser removes the authentication credentials of user from the Database.
	RemoveUser(user string) error

	// DropDatabase removes the entire Database including all of its collections.
	DropDatabase() error

	// FindRef returns a query that looks for the document in the provided reference. If the reference includes the DB field, the document will be retrieved from the respective Database.
	FindRef(ref *mgo.DBRef) IQuery

	// CollectionNames returns the collection names present in the db Database.
	CollectionNames() (names []string, err error)

	// Name returns the name of the database
	Name() string

	// Session returns the session used by the database
	Session() ISession

	// Returns the internal mgo.Database used by this implementation.
	DB() *mgo.Database
}

type database struct {
	*mgo.Database
	executor dbx.ScriptExecutor
}

func (d *database) Clone() dbx.IDatabase {
	return d.Session().Clone().DB(d.Name())
}

func (d *database) Close() {
	d.Session().Close()
}

func (d *database) Callbacks() dbx.ICallbacksManager {
	panic("implement me")
}

func (d *database) SetLogger(l log.ILogger) {
	mgo.SetLogger(wrapLogger(l))
	mgo.SetDebug(l.GetLevel() == log.DebugLevel)
}

func (d *database) Repo(name string) dbx.IRepository {
	return d.C(name)
}

func (d *database) R(name string) dbx.IRepository {
	return d.C(name)
}

func (d *database) Collection(name string) ICollection {
	return d.C(name)
}

func (d *database) C(name string) ICollection {
	return fromCollection(d.DB().C(name))
}

func (d *database) Exec(script string, result interface{}) error {
	return d.DB().Run(bson.M{"eval": script}, result)
}

func (d *database) Run(script string) error {
	if d.executor != nil {
		return d.executor(script)
	}
	return d.DB().Run(bson.M{"eval": script}, nil)
}

func (d *database) CreateRepo(name string, models ...interface{}) error {
	// MongoDB auto-creates collections so this is not required. Only present for compatibility
	return nil
}

func (d *database) HasRepo(name string) bool {
	cols, err := d.CollectionNames()
	if err != nil {
		l.Get().Error("Error obtaining Mongo Collections - ", err.Error())
	}
	return streams.From(cols).Contains(name)
}

func (d *database) Migrate(dataDir string, failOnOrderMismatch ...bool) error {
	fail := true
	if len(failOnOrderMismatch) > 0 {
		fail = failOnOrderMismatch[0]
	}
	return common.Migrate(dataDir, d, fail)
}

func (d *database) SetScriptExecutor(executor dbx.ScriptExecutor) {
	d.executor = executor
}

func (d *database) DB() *mgo.Database {
	return d.Database
}

func (d *database) With(s ISession) IDatabase {
	return FromDB(d.DB().With(s.S()))
}

func (d *database) FindRef(ref *mgo.DBRef) IQuery {
	var c ICollection
	if ref.Database == "" {
		c = d.C(ref.Collection)
	} else {
		c = d.Session().DB(ref.Database).C(ref.Collection)
	}
	return c.FindId(ref.Id)
}

func (d *database) Name() string {
	return d.DB().Name
}

func (d *database) Session() ISession {
	return fromSession(d.DB().Session)
}

func FromDB(db *mgo.Database) IDatabase {
	if db == nil {
		return nil
	}
	return &database{Database: db}
}
