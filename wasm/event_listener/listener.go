package event_listener

import (
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"syscall/js"
)

type ConnCallback struct {
	uid string
	CallbackWriter
}

func NewConnCallback(funcName string, callback *js.Value) *ConnCallback {
	return &ConnCallback{CallbackWriter: NewEventData(callback).SetEvent(funcName)}
}

func (i *ConnCallback) OnConnecting() {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()
}

func (i *ConnCallback) OnConnectSuccess() {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()

}
func (i *ConnCallback) OnConnectFailed(errCode int32, errMsg string) {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetErrCode(errCode).SetErrMsg(errMsg).SendMessage()
}

func (i *ConnCallback) OnKickedOffline() {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()
}

func (i *ConnCallback) OnUserTokenExpired() {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()
}

func (i *ConnCallback) OnSelfInfoUpdated(userInfo string) {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(userInfo).SendMessage()
}

type ConversationCallback struct {
	uid string
	CallbackWriter
}

func NewConversationCallback(callback *js.Value) *ConversationCallback {
	return &ConversationCallback{CallbackWriter: NewEventData(callback)}
}
func (c ConversationCallback) OnSyncServerStart() {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()
}

func (c ConversationCallback) OnSyncServerFinish() {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()
}

func (c ConversationCallback) OnSyncServerFailed() {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()

}

func (c ConversationCallback) OnNewConversation(conversationList string) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(conversationList).SendMessage()

}

func (c ConversationCallback) OnConversationChanged(conversationList string) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(conversationList).SendMessage()

}

func (c ConversationCallback) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(totalUnreadCount).SendMessage()
}

type AdvancedMsgCallback struct {
	CallbackWriter
}

func NewAdvancedMsgCallback(callback *js.Value) *AdvancedMsgCallback {
	return &AdvancedMsgCallback{CallbackWriter: NewEventData(callback)}
}
func (a AdvancedMsgCallback) OnRecvNewMessage(message string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(message).SendMessage()
}

func (a AdvancedMsgCallback) OnRecvC2CReadReceipt(msgReceiptList string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(msgReceiptList).SendMessage()
}

func (a AdvancedMsgCallback) OnRecvGroupReadReceipt(groupMsgReceiptList string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(groupMsgReceiptList).SendMessage()
}

func (a AdvancedMsgCallback) OnRecvMessageRevoked(msgID string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(msgID).SendMessage()
}

func (a AdvancedMsgCallback) OnNewRecvMessageRevoked(messageRevoked string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(messageRevoked).SendMessage()
}
func (a AdvancedMsgCallback) OnRecvMessageModified(message string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(message).SendMessage()
}
func (a AdvancedMsgCallback) OnRecvMessageExtensionsChanged(clientMsgID string, reactionExtensionList string) {
	m := make(map[string]interface{})
	m["clientMsgID"] = clientMsgID
	m["reactionExtensionList"] = reactionExtensionList
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(m)).SendMessage()
}

func (a AdvancedMsgCallback) OnRecvMessageExtensionsDeleted(clientMsgID string, reactionExtensionKeyList string) {
	m := make(map[string]interface{})
	m["clientMsgID"] = clientMsgID
	m["reactionExtensionKeyList"] = reactionExtensionKeyList
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(m)).SendMessage()
}

type BaseCallback struct {
	CallbackWriter
}

func (b *BaseCallback) EventData() CallbackWriter {
	return b.CallbackWriter
}

func NewBaseCallback(funcName string, _ *js.Value) *BaseCallback {
	return &BaseCallback{CallbackWriter: NewPromiseHandler().SetEvent(funcName)}
}

func (b *BaseCallback) OnError(errCode int32, errMsg string) {
	b.CallbackWriter.SetErrCode(errCode).SetErrMsg(errMsg).SendMessage()
}
func (b *BaseCallback) OnSuccess(data string) {
	b.CallbackWriter.SetData(data).SendMessage()
}

type SendMessageCallback struct {
	BaseCallback
	globalEvent CallbackWriter
	clientMsgID string
}

func (s *SendMessageCallback) SetClientMsgID(args *[]js.Value) *SendMessageCallback {
	m := sdk_struct.MsgStruct{}
	utils.JsonStringToStruct((*args)[1].String(), &m)
	s.clientMsgID = m.ClientMsgID
	return s
}
func NewSendMessageCallback(funcName string, callback *js.Value) *SendMessageCallback {
	return &SendMessageCallback{BaseCallback: BaseCallback{CallbackWriter: NewPromiseHandler().SetEvent(funcName)}, globalEvent: NewEventData(callback).SetEvent(funcName)}
}

func (s *SendMessageCallback) OnProgress(progress int) {
	mReply := make(map[string]interface{})
	mReply["progress"] = progress
	mReply["clientMsgID"] = s.clientMsgID
	s.globalEvent.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(mReply)).SendMessage()
}

type BatchMessageCallback struct {
	CallbackWriter
}

func NewBatchMessageCallback(callback *js.Value) *BatchMessageCallback {
	return &BatchMessageCallback{CallbackWriter: NewEventData(callback)}
}

func (b *BatchMessageCallback) OnRecvNewMessages(messageList string) {
	b.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(messageList).SendMessage()
}

type UserCallback struct {
	CallbackWriter
}

func NewUserCallback(callback *js.Value) *UserCallback {
	return &UserCallback{CallbackWriter: NewEventData(callback)}
}
func (u UserCallback) OnSelfInfoUpdated(userInfo string) {
	u.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(userInfo).SendMessage()
}
