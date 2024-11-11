//go:build !js

package storage

import (
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/tools/errs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"path/filepath"
)

const DBName = SQLite

func OpenGorm(userID string, dbDir string, log logger.Interface) (*gorm.DB, error) {
	path := filepath.Join(dbDir, fmt.Sprintf("OpenIM_%s_%s.db", constant.BigVersion, userID))
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{Logger: log})
	if err != nil {
		return nil, errs.WrapMsg(err, "open db failed "+path)
	}
	return db, nil
}
