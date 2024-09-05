package win

import (
	"PushMeClient/constant"
	"PushMeClient/db"
	dbMsg "PushMeClient/db/msg"
	dbPlugin "PushMeClient/db/plugin"
	"PushMeClient/utils/request"
	"PushMeClient/utils/setting"
	"PushMeClient/utils/sys"
	"encoding/json"
	"fmt"

	webview "github.com/webview/webview_go"
)

var mainCh chan bool

type Win struct {
	Webview webview.WebView
}

func (w *Win) Run() {
	w.Webview.Run()
}

func (w *Win) Destroy() {
	w.Webview.Destroy()
}

func (w *Win) Terminate() {
	w.Webview.Terminate()
}

func (w *Win) Dispatch(f func()) {
	w.Webview.Dispatch(f)
}

func (w *Win) WebNewMessage(msg db.Msg) {
	message, err := json.Marshal(msg)
	if err != nil {
		panic(err.Error())
	}
	w.Webview.Dispatch(func() {
		w.Webview.Eval("WebNewMessage(" + string(message) + ");")
	})
}

func (w *Win) WebToast(txt string) {
	w.Webview.Dispatch(func() {
		w.Webview.Eval("WebToast('" + txt + "');")
	})
}

func NewWin(title string, url string, width int, height int) Win {
	w := Win{}
	w.Webview = webview.New(true)
	w.Webview.SetTitle(title)
	w.Webview.Navigate(url)
	w.Webview.SetSize(width, height, webview.HintNone)
	w.Webview.Bind("GoNotification", func(msg db.Msg) {
		go sys.Notification(msg)
	})
	w.Webview.Bind("GoLog", func(log string) {
		go fmt.Println(log)
	})
	w.Webview.Bind("GoClose", func() {
		w.Terminate()
	})
	w.Webview.Bind("GoGetVersion", func() string {
		return constant.AppVersion
	})
	w.Webview.Bind("GoCheckVersion", func() string {
		return CheckVersion(w)
	})
	return w
}

func OpenHome(ch chan bool) Win {
	mainCh = ch
	w := NewWin(constant.AppName, setting.BaseApi, 300, 500)
	w.Webview.Bind("GoGetMessageList", func(page int, pageSize int) []db.Msg {
		return dbMsg.MsgList(page, pageSize)
	})
	w.Webview.Bind("GoAddMessage", func(msg db.Msg) db.Msg {
		return dbMsg.Add(&msg)
	})
	w.Webview.Bind("GoOpenMessage", func(id int) {
		go OpenMessage(id, w)
	})

	w.Webview.Bind("GoOpenDashboard", func() {
		go OpenDashboard(w)
	})
	w.Webview.Bind("GoOpenSetting", func() {
		go OpenSetting(w)
	})

	w.Webview.Bind("GoOnMounted", func() {
		if setting.Setting.Host.Enable {
			host, err := json.Marshal(setting.Setting.Host)
			if err != nil {
				panic(err)
			}
			go w.Webview.Dispatch(func() {
				w.Webview.Eval("WebHostConnect(" + string(host) + ");")
			})
		}
	})
	w.Webview.Bind("GoGetPluginCount", func() int {
		return dbPlugin.CountWithState(1)
	})
	w.Webview.Bind("GoGetRepostState", func() bool {
		return setting.Setting.Repost.Enable
	})
	w.Webview.Bind("GoRepostMessage", func(msg db.Msg) {
		if setting.Setting.Repost.Enable {
			go RepostMessage(msg, w)
		}
	})

	return w
}

func OpenDashboard(mWin Win) {
	w := NewWin("数据小屏", setting.BaseApi+"/view/dashboard.html", 300, 500)
	defer func() {
		w.Destroy()
		mWin.Dispatch(func() {
			mWin.Terminate()
		})
	}()
	w.Webview.Bind("GoGetDataList", func() []db.Msg {
		return dbMsg.DataList()
	})
	w.Webview.Bind("GoOpenMessage", func(id int) {
		go OpenMessage(id, w)
	})
	w.Run()
}

