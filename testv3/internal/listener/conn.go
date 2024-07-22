package listener

import "context"

type ConnListener struct {
	Ctx    context.Context
	UserID string
}

func (l *ConnListener) OnConnecting() {

}

func (l *ConnListener) OnConnectSuccess() {

}

func (l *ConnListener) OnConnectFailed(errCode int32, errMsg string) {

}

func (l *ConnListener) OnKickedOffline() {

}

func (l *ConnListener) OnUserTokenExpired() {

}

func (l *ConnListener) OnUserTokenInvalid(errMsg string) {

}
