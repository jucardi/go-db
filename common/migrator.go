package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/jucardi/go-db"
	"github.com/jucardi/go-db/logger"
	"github.com/jucardi/go-osx/paths"
	"github.com/jucardi/go-streams/streams"
	"gopkg.in/mgo.v2/bson"
)

const (
	MigrationRepo = "_migration"
)

type Migrator struct {
	// Db is the database client already initialized
	Db dbx.IDatabase

	// DataDir is the location where the migration scripts are contained
	DataDir string

	// FailOnOrderMismatch indicates whether the migration should fail if the order of previously migrated scripts
	FailOnOrderMismatch bool

	// RepoIdSuffix is an optional suffix to use for the repository name where the migration data is stored. This is
	// useful when sharing the same database instance with multiple services that own their unique repositories (tables
	// in SQL, collections in MongoDB)
	RepoIdSuffix string

	// ScriptExecutor is a custom script executor to run the scripts for a migration. If not provided uses `Db.Run`
	//
	// This custom executor is added for scenarios where Db.Run would not work as expected, for example using a MongoDB
	// client with AWS DocumentDB where `eval` is not supported, a custom executor using the `mongo` shell CLI could be
	// implemented instead.
	ScriptExecutor dbx.ScriptExecutor
}

// Migrate begins a DB migration process by migrating the scripts located in the provided data dir and storing the
// migration track in a migration repository ('_migration' by default)
//
//    {dataDir}              - The location where the migration scripts are contained
//    {db}                   - The database client already initialized
//    {failOnOrderMismatch}  - Indicates whether the migration should fail if the order of previously migrated scripts
//                             have failed (removing or adding scripts between previously migrated scripts)
//    {repoIdSuffix}         - (optional) A suffix to use for the repository name where the migration data is stored.
//                             This is useful when sharing the same database instance with multiple services that own
//                             their unique repositories (tables in SQL, collections in MongoDB)
//
func Migrate(dataDir string, db dbx.IDatabase, failOnOrderMismatch bool, repoIdSuffix ...string) error {
	migrator := &Migrator{
		Db:                  db,
		DataDir:             dataDir,
		FailOnOrderMismatch: failOnOrderMismatch,
	}
	if len(repoIdSuffix) > 0 {
		migrator.RepoIdSuffix = repoIdSuffix[0]
	}
	return migrator.Migrate()
}

// Migrate begins a DB migration process by migrating the scripts with the configuration contained my the *Migrator
// instance.
func (m *Migrator) Migrate() error {
	var infos []*MigrationInfo
	migrationRepo := MigrationRepo
	if m.RepoIdSuffix != "" {
		migrationRepo += "_" + m.RepoIdSuffix
	}

	if !m.Db.HasRepo(migrationRepo) {
		if err := m.Db.CreateRepo(migrationRepo, &MigrationInfo{}); err != nil {
			return &dbx.DbError{
				Message: fmt.Sprintf("Unable to create the required migration repository. %s", err.Error()),
				Code:    dbx.ErrDbAccess | dbx.ErrDbOperation,
			}
		}
	}

	if err := m.Db.R(migrationRepo).Where(bson.M{}).Sort("filename").All(&infos); err != nil {
		return &dbx.DbError{
			Message: fmt.Sprintf("Unable to read Database info. %s", err.Error()),
			Code:    dbx.ErrDbAccess | dbx.ErrDbOperation,
		}
	}

	objs, err := ioutil.ReadDir(m.DataDir)

	if err != nil {
		return &dbx.DbError{
			Message: fmt.Sprintf("Unable to access scripts path. %s", err.Error()),
			Code:    dbx.ErrFileAccess,
		}
	}

	foundNonMigrated := false

	var toMigrate []MigrationInfo

	for _, f := range streams.From(objs).
		Filter(func(i interface{}) bool {
			return strings.ToLower(i.(os.FileInfo).Name()) != "readme.md"
		}).
		OrderBy(func(a, b interface{}) int {
			return strings.Compare(b.(os.FileInfo).Name(), a.(os.FileInfo).Name())
		}, true).ToArray().([]os.FileInfo) {

		if f.IsDir() {
			continue
		}

		logger.Get().Info("Migrating file ", f.Name())
		fullPath := paths.Combine(m.DataDir, f.Name())
		hash, hashErr := computeHash(fullPath)

		if hashErr != nil {
			return &dbx.DbError{
				Message: fmt.Sprintf("Error computing hash for file '%s', aborting migration.", hashErr.Error()),
				Code:    dbx.ErrMigrationFailed | dbx.ErrFileAccess,
			}
		}

		// Contains the migration info
		if inf := streams.From(infos).
			Filter(
				func(obj interface{}) bool {
					return obj.(*MigrationInfo).ScriptId == f.Name()
				}).
			First(); inf != nil {

			if foundNonMigrated && m.FailOnOrderMismatch {
				return &dbx.DbError{
					Message: fmt.Sprintf("Non-Migrated file found before '%s' which has been migrated. Order import failed, unable to proceed.", f.Name()),
					Code:    dbx.ErrMigrationFailed,
				}
			}

			info := inf.(*MigrationInfo)

			if info.Hash != hash {
				return &dbx.DbError{
					Message: fmt.Sprintf("File '%s' was previously migrated but hashes don't match.", f.Name()),
					Code:    dbx.ErrMigrationFailed,
				}
			} else {
				logger.Get().Info(fmt.Sprintf("File '%s' previously migrated, continuing", f.Name()))
			}

			continue

		} else {
			foundNonMigrated = true
			toMigrate = append(toMigrate, MigrationInfo{
				ScriptId: f.Name(),
				Hash:     hash,
			})
		}
	}

	executor := m.Db.Run
	if m.ScriptExecutor != nil {
		executor = m.ScriptExecutor
	}

	for _, info := range toMigrate {
		fullPath := paths.Combine(m.DataDir, info.ScriptId)
		if content, err := ioutil.ReadFile(fullPath); err != nil {
			return &dbx.DbError{
				Message: fmt.Sprintf("Unable to read data file '%s': %s", info.ScriptId, err.Error()),
				Code:    dbx.ErrFileAccess,
			}
		} else {

			script := string(content)

			if err := executor(script); err != nil {
				return &dbx.DbError{
					Message: fmt.Sprintf("Unable to run command '%s'. %s", info.ScriptId, err.Error()),
					Code:    dbx.ErrDbOperation,
				}
			}

			info.Timestamp = time.Now()

			if err := m.Db.R(migrationRepo).Insert(info); err != nil {
				return &dbx.DbError{
					Message: fmt.Sprintf("Unable to save migration info for '%s'", info.ScriptId),
					Code:    dbx.ErrDbAccess | dbx.ErrDbOperation,
				}
			}
		}
	}

	return nil
}

func computeHash(filePath string) (string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return "", err
	}

	defer file.Close()
	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashInBytes), nil
}
