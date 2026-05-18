package api

import (
	db "PushMe/internal/models"
	"PushMe/internal/setting"
	"PushMe/internal/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	httpServer *http.Server
	serverMux  sync.Mutex
)

var apiStatus map[string]string = map[string]string{
	"empty":   "title和content不能同时为空",
	"disable": "接口服务已关闭",
	"success": "success",
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}

	var body []byte
	var msg = db.Msg{
		Title:   r.Form.Get("title"),
		Content: r.Form.Get("content"),
		Date:    r.Form.Get("date"),
		Type:    r.Form.Get("type"),
	}

	jsonParams := &db.Msg{}
	if r.Method == http.MethodPost && strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(jsonParams); err != nil {
			log.Printf("JSON decode error: %v", err)
		}
	}

	if jsonParams.Title != "" {
		msg.Title = jsonParams.Title
	}
	if jsonParams.Content != "" {
		msg.Content = jsonParams.Content
	}
	if jsonParams.Date != "" {
		msg.Date = jsonParams.Date
	}
	if jsonParams.Type != "" {
		msg.Type = jsonParams.Type
	}

	if !setting.Setting.Api.Enable {
		body = []byte(apiStatus["disable"])
	} else if msg.Title == "" && msg.Content == "" {
		body = []byte(apiStatus["empty"])
	} else {
		body = []byte(apiStatus["success"])
		utils.EventEmit("msg", msg)
	}

	w.Write(body)
}

// Start 启动服务
func Start() {
	serverMux.Lock()
	defer serverMux.Unlock()

	// 如果服务已经在运行，先停止
	if httpServer != nil {
		Stop()
	}

	addr := setting.FormatIP(setting.Setting.Api.Ip) + ":" + setting.Setting.Api.Port

	httpServer = &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(handler),
	}

	// 在 goroutine 中启动服务
	go func() {
		log.Println("HTTP 服务启动在 " + addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP 服务启动失败: %v", err)
		}
	}()
}

// Stop 停止服务
func Stop() {
	serverMux.Lock()
	defer serverMux.Unlock()

	if httpServer == nil {
		log.Println("服务未运行，无需停止")
		return
	}

	log.Println("正在停止 HTTP 服务...")

	// 设置 5 秒超时，等待现有连接处理完成
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP 服务关闭失败: %v", err)
	} else {
		log.Println("HTTP 服务已停止")
	}

	httpServer = nil
}

// Restart 重启服务
func Restart() {
	Stop()
	// 等待端口释放
	time.Sleep(1 * time.Second)
	Start()
}

// IsRunning 检查服务是否运行
func IsRunning() bool {
	serverMux.Lock()
	defer serverMux.Unlock()
	return httpServer != nil
}
