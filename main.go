package main

import (
	"fmt"

	"PushMeClient/db"
	"PushMeClient/makelnk"
	"PushMeClient/server"
	"PushMeClient/utils"
	"PushMeClient/utils/setting"
	"PushMeClient/utils/sys"
	"PushMeClient/win"
)

func main() {
	makelnk.RegisterAppUserModeId()
	db.InitDb(utils.DbPath)
	sys.GetIps()

	mainCh := make(chan bool)
	w := win.OpenHome(mainCh)
	defer func() {
		w.Destroy()
	}()

	go func() {
		for enable := range mainCh {
			if enable {
				sys.MakeShortcut()
			} else {
				sys.RemoveShortcut()
			}
		}
	}()

	if setting.Setting.Api.Enable {
		go server.Start(func(msg db.Msg) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
					w.WebToast(r.(string))
				}
			}()

			go w.WebNewMessage(msg)
		})
	}

	w.Run()

	fmt.Println("PushMe Client已退出")
}
