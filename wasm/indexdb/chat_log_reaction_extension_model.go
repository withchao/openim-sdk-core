package indexdb

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

type LocalChatLogReactionExtensions struct {
}

func (i *LocalChatLogReactionExtensions) GetMessageReactionExtension(clientMsgID string) (result *model_struct.LocalChatLogReactionExtensions, err error) {
	msg, err := Exec(clientMsgID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msg.(string); ok {
			result := model_struct.LocalChatLogReactionExtensions{}
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return &result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalChatLogReactionExtensions) InsertMessageReactionExtension(messageReactionExtension *model_struct.LocalChatLogReactionExtensions) error {
	_, err := Exec(utils.StructToJsonString(messageReactionExtension))
	return err
}
func (i *LocalChatLogReactionExtensions) GetAndUpdateMessageReactionExtension(clientMsgID string, m map[string]*server_api_params.KeyValue) error {
	_, err := Exec(clientMsgID, utils.StructToJsonString(m))
	return err
}
func (i *LocalChatLogReactionExtensions) DeleteAndUpdateMessageReactionExtension(clientMsgID string, m map[string]*server_api_params.KeyValue) error {
	_, err := Exec(clientMsgID, utils.StructToJsonString(m))
	return err
}