package httpmd5

import (
	"crypto/md5"
	"io/ioutil"
	"strings"
	"sync"
)

import (
	"fmt"
	"net/http"
)

type URLMD5 struct {
	URL string
	MD5 *string
	Err *error
}

func (this URLMD5) String() string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("URL: %s", this.URL))

	if this.MD5 == nil {
		result.WriteString(", MD5: nil")
	} else {
		result.WriteString(fmt.Sprintf(", MD5: %v", *this.MD5))
	}

	if this.Err == nil {
		result.WriteString(", Err: nil")
	} else {
		result.WriteString(fmt.Sprintf(", Err: %v", *this.Err))
	}

	return result.String()
}

type HTTPMD5 struct {
	client *http.Client
}

func NewHTTPMD5(client *http.Client) *HTTPMD5 {
	return &HTTPMD5{
		client: client,
	}
}

func (this *HTTPMD5) getURLMD5(url string) URLMD5 {
	url = strings.TrimSpace(url)

	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	resp, err := this.client.Get(url)
	if err != nil {
		return URLMD5{
			URL: url,
			MD5: nil,
			Err: &err,
		}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return URLMD5{
			URL: url,
			MD5: nil,
			Err: &err,
		}
	}

	hash := fmt.Sprintf("%x", md5.Sum(body))

	urlMD5 := URLMD5{
		URL: url,
		MD5: &hash,
		Err: nil,
	}

	return urlMD5
}

func (this *HTTPMD5) GetMD5(parallel int, urls []string) []URLMD5 {
	input := make(chan string)
	urlMD5s := make(chan URLMD5)

	var wg sync.WaitGroup

	wg.Add(parallel)

	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()

			for url := range input {
				urlMD5 := this.getURLMD5(url)

				urlMD5s <- urlMD5
			}
		}()
	}

	go func() {
		for _, url := range urls {
			input <- url
		}

		close(input)
	}()

	go func() {
		wg.Wait()
		close(urlMD5s)
	}()

	urlResults := make([]URLMD5, 0, len(urls))

	for urlResult := range urlMD5s {
		urlResults = append(urlResults, urlResult)
	}

	return urlResults
}
