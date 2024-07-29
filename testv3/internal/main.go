package internal

import (
	"context"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/openim-sdk-core/v3/testv3/internal/listener"
	"github.com/openimsdk/openim-sdk-core/v3/testv3/internal/message"
	"github.com/openimsdk/openim-sdk-core/v3/version"
	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/tools/log"
)

func InitLog(dirPath string) error {
	return log.InitFromConfig("open-im-sdk-core", "", 6, true, false, dirPath, 0, 24, version.Version)
}

type Config struct {
	WsAddr     string
	ApiAddr    string
	Secret     string
	DataDir    string
	PlatformID int32
}

func GetUserToken(ctx context.Context, req *auth.UserTokenReq) (string, error) {
	resp, err := util.CallApi[auth.UserTokenResp](ctx, constant.GetUsersToken, req)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func GetGroupMemberUserIDs(ctx context.Context, groupID string) ([]string, error) {
	resp, err := util.CallApi[group.GetFullGroupMemberUserIDsResp](ctx, constant.GetFullGroupMemberUserIDs, &group.GetFullGroupMemberUserIDsReq{
		GroupID: groupID,
	})
	if err != nil {
		return nil, err
	}
	return resp.UserIDs, nil
}

func NewUser(userID string, config *Config, msg chan<- message.Message) (*open_im_sdk.LoginMgr, error) {
	userForSDK := open_im_sdk.NewLoginMgr()
	conf := sdk_struct.IMConfig{
		ApiAddr:              config.ApiAddr,
		WsAddr:               config.WsAddr,
		PlatformID:           config.PlatformID,
		DataDir:              config.DataDir,
		LogLevel:             6,
		IsLogStandardOutput:  true,
		IsExternalExtensions: true,
	}
	userForSDK.InitSDK(conf, &listener.ConnListener{UserID: userID})
	ctx := ccontext.WithOperationID(userForSDK.BaseCtx(), utils.OperationIDGenerator())
	token, err := GetUserToken(ctx, &auth.UserTokenReq{UserID: userID, Secret: config.Secret, PlatformID: config.PlatformID})
	if err != nil {
		return nil, err
	}
	if err := userForSDK.Login(ctx, userID, token); err != nil {
		return nil, err
	}
	userForSDK.SetConversationListener(&listener.ConversationListener{Ctx: ctx, UserID: userID})
	userForSDK.SetGroupListener(&listener.GroupListener{Ctx: ctx, UserID: userID})
	userForSDK.SetAdvancedMsgListener(&listener.AdvancedMsgListener{Ctx: ctx, UserID: userID, Ch: msg})
	userForSDK.SetFriendListener(&listener.FriendListener{Ctx: ctx, UserID: userID})
	userForSDK.SetUserListener(&listener.UserListener{Ctx: ctx, UserID: userID})
	if err := userForSDK.User().SyncLoginUserInfo(ctx); err != nil {
		return nil, err
	}
	//if err := userForSDK.Friend().SyncAllFriendList(ctx); err != nil {
	//	return nil, err
	//}
	//if err := userForSDK.Friend().SyncAllBlackList(ctx); err != nil {
	//	return nil, err
	//}
	//if err := userForSDK.Group().SyncAllJoinedGroupsAndMembers(ctx); err != nil {
	//	return nil, err
	//}
	//if err := userForSDK.Conversation().SyncAllConversationHashReadSeqs(ctx); err != nil {
	//	return nil, err
	//}
	fmt.Println("user success", userID)
	return userForSDK, nil
}
