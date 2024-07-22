package listener

import "context"

type ConversationListener struct {
	Ctx    context.Context
	UserID string
}

func (l *ConversationListener) OnSyncServerStart(reinstalled bool) {

}

func (l *ConversationListener) OnSyncServerFinish(reinstalled bool) {

}

func (l *ConversationListener) OnSyncServerProgress(progress int) {

}

func (l *ConversationListener) OnSyncServerFailed(reinstalled bool) {

}

func (l *ConversationListener) OnNewConversation(conversationList string) {

}

func (l *ConversationListener) OnConversationChanged(conversationList string) {

}

func (l *ConversationListener) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {

}

func (l *ConversationListener) OnConversationUserInputStatusChanged(change string) {

}
