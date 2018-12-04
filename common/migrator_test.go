package common

import (
	"errors"
	"github.com/jucardi/go-db"
	"github.com/jucardi/go-db/logger"
	"github.com/jucardi/go-db/testutils"
	"github.com/jucardi/gotestx/mock"
	. "github.com/jucardi/gotestx/testx"
	"gopkg.in/jucardi/go-logger-lib.v1/log"
	"testing"
)

const migrationPath = "./test_assets/db_migration"

func init() {
	logger.Set(log.NewNil())
}

func TestMigrateSuccess(t *testing.T) {
	db, repo, q := testutils.MockAll()

	Convey("Migrate Success", t, func() {
		err := Migrate(migrationPath, db, true)
		ShouldEqual(1, repo.Times("Where"))
		ShouldEqual(2, repo.Times("Insert"))
		ShouldEqual(3, db.Times("R"))
		ShouldEqual(2, db.Times("Run"))
		ShouldEqual(1, q.Times("Sort"))
		ShouldEqual(1, q.Times("All"))
		ShouldBeNil(err)
	})
}

func TestReadMigrateCollectionFailed(t *testing.T) {
	msg := "some error"
	db, repo, q := testutils.MockAll()
	Convey("Migrate Failed - DB Operation", t, func() {
		q.WhenReturn("All", errors.New(msg))

		err := Migrate(migrationPath, db, true)
		ShouldError(err)
		ShouldBeTrue(err.Is(dbx.ErrDbOperation))
		ShouldEqual("Unable to read Database info. "+msg, err.Error())

		ShouldEqual(1, repo.Times("Where"))
		ShouldEqual(0, repo.Times("Insert"))
		ShouldEqual(1, db.Times("R"))
		ShouldEqual(0, db.Times("Run"))
		ShouldEqual(1, q.Times("Sort"))
		ShouldEqual(1, q.Times("All"))
	})

}

func TestReadScriptPathFailed(t *testing.T) {
	path := "some-invalid-path"
	db, repo, q := testutils.MockAll()

	Convey("Migrate Failed - DB Operation", t, func() {
		err := Migrate(path, db, true)
		ShouldError(err)
		ShouldBeTrue(err.Is(dbx.ErrFileAccess))
		ShouldEqual("Unable to access scripts path. open "+path+": no such file or directory", err.Error())

		ShouldEqual(1, repo.Times("Where"))
		ShouldEqual(0, repo.Times("Insert"))
		ShouldEqual(1, db.Times("R"))
		ShouldEqual(0, db.Times("Run"))
		ShouldEqual(1, q.Times("Sort"))
		ShouldEqual(1, q.Times("All"))
	})

}

func TestPreviousDataSuccess(t *testing.T) {
	db, repo, q := testutils.MockAll()
	Convey("Migrate Success after previous migration", t, func() {
		q.When("All", func(args ...interface{}) []interface{} {
			result := args[0]
			list := result.(*[]*MigrationInfo)
			*list = append(*list, &MigrationInfo{
				ScriptId: "script_001.js",
				Hash:     "b280f134425a4153026cf227069d4cc1",
			})
			return mock.MakeReturn(nil)
		})

		ShouldBeNil(Migrate(migrationPath, db, true))

		ShouldEqual(1, repo.Times("Where"))
		ShouldEqual(1, repo.Times("Insert"))
		ShouldEqual(2, db.Times("R"))
		ShouldEqual(1, db.Times("Run"))
		ShouldEqual(1, q.Times("Sort"))
		ShouldEqual(1, q.Times("All"))
	})
}

func TestPreviousDataFailedHash(t *testing.T) {
	db, repo, q := testutils.MockAll()
	Convey("Migrate Failed - Previous script hashes do not match", t, func() {
		q.When("All", func(args ...interface{}) []interface{} {
			result := args[0]
			list := result.(*[]*MigrationInfo)
			*list = append(*list, &MigrationInfo{
				ScriptId: "script_001.js",
				Hash:     "1234",
			})
			return mock.MakeReturn(nil)
		})

		err := Migrate(migrationPath, db, true)
		ShouldError(err)
		ShouldEqual("File 'script_001.js' was previously migrated but hashes don't match.", err.Error())

		ShouldEqual(1, repo.Times("Where"))
		ShouldEqual(0, repo.Times("Insert"))
		ShouldEqual(1, db.Times("R"))
		ShouldEqual(0, db.Times("Run"))
		ShouldEqual(1, q.Times("Sort"))
		ShouldEqual(1, q.Times("All"))
	})
}

func TestRunFailed(t *testing.T) {
	db, repo, q := testutils.MockAll()
	Convey("Migrate Failed - Script failed to run", t, func() {
		db.WhenReturn("Run", errors.New("some error"))

		err := Migrate(migrationPath, db, true)
		ShouldError(err)
		ShouldEqual("Unable to run command 'script_001.js'. some error", err.Error())

		ShouldEqual(1, repo.Times("Where"))
		ShouldEqual(0, repo.Times("Insert"))
		ShouldEqual(1, db.Times("R"))
		ShouldEqual(1, db.Times("Run"))
		ShouldEqual(1, q.Times("Sort"))
		ShouldEqual(1, q.Times("All"))
	})

}

func TestInsertFailed(t *testing.T) {
	db, repo, q := testutils.MockAll()
	Convey("Migrate Failed -  Failed to insert migration data", t, func() {
		repo.WhenReturn("Insert", errors.New("some error"))

		err := Migrate(migrationPath, db, true)
		ShouldError(err)
		ShouldEqual("Unable to save migration info for 'script_001.js'", err.Error())

		ShouldEqual(1, repo.Times("Where"))
		ShouldEqual(1, repo.Times("Insert"))
		ShouldEqual(2, db.Times("R"))
		ShouldEqual(1, db.Times("Run"))
		ShouldEqual(1, q.Times("Sort"))
		ShouldEqual(1, q.Times("All"))
	})

}
