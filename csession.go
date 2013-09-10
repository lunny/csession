package csession

import (
	"net/http"
	"net/http/cookiejar"
)

// a simple wrapper for http.RoundTripper to do something before and after RoundTrip
type transport struct {
	tr        http.RoundTripper
	BeforeReq func(req *http.Request)
	AfterReq  func(resp *http.Response, req *http.Request)
}

func newTransport(tr http.RoundTripper) *transport {
	t := &transport{}
	if tr == nil {
		tr = http.DefaultTransport
	}
	t.tr = tr
	return t
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t.BeforeReq(req)
	resp, err = t.tr.RoundTrip(req)
	if err != nil {
		return
	}
	t.AfterReq(resp, req)
	return
}

// struct session has an anonymous field http.Client, so it has all public methods
// of Client. And it add the referer, headers, cookies to request before it to do.
type Session struct {
	referer string
	cookies []*http.Cookie
	http.Client
	HeadersFunc func(req *http.Request)
}

func (s *Session) insertReferer(req *http.Request) {
	if s.referer != "" {
		req.Header.Set("Referer", s.referer)
	}
}

func (s *Session) insertCookie(req *http.Request) {
	for _, cookie := range s.cookies {
		req.AddCookie(cookie)
	}
}

func (s *Session) mergeCookie(resp *http.Response) {
	cookies := resp.Cookies()
	newCookies := make([]*http.Cookie, len(cookies))
	length := 0
	for _, c := range cookies {
		for idx, cs := range s.cookies {
			if c.Name == cs.Name {
				s.cookies[idx] = c
				goto next
			}
		}
		newCookies[length] = c
		length++
	next:
		continue
	}

	s.cookies = append(s.cookies, newCookies[:length]...)
}

// default headers, you can use your self headers func and set to session.HeadersFunc
func DefaultHeadersFunc(req *http.Request) {
	accept := "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
	agent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.116 Safari/537.36"
	encoding := "none" //gzip, deflate
	language := "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3"
	req.Header.Set("Accept", accept)
	req.Header.Set("User-Agent", agent)
	req.Header.Set("Accept-Encoding", encoding)
	req.Header.Set("Accept-Language", language)
	//req.Header.Set("Cache-Control", "max-age=0")
	//req.Header.Set("Connection", "keep-alive")
}

// new a default Session
func New() *Session {
	jar, err := cookiejar.New(nil)
	if err != nil {
		jar = nil
	}
	return NewSession(http.DefaultTransport, nil, jar)
}

// new a custom session, the params is the same as Client's public fields
func NewSession(transport http.RoundTripper,
	checkRedirect func(req *http.Request, via []*http.Request) error,
	jar http.CookieJar) *Session {
	s := &Session{HeadersFunc: DefaultHeadersFunc}
	newTr := newTransport(transport)
	newTr.AfterReq = func(resp *http.Response, req *http.Request) {
		s.mergeCookie(resp)
		s.referer = req.URL.String()
	}
	newTr.BeforeReq = func(req *http.Request) {
		s.HeadersFunc(req)
		s.insertCookie(req)
		s.insertReferer(req)
	}
	s.Client = http.Client{
		Transport:     newTr,
		CheckRedirect: checkRedirect,
		Jar:           jar,
	}
	return s
}
