package main

import (
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/testv3/internal"
	"github.com/openimsdk/protocol/sdkws"
	"time"
)

func main() {
	Main()
}

func Main() {
	var ip string

	ip = "172.16.8.48"

	conf := internal.Config{
		WsAddr:     fmt.Sprintf("ws://%s:10001", ip),
		ApiAddr:    fmt.Sprintf("http://%s:10002", ip),
		Secret:     "openIM123",
		DataDir:    "./testv3/dbdata",
		PlatformID: constant.AdminPlatformID,
	}

	var (
		userID  string
		groupID string
	)

	userID = "2110910952"
	groupID = "1830477527"

	var (
		users   []*open_im_sdk.LoginMgr
		userIDs []string
	)
	user, err := internal.NewUser(userID, &conf)
	if err != nil {
		panic(err)
	}
	ctx := ccontext.WithOperationID(user.Context(), "sasdkaskd")
	userIDs, err = internal.GetGroupMemberUserIDs(ctx, groupID)
	if err != nil {
		panic(err)
	}
	users = make([]*open_im_sdk.LoginMgr, 0, len(userIDs)+1)
	//users = append(users, user)
	for _, uID := range userIDs {
		if uID == userID {
			continue
		}
		u, err := internal.NewUser(uID, &conf)
		if err != nil {
			panic(err)
		}
		users = append(users, u)
	}

	for i := 0; i < 10000; i++ {
		msg, err := user.Conversation().CreateTextMessage(ctx, fmt.Sprintf("hello_%d", i))
		if err != nil {
			panic(err)
		}
		resp, err := user.Conversation().SendMessage(ctx, msg, "", groupID, &sdkws.OfflinePushInfo{}, false)
		if err != nil {
			panic(err)
		}
		fmt.Println("======>", resp)
		time.Sleep(time.Second)
	}

	time.Sleep(time.Hour)
}
