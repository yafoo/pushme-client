package api

import (
	db "PushMe/internal/models"
	"PushMe/internal/setting"
	"PushMe/internal/utils"
	"context"
	"encoding/json"
	"io"
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
	"pushkey": "push_key 错误",
	"success": "success",
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var msg = db.Msg{}
	var pushKey string

	contentType := r.Header.Get("Content-Type")
	isJSON := r.Method == http.MethodPost && strings.Contains(contentType, "application/json")

	if isJSON {
		// JSON 请求：一次性读取 body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("read body error: %v", err)
			w.Write([]byte("read body error"))
			return
		}
		jsonParams := &db.Msg{}
		if err := json.Unmarshal(bodyBytes, jsonParams); err != nil {
			log.Printf("JSON decode error: %v", err)
		}
		msg.Title = jsonParams.Title
		msg.Content = jsonParams.Content
		msg.Date = jsonParams.Date
		msg.Type = jsonParams.Type
		pushKey = jsonParams.PushKey
	} else if r.Method == http.MethodPost && strings.Contains(contentType, "multipart/form-data") {
		// multipart/form-data 请求
		// 注意：Go 1.25 中 r.ParseForm() 无法正确填充 PostForm，必须显式调用 ParseMultipartForm
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Printf("parse multipart form error: %v", err)
		}
		msg.Title = r.PostForm.Get("title")
		msg.Content = r.PostForm.Get("content")
		msg.Date = r.PostForm.Get("date")
		msg.Type = r.PostForm.Get("type")
		pushKey = r.PostForm.Get("push_key")
	} else if r.Method == http.MethodPost {
		// urlencoded 表单请求
		if err := r.ParseForm(); err != nil {
			log.Printf("parse form error: %v", err)
		}
		msg.Title = r.PostForm.Get("title")
		msg.Content = r.PostForm.Get("content")
		msg.Date = r.PostForm.Get("date")
		msg.Type = r.PostForm.Get("type")
		pushKey = r.PostForm.Get("push_key")
	}

	// URL 查询参数作为兜底（POST body 参数优先级更高）
	if msg.Title == "" {
		msg.Title = r.URL.Query().Get("title")
	}
	if msg.Content == "" {
		msg.Content = r.URL.Query().Get("content")
	}
	if msg.Date == "" {
		msg.Date = r.URL.Query().Get("date")
	}
	if msg.Type == "" {
		msg.Type = r.URL.Query().Get("type")
	}
	if pushKey == "" {
		pushKey = r.URL.Query().Get("push_key")
	}

	var body []byte

	if !setting.Setting.Api.Enable {
		body = []byte(apiStatus["disable"])
	} else if setting.Setting.Api.VerifyKey {
		expectedKey := setting.Setting.Host.PushKey
		if pushKey == "" || pushKey != expectedKey {
			body = []byte(apiStatus["pushkey"])
		} else if msg.Title == "" && msg.Content == "" {
			body = []byte(apiStatus["empty"])
		} else {
			body = []byte(apiStatus["success"])
			utils.EventEmit("message:api", msg)
		}
	} else if msg.Title == "" && msg.Content == "" {
		body = []byte(apiStatus["empty"])
	} else {
		body = []byte(apiStatus["success"])
		utils.EventEmit("message:api", msg)
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
