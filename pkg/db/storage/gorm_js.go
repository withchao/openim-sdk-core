//go:build js

package storage

import (
	"database/sql"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/storage/wasm_sqlite"
	"github.com/openimsdk/tools/errs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const wasmSQLite = "wasm_sqlite"

func init() {
	sql.Register(wasmSQLite, &wasm_sqlite.DriverContext{})
}

func OpenGorm(userID string, _ string, log logger.Interface) (*gorm.DB, error) {
	log = nil
	db, err := gorm.Open(sqlite.Dialector{
		DSN:        fmt.Sprintf("OpenIM_%s_%s.db", constant.BigVersion, userID),
		DriverName: wasmSQLite,
	}, &gorm.Config{
		Logger:                 log,
		SkipDefaultTransaction: true,
		CreateBatchSize:        100,
	})
	if err != nil {
		return nil, errs.WrapMsg(err, "open db failed")
	}
	db = db.Debug()
	return db, nil
}
