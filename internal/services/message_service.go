package services

import (
	"PushMe/constant"
	db "PushMe/internal/models"
	dbMsg "PushMe/internal/models/msg"
	"PushMe/internal/request"
	"PushMe/internal/setting"
	"PushMe/internal/toast"
	"log"
	"strings"
)

type MessageService struct{}

func (m *MessageService) GoLog(str string) {
	go log.Println(str)
}

func (m *MessageService) GoClose() {
	go log.Println("GoClose")
}

func (m *MessageService) GoGetVersion() string {
	return constant.AppVersion
}

func (m *MessageService) GoGetMessageList(page int, pageSize int) []db.Msg {
	return dbMsg.MsgList(page, pageSize)
}

func (m *MessageService) GoAddMessage(msg db.Msg) db.Msg {
	res := dbMsg.Add(&msg)
	if setting.Setting.Repost.Enable {
		go RepostMessage(res)
	}
	return res
}

func (m *MessageService) GoGetMessage(id int) db.Msg {
	return dbMsg.Get(id)
}

func (m *MessageService) GoDelMessage(id int) bool {
	return dbMsg.Del(id)
}

func (m *MessageService) GoGetMcount() int {
	return dbMsg.CountText()
}

func (m *MessageService) GoClearMessage() bool {
	return dbMsg.ClearText()
}

func (m *MessageService) GoGetDataList() []db.Msg {
	return dbMsg.DataList()
}

func (m *MessageService) GoOpenSetting(id int) {
	go log.Println("GoOpenSetting")
}
func (m *MessageService) GoOnMounted(id int) {
	go log.Println("GoOnMounted")
}

func containsAny(title string, keywords []string) bool {
	lowerTitle := strings.ToLower(title)
	for _, keyword := range keywords {
		if strings.Contains(lowerTitle, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func filterEmpty(keywords []string) []string {
	var result []string
	for _, keyword := range keywords {
		if strings.TrimSpace(keyword) != "" {
			result = append(result, keyword)
		}
	}
	return result
}

func RepostMessage(msg db.Msg) {
	if !setting.Setting.Repost.Enable || setting.Setting.Repost.Url == "" {
		return
	}

	if setting.Setting.Repost.Limit != "" {
		limitKeywords := filterEmpty(strings.Split(setting.Setting.Repost.Limit, "|"))
		containsAnyLimitKeyword := len(limitKeywords) == 0 || containsAny(msg.Title, limitKeywords)
		if !containsAnyLimitKeyword {
			return
		}
	}

	if setting.Setting.Repost.Omit != "" {
		omitKeywords := filterEmpty(strings.Split(setting.Setting.Repost.Omit, "|"))
		containsNoOmitKeywords := len(omitKeywords) == 0 || !containsAny(msg.Title, omitKeywords)
		if !containsNoOmitKeywords {
			return
		}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			toast.Toast(r.(string))
		}
	}()

	repost := request.Request{
		Url: setting.Setting.Repost.Url,
		Data: map[string]interface{}{
			"title":   msg.Title,
			"content": msg.Content,
			"date":    msg.Date,
			"type":    msg.Type,
		},
	}
	switch setting.Setting.Repost.Method {
	case "GET":
		repost.Get()
	case "POST/JSON":
		repost.Post()
	case "POST/FORM":
		repost.PostForm()
	default:
		repost.Post()
	}
	log.Println("RepostMessage")
}
