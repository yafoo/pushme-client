package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	Url  string
	Data map[string]interface{}
}

func (r *Request) Get() string {
	uri := r.Url
	if r.Data != nil {
		u, _ := url.ParseRequestURI(r.Url)
		data := url.Values{}
		for k, v := range r.Data {
			data.Set(k, v.(string))
		}
		u.RawQuery = data.Encode()
		uri = u.String()
	}

	resp, err := http.Get(uri)
	if err != nil {
		panic(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	return string(body)
}

func (r *Request) PostForm() string {
	data := url.Values{}
	if r.Data != nil {
		for k, v := range r.Data {
			data.Set(k, v.(string))
		}
	}

	resp, err := http.PostForm(r.Url, data)
	if err != nil {
		panic(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	return string(body)
}

func (r *Request) Post() string {
	var jsonData []byte
	var err error
	if r.Data != nil {
		jsonData, err = json.Marshal(r.Data)
		if err != nil {
			panic(err.Error())
		}
	}

	resp, err := http.Post(r.Url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		panic(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	return string(body)
}
