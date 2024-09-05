package server

import (
	"PushMeClient/db"
	"PushMeClient/db/plugin"
	"PushMeClient/server/html"
	"PushMeClient/utils/setting"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Start(callback func(msg db.Msg)) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		r.ParseForm()
		if r.Form.Get("title") == "" && r.Form.Get("content") == "" {
			if r.ContentLength == 0 {
				index, err := html.View.ReadFile("view/index.html")
				if err != nil {
					fmt.Println(err)
				}
				w.Write(index)
			} else if r.Header.Get("Content-Type") == "application/json" {
				body, _ := io.ReadAll(r.Body)
				msg := db.Msg{}
				err := json.Unmarshal(body, &msg)
				if err == nil {
					callback(msg)
					w.Write([]byte("success"))
				} else {
					fmt.Println(err)
					w.Write([]byte(err.Error()))
				}
			}
		} else {
			msg := db.Msg{
				Title:   r.Form.Get("title"),
				Content: r.Form.Get("content"),
				Date:    r.Form.Get("date"),
				Type:    r.Form.Get("type"),
			}
			callback(msg)
		}
	})

	http.Handle("/view/", http.FileServer(http.FS(html.View)))
	http.Handle("/static/", http.FileServer(http.FS(html.Static)))
	http.HandleFunc("/static/plugin.js", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		pluginTpl, err := html.Static.ReadFile("static/plugin.js")
		if err != nil {
			fmt.Println(err)
		}

		pluginList := plugin.ListWithState(1)
		if len(pluginList) == 0 {
			w.Write(pluginTpl)
		} else {
			var pluginStr string
			for _, plugin := range pluginList {
				pluginStr += fmt.Sprintf("plugins.push(\n    %s\n);\n", strings.ReplaceAll(plugin.Content, "\n", "\n    "))
			}
			pluginJs := strings.Replace(string(pluginTpl), "// plugin_list", pluginStr, 1)
			w.Write([]byte(pluginJs))
		}
	})

	fmt.Println("http started on " + setting.BaseApi)
	http.ListenAndServe(":"+setting.Setting.Api.Port, nil)
}
