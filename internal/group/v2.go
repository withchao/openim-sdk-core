package group

import (
	"crypto/md5"
	"encoding/binary"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	pconstant "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/utils/datautil"
	"strconv"
	"strings"
)

func (g *Group) CalculateHash(members []*model_struct.LocalGroupMember) uint64 {
	datautil.SortAny(members, func(a, b *model_struct.LocalGroupMember) bool {
		return a.JoinTime > b.JoinTime
	})
	if len(members) > pconstant.MaxSyncPullNumber {
		members = members[:pconstant.MaxSyncPullNumber]
	}
	hashStr := strings.Join(datautil.Slice(members, func(m *model_struct.LocalGroupMember) string {
		return strings.Join([]string{
			m.UserID,
			m.Nickname,
			m.FaceURL,
			strconv.FormatInt(int64(m.RoleLevel), 10),
			strconv.FormatInt(m.JoinTime, 10),
			strconv.FormatInt(int64(m.JoinSource), 10),
			m.InviterUserID,
			strconv.FormatInt(m.MuteEndTime, 10),
			m.OperatorUserID,
			m.Ex,
		}, ",")
	}), ";")
	sum := md5.Sum([]byte(hashStr))
	return binary.BigEndian.Uint64(sum[:])
}
