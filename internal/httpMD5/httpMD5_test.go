package httpmd5

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func makeTestHTTPServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		// Send response to be tested
		reqURL := req.URL.String()
		if reqURL == "/" || strings.HasPrefix(reqURL, "/path") {
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte(reqURL))
			return
		}

		rw.WriteHeader(http.StatusNotFound)
		_, _ = rw.Write([]byte("not found"))
	}))

	return server
}

func makeMD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func TestOneURL(t *testing.T) {
	// Close the server when test finishes
	httpServer := makeTestHTTPServer()
	defer httpServer.Close()

	httpMD5 := NewHTTPMD5(httpServer.Client())
	testURL := httpServer.URL + "/"
	md5s := httpMD5.GetMD5(1, []string{testURL})

	if len(md5s) != 1 {
		t.Error("one answer is expected")
		return
	}

	if md5s[0].URL != testURL {
		t.Errorf("url is not correct: %v", md5s[0].URL)
	}

	if md5s[0].Err != nil || md5s[0].MD5 == nil {
		t.Errorf("error is not expected")
	}

	if *md5s[0].MD5 != makeMD5("/") {
		t.Errorf("md5 is not correct")
	}
}

func TestWrongURL(t *testing.T) {
	// Close the server when test finishes
	httpServer := makeTestHTTPServer()
	defer httpServer.Close()

	httpMD5 := NewHTTPMD5(httpServer.Client())
	testURL := "http://www.not-existing-host-anywhere"
	md5s := httpMD5.GetMD5(1, []string{testURL})

	if len(md5s) != 1 {
		t.Error("one answer is expected")
		return
	}

	if md5s[0].URL != testURL {
		t.Errorf("url is not correct: %v", md5s[0].URL)
	}

	if md5s[0].Err == nil || md5s[0].MD5 != nil {
		t.Errorf("error is expected")
	}
}

func TestManyURLs(t *testing.T) {
	// Close the server when test finishes
	httpServer := makeTestHTTPServer()
	defer httpServer.Close()

	httpMD5 := NewHTTPMD5(httpServer.Client())

	testURLs := []string{}

	const numReqests = 1000

	for i := 0; i < numReqests; i++ {
		testURLs = append(testURLs, fmt.Sprintf("%s%s%d", httpServer.URL, "/path", i))
	}

	md5s := httpMD5.GetMD5(10, testURLs)

	if len(md5s) != numReqests {
		t.Errorf("wrong number of responses: %v", len(md5s))
		return
	}

	for _, urlMD5 := range md5s {
		if urlMD5.Err != nil || urlMD5.MD5 == nil {
			t.Errorf("error is not expected")
		}

		body := strings.TrimPrefix(urlMD5.URL, httpServer.URL)

		bodyMD5 := makeMD5(body)

		if *urlMD5.MD5 != bodyMD5 {
			t.Errorf("md5 is not correct: %v", urlMD5)
			return
		}
	}
}
