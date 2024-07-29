package listener

import "context"

type FriendListener struct {
	Ctx    context.Context
	UserID string
}

func (l *FriendListener) OnFriendApplicationAdded(friendApplication string) {

}

func (l *FriendListener) OnFriendApplicationDeleted(friendApplication string) {

}

func (l *FriendListener) OnFriendApplicationAccepted(friendApplication string) {

}

func (l *FriendListener) OnFriendApplicationRejected(friendApplication string) {

}

func (l *FriendListener) OnFriendAdded(friendInfo string) {

}

func (l *FriendListener) OnFriendDeleted(friendInfo string) {

}

func (l *FriendListener) OnFriendInfoChanged(friendInfo string) {

}

func (l *FriendListener) OnBlackAdded(blackInfo string) {

}

func (l *FriendListener) OnBlackDeleted(blackInfo string) {

}
