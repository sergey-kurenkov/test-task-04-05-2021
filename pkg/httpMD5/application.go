package httpmd5

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Application struct {
	out    io.Writer
	client *http.Client
}

func NewApplication(out io.Writer, client *http.Client) *Application {
	return &Application{
		out:    out,
		client: client,
	}
}

func (this *Application) Run(parallel int, timeout time.Duration, urls []string) {
	httpMD5Client := NewHTTPMD5(this.client, timeout)

	urlMD5s := httpMD5Client.GetMD5(parallel, urls)

	for _, urlMD5 := range urlMD5s {
		if urlMD5.MD5 != nil {
			fmt.Fprintf(this.out, "%s %s\n", urlMD5.URL, *urlMD5.MD5)
			continue
		}

		if urlMD5.Err != nil {
			fmt.Fprintf(this.out, "%s, error: %s\n", urlMD5.URL, *urlMD5.Err)
		}
	}
}
