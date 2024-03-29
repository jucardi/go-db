package mgo

import (
	"fmt"
	"github.com/jucardi/go-db"
	"github.com/jucardi/go-db/logger"
	"github.com/jucardi/go-strings/stringx"
	"time"
	"gopkg.in/mgo.v2"
)

var (
	// ErrNotFound is the error returned when no results are found in a mongo operation.
	ErrNotFound = mgo.ErrNotFound

	// ErrCursor is the error returned when the cursor used in a mongo operation is not valid.
	ErrCursor = mgo.ErrCursor
)

type provider struct {
}

func (p *provider) Dial(cfg *dbx.DbConfig) (dbx.IDatabase, error) {
	s, err := Dial(cfg)
	if err != nil {
		return nil, err
	}
	return s.DB(""), nil
}

func init() {
	if err := dbx.Register("MongoDB", &provider{}, true); err != nil {
		logger.Get().Error("Unable to register MongoDB provider, ", err.Error())
	}
}

// Dial establishes a new session to the cluster identified by the given seed
// server(s). The session will enable communication with all of the servers in
// the cluster, so the seed servers are used only to find out about the cluster
// topology.
func Dial(cfg *dbx.DbConfig) (ISession, error) {
	url := toUrl(cfg)
	s, err := mgo.Dial(url)

	for i := 1; err != nil && i <= cfg.DialMaxRetries; i++ {
		logger.Get().Error(fmt.Sprintf("Unable to connect to mongo on '%s': %v. Retrying in %v", cfg.Host, err, cfg.DialRetryTimeout))
		time.Sleep(time.Duration(cfg.DialRetryTimeout)*time.Millisecond)
		logger.Get().Warn(fmt.Sprintf("Retrying to connect to mongo, attempt %d of %d", i, cfg.DialMaxRetries))
		s, err = mgo.Dial(url)
	}

	return fromSession(s), err
}

// DialWithTimeout works like Dial, but uses timeout as the amount of time to
// wait for a server to respond when first connecting and also on follow up
// operations in the session. If timeout is zero, the call may block
// forever waiting for a connection to be made.
//
// See SetSyncTimeout for customizing the timeout for the session.
func DialWithTimeout(url string, timeout time.Duration) (ISession, error) {
	s, err := mgo.DialWithTimeout(url, timeout)
	return fromSession(s), err
}

// DialWithInfo establishes a new session to the cluster identified by info.
func DialWithInfo(info *mgo.DialInfo) (ISession, error) {
	s, err := mgo.DialWithInfo(info)
	return fromSession(s), err
}

// IsDup returns whether err informs of a duplicate key error because
// a primary key index or a secondary unique index already has an entry
// with the given value.
func IsDup(err error) bool {
	return mgo.IsDup(err)
}

func toUrl(cfg *dbx.DbConfig) string {
	builder := stringx.Builder().Append("mongodb://")

	if cfg.Username != "" {
		builder.Append(cfg.Username)

		if cfg.Password != "" {
			builder.Append(":").Append(cfg.Password)
		}
		builder.Append("@")
	}

	builder.Append(cfg.Host)
	if cfg.Port > 0 {
		builder.Appendf(":%d", cfg.Port)
	}
	if cfg.Database != "" {
		builder.Append("/").Append(cfg.Database)
	}
	return builder.Build()
}
