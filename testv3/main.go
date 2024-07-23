package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/openim-sdk-core/v3/testv3/internal"
	"github.com/openimsdk/openim-sdk-core/v3/testv3/internal/message"
	"github.com/openimsdk/protocol/sdkws"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

const delimiter = ";;;"

func main() {
	Main()
}

func Main() {
	var ip string

	ip = "172.16.8.48"
	//ip = "14.29.213.197"

	conf := internal.Config{
		WsAddr:  fmt.Sprintf("ws://%s:10001", ip),
		ApiAddr: fmt.Sprintf("http://%s:10002", ip),
		//WsAddr:     fmt.Sprintf("ws://%s:50001", ip),
		//ApiAddr:    fmt.Sprintf("http://%s:50002", ip),
		Secret:     "openIM123",
		DataDir:    "./testv3/dbdata",
		PlatformID: constant.AdminPlatformID,
	}

	var (
		userID  string
		groupID string
	)

	userID = "2110910952" // vm
	//userID = "2890713225" // 14

	groupID = "2783684806" // vm
	//groupID = "2783684806" // 14

	_, _ = userID, groupID
	var (
		users   []*open_im_sdk.LoginMgr
		userIDs []string
	)
	ch := make(chan message.Message, 1024*8)
	user, err := internal.NewUser(userID, &conf, ch)
	if err != nil {
		panic(err)
	}
	ctx := ccontext.WithOperationID(user.Context(), "sasdkaskd")
	//userIDs, err = internal.GetGroupMemberUserIDs(ctx, groupID)
	//if err != nil {
	//	panic(err)
	//}
	userIDs = []string{"1814217053"} // vm
	//userIDs = []string{"1112594574"} // 14
	fmt.Println("userIDs", userIDs)

	msgPrefix := strconv.FormatUint(rand.Uint64(), 10)
	users = make([]*open_im_sdk.LoginMgr, 0, len(userIDs)+1)
	for _, uID := range userIDs {
		if uID == userID {
			continue
		}
		u, err := internal.NewUser(uID, &conf, ch)
		if err != nil {
			panic(err)
		}
		users = append(users, u)
	}

	var sendMsgNum int64

	go sendUserMsg(ctx, user, userIDs[0], msgPrefix, &sendMsgNum)
	//go sendGroupMsg(ctx, user, groupID, msgPrefix, &sendMsgNum)

	signalCh := make(chan os.Signal, 1)
	go func() {
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	}()

	var (
		start     = time.Now()
		total     time.Duration
		recvCount int64
	)

	defer func() {
		if recvCount == 0 {
			return
		}
		end := time.Now()
		runTime := end.Sub(start)
		fmt.Println("user", len(users)+1)
		fmt.Println("send", atomic.LoadInt64(&sendMsgNum))
		fmt.Println("recv", recvCount)
		fmt.Println("runtime", runTime.String())
		fmt.Println("avg", (total / time.Duration(recvCount)).String())
	}()
	for {
		select {
		case val := <-signalCh:
			fmt.Println("exit", val.String())
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			var data sdk_struct.MsgStruct
			if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
				panic(err)
			}
			sendTime, ok := parseTestMsg(msgPrefix, &data)
			if !ok {
				continue
			}
			duration := time.Now().Sub(sendTime)
			recvCount++
			total += duration
			fmt.Printf("[%s] %s -> %s time %s\n", data.ClientMsgID, user.GetLoginUserID(), msg.UserID, duration.String())
		}
	}
}

func sendGroupMsg(ctx context.Context, user *open_im_sdk.LoginMgr, groupID string, msgPrefix string, sendMsgNum *int64) {
	for i := 0; i < 100000; i++ {
		atomic.AddInt64(sendMsgNum, 1)
		msg, err := user.Conversation().CreateTextMessage(ctx, buildTestMsg(msgPrefix, time.Now()))
		if err != nil {
			panic(err)
		}
		resp, err := user.Conversation().SendMessage(ctx, msg, "", groupID, &sdkws.OfflinePushInfo{}, false)
		if err != nil {
			panic(err)
		}
		_ = resp
		//fmt.Println("======>", resp)
		//fmt.Println("send success", i+1)
		time.Sleep(time.Second)
	}
}

func sendUserMsg(ctx context.Context, user *open_im_sdk.LoginMgr, userID string, msgPrefix string, sendMsgNum *int64) {
	for i := 0; i < 100000; i++ {
		atomic.AddInt64(sendMsgNum, 1)
		msg, err := user.Conversation().CreateTextMessage(ctx, buildTestMsg(msgPrefix, time.Now()))
		if err != nil {
			panic(err)
		}
		resp, err := user.Conversation().SendMessage(ctx, msg, userID, "", &sdkws.OfflinePushInfo{}, false)
		if err != nil {
			panic(err)
		}
		_ = resp
		//fmt.Println("======>", resp)
		//fmt.Println("send success", i+1)
		time.Sleep(time.Second)
	}
}

func buildTestMsg(prefix string, now time.Time) string {
	return strings.Join([]string{prefix, strconv.FormatInt(now.UnixMilli(), 10)}, delimiter)
}

func parseTestMsg(prefix string, data *sdk_struct.MsgStruct) (time.Time, bool) {
	if data.TextElem == nil {
		return time.Time{}, false
	}
	arr := strings.Split(data.TextElem.Content, delimiter)
	if len(arr) != 2 {
		return time.Time{}, false
	}
	if arr[0] != prefix {
		return time.Time{}, false
	}
	val, err := strconv.ParseInt(arr[1], 10, 64)
	if err != nil {
		return time.Time{}, false
	}
	return time.UnixMilli(val), true
}
