package webRequest

import (
	"net/http"
	"net/url"

	"github.com/samber/lo"
	"github.com/tocha688/go-cookiejar"
	"github.com/tocha688/go_web_reuqest/utils"
)

type WebClient struct {
	Proxy         string
	TlsRandom     bool
	Http1Force    bool
	CookieJar     *cookiejar.Jar
	Cookies       []*http.Cookie
	Header        *WebHeader
	IsDebug       bool
	SkipTlsVerify bool
	//指纹
	SpecJa3 string
	SpecH2  string
	//禁用保活
	DisAlive bool
}

func (p *WebClient) init() {
	if p.Header == nil {
		p.Header = &WebHeader{}
	}
	if p.CookieJar == nil {
		p.CookieJar, _ = cookiejar.New(&cookiejar.Options{
			NoPersist: true,
		})
	}
}

// ---------------cookie---------------
//
//	func (p *WebClient) SetCookie(key string, value string) {
//		p.init()
//		p.Cookies = append(p.Cookies, &http.Cookie{Name: key, Value: value})
//	}
func (p *WebClient) SetCookie(cookie *http.Cookie) *WebClient {
	p.init()
	p.Cookies = append(p.Cookies, cookie)
	return p
}
func (p *WebClient) SetJarCookie(ul *url.URL, cos ...*http.Cookie) *WebClient {
	p.init()
	p.CookieJar.SetCookies(ul, cos)
	return p
}
func (p *WebClient) SetJarCookies(ul *url.URL, cookies []*http.Cookie) *WebClient {
	p.init()
	for _, cookie := range cookies {
		p.SetJarCookie(ul, cookie)
	}
	return p
}
func (p *WebClient) GetJarCookies(ul *url.URL) []*http.Cookie {
	p.init()
	rcos := []*http.Cookie{}
	for _, co := range p.CookieJar.Cookies(ul) {
		cookie := &http.Cookie{
			Name:       co.Name,
			Value:      co.Value,
			Domain:     co.Domain,
			Path:       co.Path,
			Expires:    co.Expires,
			RawExpires: co.RawExpires,
			MaxAge:     co.MaxAge,
			Secure:     co.Secure,
			HttpOnly:   co.HttpOnly,
			SameSite:   http.SameSite(co.SameSite),
			Raw:        co.Raw,
			Unparsed:   co.Unparsed,
		}
		rcos = append(rcos, cookie)
	}
	return rcos
}
func (p *WebClient) GetAllCookies() []*http.Cookie {
	p.init()
	rcos := p.CookieJar.AllCookies()
	return utils.CookieMergen(p.Cookies, rcos)
}
func (p *WebClient) SetCookies(cookies []*http.Cookie) *WebClient {
	p.init()
	p.Cookies = utils.CookieMergen(p.Cookies, cookies)
	return p
}
func (p *WebClient) ClearCookies() *WebClient {
	p.init()
	p.Cookies = []*http.Cookie{}
	p.CookieJar, _ = cookiejar.New(nil)
	return p
}
func (p *WebClient) GetCookie(name string) *http.Cookie {
	co, isFind := lo.Find(p.GetAllCookies(), func(it *http.Cookie) bool {
		return it.Name == name
	})
	if isFind {
		return co
	}
	return nil
}
func (p *WebClient) RemoveCookie(names ...string) *WebClient {
	p.Cookies = lo.Filter(p.Cookies, func(it *http.Cookie, ix int) bool {
		_, i2 := lo.Find(names, func(name string) bool {
			return it.Name == name
		})
		return !i2
	})
	for _, co := range p.CookieJar.AllCookies() {
		_, i2 := lo.Find(names, func(name string) bool {
			return co.Name == name
		})
		if i2 {
			p.CookieJar.RemoveCookie(co)
		}
	}
	return p
}

// 请求前合并cookie
func (p *WebClient) GetCookies() []*http.Cookie {
	p.init()
	//获取全部cookie
	rcos := p.CookieJar.AllCookies()
	rcos = utils.CookieMergen(p.Cookies, rcos)
	//转为字符串
	return rcos
}

// ---------------head---------------
func (p *WebClient) SetHeader(k string, v string) *WebClient {
	p.init()
	p.Header.Set(k, v)
	return p
}
func (p *WebClient) DelHeader(k string) *WebClient {
	p.init()
	p.Header.Del(k)
	return p
}
func (p *WebClient) SetHeaders(headers map[string]string) *WebClient {
	p.init()
	for k, v := range headers {
		p.Header.Set(k, v)
	}
	return p
}

// ---------------others---------------
func (p *WebClient) R() *WebRequest {
	p.init()
	return &WebRequest{
		Header:        &WebHeader{},
		Body:          "",
		Cookies:       []*http.Cookie{},
		IsNotRedirect: false,
		Timeout:       30,
	}
}

func (p *WebClient) SetDebug(debug bool) *WebClient {
	p.init()
	p.IsDebug = debug
	return p
}

func (p *WebClient) SetProxy(u string) *WebClient {
	p.init()
	p.Proxy = u
	return p
}
func (p *WebClient) SetDisAlive(b bool) *WebClient {
	p.init()
	p.DisAlive = b
	return p
}
func (p *WebClient) SetSkipTlsVerify(b bool) *WebClient {
	p.init()
	p.SkipTlsVerify = b
	return p
}

// ----------------
func New() *WebClient {
	r := &WebClient{}
	r.init()
	return r
}
