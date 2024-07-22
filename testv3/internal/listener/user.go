package listener

import "context"

type UserListener struct {
	Ctx    context.Context
	UserID string
}

func (l *UserListener) OnSelfInfoUpdated(userInfo string) {

}

func (l *UserListener) OnUserStatusChanged(userOnlineStatus string) {

}

func (l *UserListener) OnUserCommandAdd(userCommand string) {

}

func (l *UserListener) OnUserCommandDelete(userCommand string) {

}

func (l *UserListener) OnUserCommandUpdate(userCommand string) {

}
