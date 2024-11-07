package db

import (
	"context"
	"errors"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/storage"
	"sync"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/version"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

type TableChecker struct {
	tableCache map[string]bool
	mu         sync.RWMutex
}

func NewTableChecker(tables []string) *TableChecker {
	tc := &TableChecker{
		tableCache: make(map[string]bool),
	}
	tc.InitTableCache(tables)
	return tc
}

func (tc *TableChecker) InitTableCache(tables []string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	for _, table := range tables {
		tc.tableCache[table] = true
	}
}

func (tc *TableChecker) HasTable(tableName string) bool {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.tableCache[tableName]
}
func (tc *TableChecker) UpdateTable(tableName string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.tableCache[tableName] = true
}

type DataBase struct {
	loginUserID  string
	dbDir        string
	conn         *gorm.DB
	tableChecker *TableChecker
	mRWMutex     sync.RWMutex
}

func (d *DataBase) GetMultipleMessageReactionExtension(ctx context.Context, msgIDList []string) (result []*model_struct.LocalChatLogReactionExtensions, err error) {
	// TODO implement me
	panic("implement me")
}

func (d *DataBase) InitSuperLocalErrChatLog(ctx context.Context, groupID string) {
	panic("implement me")
}

func (d *DataBase) InitSuperLocalChatLog(ctx context.Context, groupID string) {
	panic("implement me")
}

func (d *DataBase) SetChatLogFailedStatus(ctx context.Context) {
	panic("implement me")
}

func (d *DataBase) InitDB(ctx context.Context, userID string, dataDir string) error {
	panic("implement me")
}

func (d *DataBase) Close(ctx context.Context) error {
	dbConn, err := d.conn.WithContext(ctx).DB()
	if err != nil {
		return err
	} else {
		if dbConn != nil {
			err := dbConn.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewDataBase(ctx context.Context, loginUserID string, dbDir string, logLevel int) (*DataBase, error) {
	dataBase := &DataBase{loginUserID: loginUserID, dbDir: dbDir}
	err := dataBase.initDB(ctx, logLevel)
	if err != nil {
		return dataBase, errs.WrapMsg(err, "initDB failed "+dbDir)
	}
	tables, err := dataBase.GetExistTables(ctx)
	if err != nil {
		return dataBase, errs.Wrap(err)
	}
	dataBase.tableChecker = NewTableChecker(tables)

	return dataBase, nil
}

func (d *DataBase) initDB(ctx context.Context, logLevel int) error {
	var zLogLevel logger.LogLevel
	if d.loginUserID == "" {
		return errors.New("no uid")
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	if logLevel > 5 {
		zLogLevel = logger.Info
	} else {
		zLogLevel = logger.Silent
	}
	db, err := storage.OpenGorm(d.loginUserID, d.dbDir, log.NewSqlLogger(zLogLevel, false, time.Millisecond*200))
	if err != nil {
		return err
	}

	log.ZDebug(ctx, "open db success")
	sqlDB, err := db.DB()
	if err != nil {
		return errs.WrapMsg(err, "get sql db failed")
	}

	sqlDB.SetConnMaxLifetime(time.Hour * 1)
	sqlDB.SetMaxOpenConns(3)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)
	d.conn = db

	if err = db.AutoMigrate(&model_struct.LocalAppSDKVersion{}); err != nil {
		log.ZWarn(ctx, "AutoMigrate LocalAppSDKVersion failed", err)
	}

	if err = d.versionDataMigrate(ctx); err != nil {
		return err
	}

	return nil
}

func (d *DataBase) versionDataMigrate(ctx context.Context) error {
	verModel, err := d.GetAppSDKVersion(ctx)
	if errs.Unwrap(err) == errs.ErrRecordNotFound {
		err = d.conn.AutoMigrate(
			&model_struct.LocalAppSDKVersion{},
			&model_struct.LocalFriend{},
			&model_struct.LocalFriendRequest{},
			&model_struct.LocalGroup{},
			&model_struct.LocalGroupMember{},
			&model_struct.LocalGroupRequest{},
			&model_struct.LocalErrChatLog{},
			&model_struct.LocalUser{},
			&model_struct.LocalBlack{},
			&model_struct.LocalConversation{},
			&model_struct.NotificationSeqs{},
			&model_struct.LocalChatLog{},
			&model_struct.LocalAdminGroupRequest{},
			&model_struct.LocalWorkMomentsNotification{},
			&model_struct.LocalWorkMomentsNotificationUnreadCount{},
			&model_struct.TempCacheLocalChatLog{},
			&model_struct.LocalChatLogReactionExtensions{},
			&model_struct.LocalUpload{},
			&model_struct.LocalStranger{},
			&model_struct.LocalSendingMessages{},
			&model_struct.LocalUserCommand{},
			&model_struct.LocalVersionSync{},
		)
		if err != nil {
			return err
		}
		err = d.SetAppSDKVersion(ctx, &model_struct.LocalAppSDKVersion{Version: version.Version})
		if err != nil {
			return err
		}

		return nil
	} else if err != nil {
		return err
	}
	if verModel.Version != version.Version {
		switch version.Version {
		case "3.8.0":
			d.conn.AutoMigrate(&model_struct.LocalAppSDKVersion{})
		}
		err = d.SetAppSDKVersion(ctx, &model_struct.LocalAppSDKVersion{Version: version.Version})
		if err != nil {
			return err
		}
	}

	return nil
}
