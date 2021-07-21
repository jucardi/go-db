package sql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jucardi/go-db"
	"github.com/jucardi/go-db/logger"
	"github.com/jucardi/go-strings/stringx"
)

type provider struct {
}

func (p *provider) Dial(cfg *dbx.DbConfig) (dbx.IDatabase, error) {
	return Dial(cfg)
}

func init() {
	if err := dbx.Register("MySQL", &provider{}, true); err != nil {
		logger.Get().Error("Unable to register MySQL provider, ", err.Error())
	}
}

// Dial establishes a new session to the cluster identified by the given seed
// server(s). The session will enable communication with all of the servers in
// the cluster, so the seed servers are used only to find out about the cluster
// topology.
func Dial(cfg *dbx.DbConfig) (IDatabase, error) {
	db, err := gorm.Open("mysql", getUrl(cfg))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to mysql, %s", err.Error())
	}
	return FromDB(db, false), nil
}

func getUrl(cfg *dbx.DbConfig) string {
	builder := stringx.Builder()

	if cfg.Username != "" {
		builder.Append(cfg.Username)

		if cfg.Password != "" {
			builder.Append(":").Append(cfg.Password)
		}
		builder.Append("@")
	}

	builder.Append("tpc(", cfg.Host)
	if cfg.Port > 0 {
		builder.Appendf(":%d", cfg.Port)
	}
	builder.Append(")")
	if cfg.Database != "" {
		builder.Append("/").Append(cfg.Database)
	}

	return builder.Append(cfg.Options).Build()
}
