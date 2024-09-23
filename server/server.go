package server

import (
	"PushMeClient/db"
	"PushMeClient/db/plugin"
	"PushMeClient/server/html"
	"PushMeClient/utils/setting"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var apiStatus map[string]string = map[string]string{
	"empty":   "title和content不能同时为空",
	"disable": "接口服务已关闭",
	"success": "success",
}

func Start(callback func(msg db.Msg)) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body []byte
		var msg = db.Msg{}

		r.ParseForm()
		if r.Form.Get("title") == "" && r.Form.Get("content") == "" && r.ContentLength == 0 {
			if r.Method == "GET" {
				tpl, err := html.View.ReadFile("view/index.html")
				if err != nil {
					fmt.Println(err)
				}
				body = tpl
			} else if !setting.Setting.Api.Enable {
				body = []byte(apiStatus["disable"])
			} else {
				body = []byte(apiStatus["empty"])
			}
		} else if !setting.Setting.Api.Enable {
			body = []byte(apiStatus["disable"])
		} else {
			if r.Form.Get("title") != "" {
				msg.Title = r.Form.Get("title")
			}
			if r.Form.Get("content") != "" {
				msg.Content = r.Form.Get("content")
			}
			if r.Form.Get("content") != "" {
				msg.Date = r.Form.Get("date")
			}
			if r.Form.Get("content") != "" {
				msg.Type = r.Form.Get("type")
			}

			if strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "application/json") {
				data := db.Msg{}
				err := json.NewDecoder(r.Body).Decode(&data)
				if err != nil {
					fmt.Println(err)
					body = []byte(err.Error())
				} else {
					if data.Title != "" {
						msg.Title = data.Title
					}
					if data.Content != "" {
						msg.Content = data.Content
					}
					if data.Date != "" {
						msg.Date = data.Date
					}
					if data.Type != "" {
						msg.Type = data.Type
					}
				}
			} else if r.Method == "POST" {
				if r.PostForm.Get("title") != "" {
					msg.Title = r.Form.Get("title")
				}
				if r.PostForm.Get("content") != "" {
					msg.Content = r.Form.Get("content")
				}
				if r.PostForm.Get("date") != "" {
					msg.Date = r.Form.Get("date")
				}
				if r.PostForm.Get("type") != "" {
					msg.Type = r.Form.Get("type")
				}
			}
		}

		if len(body) == 0 {
			if msg.Title == "" && msg.Content == "" {
				body = []byte(apiStatus["empty"])
			} else {
				callback(msg)
				body = []byte(apiStatus["success"])
			}
		}

		w.Write(body)
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
