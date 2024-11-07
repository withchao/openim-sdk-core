//go:build js

package storage

import (
	"database/sql"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/storage/wasm_pglite"
	"github.com/openimsdk/tools/errs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const wasmPGLite = "wasm_pglite"

func init() {
	sql.Register(wasmPGLite, &wasm_pglite.DriverContext{})
}

func OpenGorm(userID string, _ string, log logger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:        fmt.Sprintf("OpenIM_%s_%s.db", constant.BigVersion, userID),
		DriverName: wasmPGLite,
	}), &gorm.Config{Logger: log})
	if err != nil {
		return nil, errs.WrapMsg(err, "open db failed")
	}
	return db, nil
}