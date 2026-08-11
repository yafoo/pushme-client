package services

import (
	db "PushMe/internal/models"
	dbMsg "PushMe/internal/models/msg"
	"PushMe/internal/setting"
	"PushMe/internal/utils"
)

type MessageService struct{}

func (m *MessageService) GoGetMessageList(page int, pageSize int) []db.Msg {
	return dbMsg.MsgList(page, pageSize)
}

func (m *MessageService) GoAddMessage(msg db.Msg) db.Msg {
	res := dbMsg.Add(&msg)
	if res.ID > 0 && setting.Setting.Repost.Enable {
		go utils.RepostMessage(res)
	}
	return res
}

func (m *MessageService) GoGetMessage(id int) db.Msg {
	return dbMsg.Get(id)
}

func (m *MessageService) GoDelMessage(id int) bool {
	return dbMsg.Del(id)
}

func (m *MessageService) GoUpdateMessage(msg db.Msg) bool {
	res := dbMsg.Update(&msg)
	if res && setting.Setting.Repost.Enable {
		go utils.RepostMessage(msg)
	}
	return res
}

func (m *MessageService) GoGetMessageCount() int {
	return dbMsg.CountText()
}

func (m *MessageService) GoClearMessage() bool {
	return dbMsg.ClearText()
}

func (m *MessageService) GoGetDataList() []db.Msg {
	return dbMsg.DataList()
}
