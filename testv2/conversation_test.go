// Copyright © 2023 OpenIM SDK.
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

package testv2

import (
	"context"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/sdk_struct"
	"strings"
	"testing"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

func Test_GetAllConversationList(t *testing.T) {
	conversations, err := open_im_sdk.UserForSDK.Conversation().GetAllConversationList(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, conversation := range conversations {
		t.Log(conversation)
	}
}

func Test_GetConversationListSplit(t *testing.T) {
	conversations, err := open_im_sdk.UserForSDK.Conversation().GetConversationListSplit(ctx, 0, 20)
	if err != nil {
		t.Fatal(err)
	}
	for _, conversation := range conversations {
		t.Log(conversation)
	}
}

func Test_SetConversationRecvMessageOpt(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetConversationRecvMessageOpt(ctx, []string{"asdasd"}, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SetSetGlobalRecvMessageOpt(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetGlobalRecvMessageOpt(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_HideConversation(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().HideConversation(ctx, "asdasd")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GetConversationRecvMessageOpt(t *testing.T) {
	opts, err := open_im_sdk.UserForSDK.Conversation().GetConversationRecvMessageOpt(ctx, []string{"asdasd"})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range opts {
		t.Log(v.ConversationID, *v.Result)
	}
}

func Test_GetGlobalRecvMessageOpt(t *testing.T) {
	opt, err := open_im_sdk.UserForSDK.Conversation().GetOneConversation(ctx, 2, "1772958501")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*opt)
}

func Test_GetGetMultipleConversation(t *testing.T) {
	conversations, err := open_im_sdk.UserForSDK.Conversation().GetMultipleConversation(ctx, []string{"asdasd"})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range conversations {
		t.Log(v)
	}
}

func Test_DeleteConversation(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().DeleteConversation(ctx, "group_17729585012")
	if err != nil {
		if !strings.Contains(err.Error(), "no update") {
			t.Fatal(err)
		}
	}
}

func Test_DeleteAllConversationFromLocal(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().DeleteAllConversationFromLocal(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SetConversationDraft(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetConversationDraft(ctx, "group_17729585012", "draft")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ResetConversationGroupAtType(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().ResetConversationGroupAtType(ctx, "group_17729585012")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_PinConversation(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().PinConversation(ctx, "group_17729585012", true)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SetOneConversationPrivateChat(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetOneConversationPrivateChat(ctx, "single_3411008330", true)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SetOneConversationBurnDuration(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetOneConversationBurnDuration(ctx, "single_3411008330", 10)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SetOneConversationRecvMessageOpt(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetOneConversationRecvMessageOpt(ctx, "single_3411008330", 1)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GetTotalUnreadMsgCount(t *testing.T) {
	count, err := open_im_sdk.UserForSDK.Conversation().GetTotalUnreadMsgCount(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func Test_SendMessage(t *testing.T) {
	ctx = context.WithValue(ctx, "callback", TestSendMsg{})
	msg, _ := open_im_sdk.UserForSDK.Conversation().CreateTextMessage(ctx, "textMsg")
	_, err := open_im_sdk.UserForSDK.Conversation().SendMessage(ctx, msg, "3411008330", "", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SendMessageNotOss(t *testing.T) {
	ctx = context.WithValue(ctx, "callback", TestSendMsg{})
	msg, _ := open_im_sdk.UserForSDK.Conversation().CreateTextMessage(ctx, "textMsg")
	_, err := open_im_sdk.UserForSDK.Conversation().SendMessageNotOss(ctx, msg, "3411008330", "", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SendMessageByBuffer(t *testing.T) {
	ctx = context.WithValue(ctx, "callback", TestSendMsg{})
	msg, _ := open_im_sdk.UserForSDK.Conversation().CreateTextMessage(ctx, "textMsg")
	_, err := open_im_sdk.UserForSDK.Conversation().SendMessageByBuffer(ctx, msg, "3411008330", "", &sdkws.OfflinePushInfo{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_FindMessageList(t *testing.T) {
	msgs, err := open_im_sdk.UserForSDK.Conversation().FindMessageList(ctx, []*sdk_params_callback.ConversationArgs{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(msgs.TotalCount)
	for _, v := range msgs.FindResultItems {
		t.Log(v)
	}
}

func Test_GetHistoryMessageList(t *testing.T) {
	msgs, err := open_im_sdk.UserForSDK.Conversation().GetHistoryMessageList(ctx, sdk_params_callback.GetHistoryMessageListParams{})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range msgs {
		t.Log(v)
	}
}

func Test_GetAdvancedHistoryMessageList(t *testing.T) {
	msgs, err := open_im_sdk.UserForSDK.Conversation().GetAdvancedHistoryMessageList(ctx, sdk_params_callback.GetAdvancedHistoryMessageListParams{})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range msgs.MessageList {
		t.Log(v)
	}
}

func Test_GetAdvancedHistoryMessageListReverse(t *testing.T) {
	msgs, err := open_im_sdk.UserForSDK.Conversation().GetAdvancedHistoryMessageListReverse(ctx, sdk_params_callback.GetAdvancedHistoryMessageListParams{})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range msgs.MessageList {
		t.Log(v)
	}
}

func Test_DeleteMessageFromLocalStorage(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().DeleteMessageFromLocalStorage(ctx, &sdk_struct.MsgStruct{})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ClearC2CHistoryMessage(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().ClearC2CHistoryMessage(ctx, "3411008330")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ClearGroupHistoryMessage(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().ClearGroupHistoryMessage(ctx, "group_17729585012")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ClearC2CHistoryMessageFromLocalAndSvr(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().ClearC2CHistoryMessageFromLocalAndSvr(ctx, "3411008330")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ClearGroupHistoryMessageFromLocalAndSvr(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().ClearGroupHistoryMessageFromLocalAndSvr(ctx, "group_17729585012")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_InsertSingleMessageToLocalStorage(t *testing.T) {
	_, err := open_im_sdk.UserForSDK.Conversation().InsertSingleMessageToLocalStorage(ctx, &sdk_struct.MsgStruct{}, "3411008330", "")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_InsertGroupMessageToLocalStorage(t *testing.T) {
	_, err := open_im_sdk.UserForSDK.Conversation().InsertGroupMessageToLocalStorage(ctx, &sdk_struct.MsgStruct{}, "group_17729585012", "")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SearchLocalMessages(t *testing.T) {
	msgs, err := open_im_sdk.UserForSDK.Conversation().SearchLocalMessages(ctx, &sdk_params_callback.SearchLocalMessagesParams{})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range msgs.SearchResultItems {
		t.Log(v)
	}
}

func Test_DeleteConversationFromLocalAndSvr(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().DeleteConversationFromLocalAndSvr(ctx, "group_17729585012")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_DeleteMessageFromLocalAndSvr(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().DeleteMessageFromLocalAndSvr(ctx, &sdk_struct.MsgStruct{})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_DeleteAllMsgFromLocalAndSvr(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().DeleteAllMsgFromLocalAndSvr(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_DeleteAllMsgFromLocal(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().DeleteAllMsgFromLocal(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
