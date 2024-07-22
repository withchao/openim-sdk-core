package listener

import (
	"context"
	"github.com/openimsdk/tools/log"
)

type AdvancedMsgListener struct {
	Ctx    context.Context
	UserID string
}

func (l *AdvancedMsgListener) OnRecvNewMessage(message string) {
	log.ZDebug(l.Ctx, "OnRecvNewMessage", "message", message)
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
