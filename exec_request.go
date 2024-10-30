package webRequest

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	"github.com/gospider007/ja3"
	"github.com/gospider007/requests"
	utls "github.com/refraction-networking/utls"
	"github.com/tocha688/go_web_request/utils"
)

func ja3_random(ja3s string) string {
	arr := strings.Split(ja3s, ",")
	if arr[2] != "" {
		arr[2] = strings.Join(utils.ShuffleStrings[string](strings.Split(arr[2], "-")), "-")
	}
	return strings.Join(arr, ",")
}

func ja3_extension[V1 utls.TLSExtension](ja3s ja3.Ja3Spec) (V1, int) {
	for i, k := range ja3s.Extensions {
		if n, isok := k.(V1); isok {
			return n, i
		}
	}
	var zero V1
	return zero, -1
}

// ja3 修补程序
func ja3_repair(p *WebRequest, ja3s ja3.Ja3Spec) ja3.Ja3Spec {
	if p.IsChrome() {
		//修补 KeyShare
		tks, index := ja3_extension[*utls.KeyShareExtension](ja3s)
		if index != -1 {
			tks.KeyShares = []utls.KeyShare{
				{Group: utls.CurveID(utls.GREASE_PLACEHOLDER), Data: []byte{0}},
				{Group: utls.X25519Kyber768Draft00},
				{Group: utls.CurveP256},
			}
		}
	} else if p.IsFirefox() {
		tks, index := ja3_extension[*utls.KeyShareExtension](ja3s)
		if index != -1 {
			tks.KeyShares = []utls.KeyShare{
				{Group: utls.X25519},
				{Group: utls.CurveP256},
			}
		}
	}
	return ja3s
}

func (p *WebRequest) execute_requests(target string, method string) (res *WebResponse, err error) {
	p.Url = target
	p.Method = method
	ul, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	p.targetUrl = ul
	///
	ops := requests.RequestOption{
		TlsConfig:  &tls.Config{},
		UtlsConfig: &utls.Config{},
		//默认不重定向，因为要合并cookie
		MaxRedirect: -1,
	}
	//初始化指纹
	if p.client.SpecJa3 != "" {
		ja3d := p.client.SpecJa3
		//随机扩展
		if p.client.TlsRandom {
			ja3d = ja3_random(ja3d)
		}
		ja3Spec, err := ja3.CreateSpecWithStr(ja3d) //create ja3 spec with string
		if err != nil {
			return nil, err
		}
		//ja3修补
		ja3Spec = ja3_repair(p, ja3Spec)
		ops.Ja3Spec = ja3Spec
		ops.Ja3 = true
	}
	if p.client.SpecH2 != "" {
		sh2 := p.client.SpecH2
		sh2 = strings.ReplaceAll(sh2, ";", ",")
		h2Spec, err := ja3.CreateH2SpecWithStr(sh2) //create h2 spec with string
		if err != nil {
			return nil, err
		}
		ops.H2Ja3Spec = h2Spec
	}
	//代理
	if p.Proxy != "" {
		ops.Proxy = p.Proxy
	} else if p.client.Proxy != "" {
		ops.Proxy = p.client.Proxy
	}
	// if p.IsNotRedirect {
	// 	ops.MaxRedirect = -1
	// }
	if p.client.Http1Force || p.Http1Force {
		ops.ForceHttp1 = true
	}
	if p.client.SkipTlsVerify {
		ops.TlsConfig.InsecureSkipVerify = true
		ops.UtlsConfig.InsecureSkipVerify = true
	}
	ops.DisAlive = p.client.DisAlive
	//加载系统证书
	// x509pool, err := x509.SystemCertPool()
	// if err != nil {
	// 	return nil, err
	// }
	// ops.TlsConfig.RootCAs = x509pool
	// ops.UtlsConfig.RootCAs = x509pool
	//请求头
	hds := p.loadRequestHeader(target)
	//cookie
	if hds.Get("cookie") == "" {
		cookies := p.GetCookies()
		if len(cookies) > 0 {
			hds.Set("cookie", utils.CookieToString(cookies))
		}
	}
	ops.Headers = hds
	//重写请求头
	ops.OrderHeaders = DefaultHeaderOrders
	if p.Body != "" && p.isLoadBody(target) {
		ops.Body = p.Body
	}
	//开始请求
	resp, err := requests.Request(nil, method, target, ops)
	res = &WebResponse{
		Request:             p,
		request_headers:     http.Header(hds),
		RawRequestsResponse: resp,
	}
	if err != nil {
		return res, err
	}
	//回填cookie
	cookies := resp.Cookies()
	if cookies != nil {
		p.client.CookieJar.SetCookies(resp.Url(), cookies)
	}
	res.ResponseBytes = resp.Content()
	//判断是否有重定向
	loction := res.GetHeader("location")
	if !p.IsNotRedirect && loction != "" {
		//打印
		p.client.Println("Redirect: " + res.Url().String() + " --> " + loction)
		p.after_fn(res)
		//复制
		p2 := p.Clone()
		p2.SetHeader("referer", target)
		//删除数据
		p2.Body = ""
		p2.SetHeader("content-type", "")
		return p2.execute_requests(loction, "GET")
	} else {
		p.after_fn(res)
	}

	return res, nil
}
