package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	port = "3200"
)

func initServer(callback func(msg Msg)) {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exeName := filepath.Base(exePath)
	exeName = strings.ReplaceAll(exeName, ".exe", "")
	params := strings.Split(exeName, "-")
	if len(params) > 1 {
		port = params[1]
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := io.ReadAll(r.Body)
		fmt.Println(string(body))
		msg := Msg{}
		err := json.Unmarshal(body, &msg)
		if err == nil {
			callback(msg)
		}
	})

	fmt.Println("http started on " + port)
	http.ListenAndServe(":"+port, nil)
}
