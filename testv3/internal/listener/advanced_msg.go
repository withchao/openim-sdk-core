package listener

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/testv3/internal/message"
)

type AdvancedMsgListener struct {
	Ctx    context.Context
	UserID string
	Ch     chan<- message.Message
}

func (l *AdvancedMsgListener) OnRecvNewMessage(msg string) {
	l.Ch <- message.Message{UserID: l.UserID, Data: msg}
	//log.ZDebug(l.Ctx, "OnRecvNewMessage", "message", message)
}

func (l *AdvancedMsgListener) OnRecvC2CReadReceipt(msgReceiptList string) {

}

func (l *AdvancedMsgListener) OnRecvGroupReadReceipt(groupMsgReceiptList string) {

}

func (l *AdvancedMsgListener) OnNewRecvMessageRevoked(messageRevoked string) {

}

func (l *AdvancedMsgListener) OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string) {

}

func (l *AdvancedMsgListener) OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string) {

}

func (l *AdvancedMsgListener) OnRecvMessageExtensionsAdded(msgID string, reactionExtensionList string) {

}

func (l *AdvancedMsgListener) OnRecvOfflineNewMessage(message string) {

}

func (l *AdvancedMsgListener) OnMsgDeleted(message string) {

}

func (l *AdvancedMsgListener) OnRecvOnlineOnlyMessage(message string) {

}