func OpenMessage(id int, mWin Win) {
	w := NewWin("消息内容", setting.BaseApi+"/view/message.html?id="+fmt.Sprint(id), 300, 500)
	defer func() {
		w.Destroy()
		mWin.Dispatch(func() {
			mWin.Terminate()
		})
	}()
	w.Webview.Bind("GoGetMessage", func(id int) db.Msg {
		return dbMsg.Get(id)
	})
	w.Webview.Bind("GoDelMessage", func(id int) bool {
		return dbMsg.Del(id)
	})
	w.Run()
}

func OpenSetting(mWin Win) {
	w := NewWin("系统设置", setting.BaseApi+"/view/setting.html", 300, 500)
	defer func() {
		w.Destroy()
		mWin.Dispatch(func() {
			mWin.Terminate()
		})
	}()
	w.Webview.Bind("GoGetSetting", func() setting.SettingType {
		return setting.Setting
	})
	w.Webview.Bind("GoGetSettingDefault", func() setting.SettingType {
		return setting.GetSettingDefault()
	})
	w.Webview.Bind("GoSaveSetting", func(settingData setting.SettingType) bool {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
				w.WebToast(r.(string))
			}
		}()
		setting.SaveSetting(settingData)
		mainCh <- settingData.System.Enable
		return true
	})
	w.Webview.Bind("GoOpenPlugin", func() {
		go OpenPlugin(mWin)
	})
	w.Webview.Bind("GoGetIps", func() []string {
		return sys.GetIps()
	})
	w.Webview.Bind("GoGetMcount", func() int {
		return dbMsg.CountText()
	})
	w.Webview.Bind("GoClearMessage", func() bool {
		return dbMsg.ClearText()
	})
	w.Webview.Bind("GoRestart", func() {
		go sys.Restart()
	})
	w.Run()
}

func OpenPlugin(mWin Win) {
	w := NewWin("消息插件", setting.BaseApi+"/view/plugin.html", 300, 500)
	defer func() {
		w.Destroy()
		mWin.Dispatch(func() {
			mWin.Terminate()
		})
	}()
	w.Webview.Bind("GoGetPluginList", func() []db.Plugin {
		return dbPlugin.List()
	})
	w.Webview.Bind("GoOpenPluginEdit", func(id int) {
		go OpenPluginEdit(id, mWin)
	})
	w.Webview.Bind("GoEditPluginState", func(plugin db.Plugin) int {
		return dbPlugin.Add(&plugin)
	})
	w.Webview.Bind("GoSwapPlugin", func(plugin1 db.Plugin, plugin2 db.Plugin) bool {
		return dbPlugin.SwapSort(&plugin1, &plugin2)
	})
	w.Run()
}

func OpenPluginEdit(id int, mWin Win) {
	w := NewWin("插件编辑", setting.BaseApi+"/view/plugin_edit.html?id="+fmt.Sprint(id), 300, 500)
	defer func() {
		w.Destroy()
		mWin.Dispatch(func() {
			mWin.Terminate()
		})
	}()
	w.Webview.Bind("GoGetPlugin", func(id int) db.Plugin {
		return dbPlugin.Get(id)
	})
	w.Webview.Bind("GoEditPlugin", func(plugin db.Plugin) int {
		return dbPlugin.Add(&plugin)
	})
	w.Webview.Bind("GoDelPlugin", func(id int) bool {
		return dbPlugin.Del(id)
	})
	w.Run()
}

func RepostMessage(msg db.Msg, w Win) {
	if !setting.Setting.Repost.Enable || setting.Setting.Repost.Url == "" {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			w.WebToast(r.(string))
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
	fmt.Println("RepostMessage")
}

func CheckVersion(w Win) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			w.WebToast(r.(string))
		}
	}()

	req := request.Request{
		Url: constant.UrlVersion,
		Data: map[string]interface{}{
			"version": constant.AppVersion,
		},
	}
	result := req.Post()
	return result
}
