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
	log = nil
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:        fmt.Sprintf("OpenIM_%s_%s.db", constant.BigVersion, userID),
		DriverName: wasmPGLite,
	}), &gorm.Config{
		//Logger:                 log,
		SkipDefaultTransaction: true,
		CreateBatchSize:        1,
	})
	if err != nil {
		return nil, errs.WrapMsg(err, "open db failed")
	}
	db = db.Debug()
	//conn, err := db.DB()
	//if err != nil {
	//	return nil, err
	//}
	//conn.SetMaxOpenConns(1)
	//conn.SetMaxIdleConns(1)

	//if err := db.AutoMigrate(&UserTest{}); err != nil {
	//	return nil, err
	//}
	//if err := db.Create(&UserTest{ID: strconv.Itoa(int(time.Now().Unix())), Name: "test", Age: 10, Text: "hello world!", CreateTime: time.Now()}).Error; err != nil {
	//	return nil, err
	//}
	return db, nil
}

//type UserTest struct {
//	ID         string    `gorm:"column:id"`
//	Name       string    `gorm:"column:name;type:varchar(255)"`
//	Age        int       `gorm:"column:age"`
//	Text       string    `gorm:"column:text"`
//	CreateTime time.Time `gorm:"column:create_time"`
//}
