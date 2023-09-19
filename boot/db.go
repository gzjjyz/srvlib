package boot

import (
	"errors"
	"github.com/gzjjyz/logger"

	"github.com/gzjjyz/micro/env"
	"github.com/gzjjyz/srvlib/db"
)

const (
	defaultMysqlConnName = "default"
)

func InitOrmMysql(connName string) error {
	if connName == "" {
		connName = defaultMysqlConnName
	}
	connCfg, ok := env.MustMeta().DBConnections.GetMysqlConn(connName)
	if !ok {
		err := errors.New("mysql connection config not found")
		logger.Errorf(err.Error())
		return err
	}

	if err := db.InitOrmMysqlV2(connCfg.User, connCfg.Password, connCfg.Host, connCfg.Port, connCfg.Databases); err != nil {
		logger.Errorf(err.Error())
		return err
	}

	return nil
}
