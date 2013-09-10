# csession


csession provides a client session for http client. It has the same functions as http.Client, but it will automatically save cookies, headers and set referer. 


## Installing csession

	go get github.com/lunny/csession

## Quick Start

1. Before use please New or NewSession
```Go
import (
	"github.com/lunny/csession"
)
session := csession.New()
```

or

```Go
import (
	"github.com/lunny/csession"
)
session := csession.NewSession(transport, checkRedirect, jar)
```
NewSession's params are the same as Client's fields

2.If you want to customize your headers, your can set HeadersFunc as your func. Default, session.HeadersFunc = session.DefaultHeadersFunc.

```Go
session.HeadersFunc = func(req *http.Request) {
	session.DefaultHeadersFunc(req)
	req.Header.Set("Cache-Control", "max-age=0")
}
```

3.use session like use client

```Go
resp, err := session.Get("http://www.google.com")

forms := url.Values{
		"username": {"username"}, "password": {"password"},
	}


resp, err := session.PostForm("http://www.google.com", forms)

resp, err := session.Post("http://www.google.com")

resp, err := session.Head(...)

resp, err := sesion.Do(req)

```

## Documents 

Please visit [GoWalker](http://gowalker.org/github.com/lunny/csession)


## LICENSE

 BSD License
 [http://creativecommons.org/licenses/BSD/](http://creativecommons.org/licenses/BSD/)
