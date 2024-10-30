package webRequest

var DefaultHeaderOrders = []string{
	"host",
	"connection",
	"content-length",
	"pragma",
	"cache-control",
	"device-memory",
	"viewport-width",
	"rtt",
	"downlink",
	"ect",
	"sec-ch-ua",
	"sec-ch-ua-mobile",
	"sec-ch-ua-full-version",
	"sec-ch-ua-arch",
	"sec-ch-ua-platform",
	"sec-ch-ua-platform-version",
	"sec-ch-ua-model",
	"origin",
	"content-type",
	"upgrade-insecure-requests",
	"user-agent",
	"accept",
	"sec-fetch-site",
	"sec-fetch-mode",
	"sec-fetch-user",
	"sec-fetch-dest",
	"referer",
	"accept-encoding",
	"accept-language",
	"cookie",
	"priority",
}

func (p *WebRequest) after_fn(res *WebResponse) {
	if p.client.IsDebug {
		res.PrintDebugger()
	}
}

func (p *WebRequest) Execute(target string, method string) (*WebResponse, error) {
	//先判断是否有cors
	if p.IsCors && method != "OPTIONS" {
		resp, err := p.Clone().Options(target, method, p.Header.GetValue("content-type"))
		if err != nil {
			return resp, err
		}
	}
	return p.execute_requests(target, method)
}

func (p *WebRequest) Post(target string) (*WebResponse, error) {
	return p.Execute(target, "POST")
}
func (p *WebRequest) Get(target string) (*WebResponse, error) {
	return p.Execute(target, "GET")
}
func (p *WebRequest) Patch(target string) (*WebResponse, error) {
	return p.Execute(target, "PATCH")
}
func (p *WebRequest) Options(target string, method string, headers string) (*WebResponse, error) {
	p.SetHeaders(map[string]string{
		"access-control-request-headers": headers,
		"access-control-request-method":  method,
		"sec-fetch-dest":                 "empty",
		"sec-fetch-mode":                 "cors",
		"sec-fetch-site":                 "cross-site",
	})
	return p.Execute(target, "OPTIONS")
}
