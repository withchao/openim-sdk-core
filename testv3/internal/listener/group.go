package listener

import "context"

type GroupListener struct {
	Ctx    context.Context
	UserID string
}

func (l *GroupListener) OnJoinedGroupAdded(groupInfo string) {

}

func (l *GroupListener) OnJoinedGroupDeleted(groupInfo string) {

}

func (l *GroupListener) OnGroupMemberAdded(groupMemberInfo string) {

}

func (l *GroupListener) OnGroupMemberDeleted(groupMemberInfo string) {

}

func (l *GroupListener) OnGroupApplicationAdded(groupApplication string) {

}

func (l *GroupListener) OnGroupApplicationDeleted(groupApplication string) {

}

func (l *GroupListener) OnGroupInfoChanged(groupInfo string) {

}

func (l *GroupListener) OnGroupDismissed(groupInfo string) {

}

func (l *GroupListener) OnGroupMemberInfoChanged(groupMemberInfo string) {

}

func (l *GroupListener) OnGroupApplicationAccepted(groupApplication string) {

}

func (l *GroupListener) OnGroupApplicationRejected(groupApplication string) {

}
