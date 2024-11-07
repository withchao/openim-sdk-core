// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) BatchInsertTempCacheMessageList(ctx context.Context, MessageList []*model_struct.TempCacheLocalChatLog) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	if MessageList == nil {
		return nil
	}
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(MessageList).Error, "BatchInsertTempCacheMessageList failed")
}
func (d *DataBase) InsertTempCacheMessage(ctx context.Context, Message *model_struct.TempCacheLocalChatLog) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(Message).Error, "InsertTempCacheMessage failed")

}
