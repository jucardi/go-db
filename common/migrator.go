package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jucardi/go-db"
	"github.com/jucardi/go-db/logger"
	"gopkg.in/jucardi/go-osx.v1/paths"
	"gopkg.in/jucardi/go-streams.v1/streams"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	MigrationRepo = "_migration"
)

func Migrate(dataDir string, db dbx.IDatabase, failOnOrderMismatch bool) *dbx.DbError {
	var infos []*MigrationInfo

	if !db.HasRepo(MigrationRepo) {
		if err := db.CreateRepo(MigrationRepo, &MigrationInfo{}); err != nil {
			return &dbx.DbError{
				Message: fmt.Sprintf("Unable to create the required migration repository. %s", err.Error()),
				Code:    dbx.ErrDbAccess | dbx.ErrDbOperation,
			}
		}
	}

	if err := db.R(MigrationRepo).Where(bson.M{}).Sort("filename").All(&infos); err != nil {
		return &dbx.DbError{
			Message: fmt.Sprintf("Unable to read Database info. %s", err.Error()),
			Code:    dbx.ErrDbAccess | dbx.ErrDbOperation,
		}
	}

	objs, err := ioutil.ReadDir(dataDir)

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

		logger.Get().Infof("Migrating file '%s'", f.Name())
		fullPath := paths.Combine(dataDir, f.Name())
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

			if foundNonMigrated && failOnOrderMismatch {
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
				logger.Get().Infof("File '%s' previously migrated, continuing", f.Name())
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

	for _, info := range toMigrate {
		fullPath := paths.Combine(dataDir, info.ScriptId)
		if content, err := ioutil.ReadFile(fullPath); err != nil {
			return &dbx.DbError{
				Message: fmt.Sprintf("Unable to read data file '%s': %s", info.ScriptId, err.Error()),
				Code:    dbx.ErrFileAccess,
			}
		} else {

			script := string(content)

			if err := db.Run(script); err != nil {
				return &dbx.DbError{
					Message: fmt.Sprintf("Unable to run command '%s'. %s", info.ScriptId, err.Error()),
					Code:    dbx.ErrDbOperation,
				}
			}

			info.Timestamp = time.Now()

			if err := db.R(MigrationRepo).Insert(info); err != nil {
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
