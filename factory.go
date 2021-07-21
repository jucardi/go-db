package dbx

import (
	"errors"
	"github.com/jucardi/go-beans/beans"
	"github.com/jucardi/go-strings/stringx"
)

type IDbProvider interface {
	Dial(cfg *DbConfig) (IDatabase, error)
}

func Register(dbType string, provider IDbProvider, setPrimary ...bool) error {
	if err := beans.Register((*IDbProvider)(nil), dbType, provider); err != nil {
		return err
	}
	if len(setPrimary) > 0 && setPrimary[0] {
		return beans.SetPrimary((*IDbProvider)(nil), dbType)
	}
	return nil
}

func Dial(cfg *DbConfig, dbType ...string) (IDatabase, error) {
	provider := GetProvider(dbType...)
	if provider == nil {
		return nil, errors.New("database provider was not found.")
	}
	return provider.Dial(cfg)
}

func GetProvider(dbType ...string) IDbProvider {
	ret := beans.Resolve((*IDbProvider)(nil), stringx.GetOrDefault("", dbType...))
	if ret != nil {
		return ret.(IDbProvider)
	}
	return nil
}
