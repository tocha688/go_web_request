package webRequest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gospider007/requests"
)

type WebResponse struct {
	Request             *WebRequest
	RawRequestsResponse *requests.Response
	ResponseBytes       []byte
	request_headers     http.Header
}

func (p *WebResponse) String() string {
	return string(p.ResponseBytes)
}
func (p *WebResponse) Status() string {
	return p.RawRequestsResponse.Status()
}
func (p *WebResponse) StatusCode() int {
	return p.RawRequestsResponse.StatusCode()
}
func (p *WebResponse) GetHeader(name string) string {
	return p.GetHeaders().Get(name)
}
func (p *WebResponse) GetHeaders() http.Header {
	return p.RawRequestsResponse.Headers()
}
func (p *WebResponse) Cookies() []*http.Cookie {
	return p.RawRequestsResponse.Cookies()
}
func (p *WebResponse) Url() *url.URL {
	return p.RawRequestsResponse.Url()
}
func (p *WebResponse) GetRequestHeader() http.Header {
	//暂时无法requests获取请求头,因为他没有暴露原始response
	return p.Request.GetHttpHeader()
}

func (p *WebResponse) PrintDebugger() {
	str := fmt.Sprintln("Request: ", p.Request.targetUrl)
	str += fmt.Sprintln("")
	for k, v := range p.GetRequestHeader() {
		str = str + fmt.Sprintln(k, ":", v)
	}
	//cookie
	// str = str + fmt.Sprintln("Cookies: ")
	// for _, cok := range req.Cookies() {
	// 	str = str + fmt.Sprintln("\t", cok.Name, ":", cok.Value)
	// }
	str = str + fmt.Sprintln("\n", p.Request.Body)
	str = str + fmt.Sprintln("-------------Response-------------")
	str = str + fmt.Sprintln("Url: ", p.Url())
	str += fmt.Sprintln("")
	for k, v := range p.GetHeaders() {
		str = str + fmt.Sprintln(k, ":", v[0])
	}
	str = str + fmt.Sprintln("Status: ", p.Status())
	str = str + fmt.Sprintln("StatusCode: ", p.StatusCode())
	body := p.String()
	if len(body) < 1000 {
		str += fmt.Sprintln("\n", body)
	} else {
		str += fmt.Sprintln(body[:1000], "\n - ResponseBodySize:", len(body))
	}
	fmt.Println(str + "\n")
}
