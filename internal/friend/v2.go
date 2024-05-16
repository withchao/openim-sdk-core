package friend

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	sdk "github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/server_api_params"
	pconstant "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/friend"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"strconv"
	"strings"
)

func (f *Friend) CalculateHash(friends []*model_struct.LocalFriend) uint64 {
	datautil.SortAny(friends, func(a, b *model_struct.LocalFriend) bool {
		return a.CreateTime > b.CreateTime
	})
	if len(friends) > pconstant.MaxSyncPullNumber {
		friends = friends[:pconstant.MaxSyncPullNumber]
	}
	hashStr := strings.Join(datautil.Slice(friends, func(f *model_struct.LocalFriend) string {
		return strings.Join([]string{
			f.FriendUserID,
			f.Remark,
			strconv.FormatInt(f.CreateTime, 10),
			strconv.Itoa(int(f.AddSource)),
			f.OperatorUserID,
			f.Ex,
			strconv.FormatBool(f.IsPinned),
		}, ",")
	}), ";")
	sum := md5.Sum([]byte(hashStr))
	return binary.BigEndian.Uint64(sum[:])
}

func (f *Friend) getFriendServerHash(ctx context.Context) (*friend.GetFriendHashResp, error) {
	resp, err := util.CallApi[friend.GetFriendHashResp](ctx, constant.GetFriendHash, &friend.GetFriendHashReq{UserID: f.loginUserID})
	if err != nil {
		return nil, err
	}
	f.friendNum = resp.Total
	return resp, nil
}

func (f *Friend) SyncFriendPart(ctx context.Context) error {
	hashResp, err := f.getFriendServerHash(ctx)
	if err != nil {
		return err
	}
	friends, err := f.db.GetAllFriendList(ctx)
	if err != nil {
		return err
	}
	hashCode := f.CalculateHash(friends)
	log.ZDebug(ctx, "SyncFriendPart", "serverHash", hashResp.Hash, "serverTotal", hashResp.Total, "localHash", hashCode, "localTotal", len(friends))
	if hashCode == hashResp.Hash {
		return nil
	}
	req := &friend.GetPaginationFriendsReq{
		UserID:     f.loginUserID,
		Pagination: &sdkws.RequestPagination{PageNumber: pconstant.FirstPageNumber, ShowNumber: pconstant.MaxSyncPullNumber},
	}
	resp, err := util.CallApi[friend.GetPaginationFriendsResp](ctx, constant.GetFriendListRouter, req)
	if err != nil {
		return err
	}
	serverFriends := util.Batch(ServerFriendToLocalFriend, resp.FriendsInfo)
	return f.friendSyncer.Sync(ctx, serverFriends, friends, nil)
}

func (f *Friend) IsLocalFriend(ctx context.Context) (bool, error) {
	if f.friendNum >= 0 {
		return f.friendNum < pconstant.MaxSyncPullNumber, nil
	}
	fs, err := f.db.GetAllFriendList(ctx)
	if err != nil {
		return false, err
	}
	return len(fs) < pconstant.MaxSyncPullNumber, nil
}

func (f *Friend) SearchFriendsV2(ctx context.Context, req *friend.SearchFriendsReq) ([]*sdk.SearchFriendItem, error) {
	local, err := f.IsLocalFriend(ctx)
	if err != nil {
		return nil, err
	}
	if local {
		fs, err := f.SearchFriends(ctx, &sdk.SearchFriendsParam{
			KeywordList:      []string{req.Keyword},
			IsSearchNickname: true,
			IsSearchRemark:   true,
			IsSearchUserID:   true,
		})
		if err != nil {
			return nil, err
		}
		return datautil.SlicePaginate(fs, req.Pagination), nil
	}
	resp, err := util.CallApi[friend.SearchFriendsResp](ctx, constant.SearchFriends, req)
	if err != nil {
		return nil, err
	}
	return f.toSearchFriendItem(ctx, util.Batch(ServerFriendToLocalFriend, resp.Friends))
}

func (f *Friend) SyncAllFriendList(ctx context.Context) error {
	return f.SyncFriendPart(ctx)
}

func (f *Friend) deleteLocalFriend(ctx context.Context, friendUserIDs []string) error {
	local, err := f.IsLocalFriend(ctx)
	if err != nil {
		return err
	}
	if local {
		friends, err := f.db.GetFriendInfoList(ctx, friendUserIDs)
		if err != nil {
			return err
		}
		for _, lf := range friends {
			if err := f.db.DeleteFriendDB(ctx, lf.FriendUserID); err != nil {
				return err
			}
			f.friendListener.OnFriendDeleted(*lf)
		}
	} else {
		lfs, err := f.db.GetFriendInfoList(ctx, friendUserIDs)
		if err != nil {
			return err
		}
		if len(lfs) == 0 {
			return nil
		}
		return f.SyncFriendPart(ctx)
	}
	return nil
}

func (f *Friend) SyncFriends(ctx context.Context, friendIDs []string) error {
	var resp friend.GetDesignatedFriendsResp
	if err := util.ApiPost(ctx, constant.GetDesignatedFriendsRouter, &friend.GetDesignatedFriendsReq{OwnerUserID: f.loginUserID, FriendUserIDs: friendIDs}, &resp); err != nil {
		return err
	}
	return f.deleteLocalFriend(ctx, friendIDs)
}

func (f *Friend) GetFriendListV2(ctx context.Context, pagination *sdkws.RequestPagination) ([]*server_api_params.FullUserInfo, error) {
	if util.LocalQuery(pagination) {
		users, err := f.GetFriendList(ctx)
		if err != nil {
			return nil, err
		}
		return datautil.SlicePaginate(users, pagination), nil
	}
	req := &friend.GetPaginationFriendsReq{
		UserID:     f.loginUserID,
		Pagination: &sdkws.RequestPagination{PageNumber: pconstant.FirstPageNumber, ShowNumber: pconstant.MaxSyncPullNumber},
	}
	resp, err := util.CallApi[friend.GetPaginationFriendsResp](ctx, constant.GetFriendListRouter, req)
	if err != nil {
		return nil, err
	}
	serverFriends := util.Batch(ServerFriendToLocalFriend, resp.FriendsInfo)
	return f.localFriendToFullUserInfo(ctx, serverFriends)
}
