package webRequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/tocha688/go_web_request/utils"
)

type WebRequest struct {
	client        *WebClient
	Url           string
	targetUrl     *url.URL
	Method        string
	Body          string
	Header        *WebHeader
	Cookies       []*http.Cookie
	IsNotRedirect bool
	IsCors        bool
	//秒
	Timeout          int
	ServerName       string
	_NotClientHeader bool
	Http1Force       bool
	Proxy            string
}

func (p *WebRequest) Clone() *WebRequest {
	return &WebRequest{
		Url:           p.Url,
		Method:        p.Method,
		Body:          p.Body,
		Header:        utils.DeepClone(p.Header),
		Cookies:       p.Cookies,
		IsNotRedirect: p.IsNotRedirect,
		IsCors:        p.IsCors,
		Timeout:       p.Timeout,
		// targetUrl:        utils.DeepClone(p.targetUrl),
		_NotClientHeader: p._NotClientHeader,
		Http1Force:       p.Http1Force,
		Proxy:            p.Proxy,
		client:           p.client,
	}
}

func (p *WebRequest) SetProxy(u string) *WebRequest {
	p.Proxy = u
	return p
}
func (p *WebRequest) SetCors() *WebRequest {
	p.IsCors = true
	return p
}
func (p *WebRequest) NotBaseHeader() *WebRequest {
	p._NotClientHeader = false
	return p
}
func (p *WebRequest) SetBodyUrlencoded(body map[string]string) *WebRequest {
	p.Body = utils.QueryStringify(body)
	p.Header.Set("content-type", "application/x-www-form-urlencoded")
	return p
}
func (p *WebRequest) SetBody(body string) *WebRequest {
	p.Body = body
	return p
}
func (p *WebRequest) SetBodyJson(body map[string]any) *WebRequest {
	jsonBytes, _ := json.Marshal(body)
	p.Body = string(jsonBytes)
	p.Header.Set("content-type", "application/json")
	return p
}
func (p *WebRequest) SetBodyJsonV1(body map[string]any) *WebRequest {
	jsonBytes, _ := json.Marshal(body)
	p.Body = string(jsonBytes)
	p.Header.Set("content-type", "application/vnd.bc.v1+json")
	return p
}
func (p *WebRequest) SetFormData(data map[string]string) *WebRequest {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	writer.SetBoundary("WebKitFormBoundary" + utils.RandomHex(16))
	for key, value := range data {
		part, err := writer.CreateFormField(key)
		if err != nil {
			fmt.Println("Error creating form field:", err)
		}
		_, err = io.WriteString(part, value)
		if err != nil {
			fmt.Println("Error writing to form field:", err)
		}
	}
	p.Header.Set("content-type", writer.FormDataContentType())
	// p.SetBody(requestBody.String() + "\n--" + writer.Boundary() + "--")
	err := writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
	}
	p.SetBody(requestBody.String())
	return p
}
func (p *WebRequest) NotRedirect() *WebRequest {
	p.IsNotRedirect = true
	return p
}
func (p *WebRequest) SetCookie(cok *http.Cookie) *WebRequest {
	p.Cookies = append(p.Cookies, cok)
	return p
}
func (p *WebRequest) SetCookies(cok []*http.Cookie) *WebRequest {
	p.Cookies = append(p.Cookies, cok...)
	return p
}
func (p *WebRequest) SetHeader(key string, value string) *WebRequest {
	p.Header.Set(key, value)
	return p
}
func (p *WebRequest) SetHeaders(headers map[string]string) *WebRequest {
	for k, v := range headers {
		p.Header.Set(k, v)
	}
	return p
}

// --------------- 合并获取请求头
func (req *WebRequest) GetHeaders() map[string]string {
	headers := make(map[string]string)
	if !req._NotClientHeader {
		for _, hd := range req.client.Header.headers {
			headers[hd.Name] = hd.Value
		}
	}
	for _, hd := range req.Header.headers {
		headers[hd.Name] = hd.Value
	}
	return headers
}
func (p *WebRequest) GetHttpHeader() http.Header {
	headers := http.Header{}
	for key, co := range p.GetHeaders() {
		headers.Set(key, co)
	}
	return headers
}

// ---------------- 合并获取cookie
func (p *WebRequest) GetCookies() []*http.Cookie {
	cookies := p.client.GetAllCookies()
	cookies = utils.CookieMergen(cookies, p.Cookies)
	return cookies
}

// ----检测浏览器类型----
func (p *WebRequest) UaIs(s string) bool {
	h := p.Header.Get("user-agent")
	if h == nil {
		return false
	}
	return strings.Contains(h.Value, s)
}
func (p *WebRequest) IsChrome() bool {
	return p.UaIs("Chrome")
}
func (p *WebRequest) IsFirefox() bool {
	return p.UaIs("Firefox")
}
func (p *WebRequest) isLoadBody(target string) bool {
	return target == "POST" || target == "PUT" || target == "PATCH"
}
func (p *WebRequest) loadRequestHeader(target string) http.Header {
	headers := http.Header{}
	for key, co := range p.GetHeaders() {
		if key == "content-type" && !p.isLoadBody(target) {
			continue
		}
		headers.Set(key, co)
	}
	return headers
}
