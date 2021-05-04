package httpmd5

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestApplicationOneURL(t *testing.T) {
	// Close the server when test finishes
	httpServer := makeTestHTTPServer()
	defer httpServer.Close()

	buffer := &bytes.Buffer{}

	app := NewApplication(buffer, httpServer.Client())

	testURL := httpServer.URL + "/"
	app.Run(10, []string{testURL})

	expectedAnswer := fmt.Sprintf("%s%s %s\n", httpServer.URL, "/", makeMD5("/"))

	receivedAnswer := buffer.String()

	if expectedAnswer != receivedAnswer {
		t.Errorf("wrong result:\nexpected answer:%s\nreceived answer:%s",
			expectedAnswer, receivedAnswer)
	}
}

func TestApplicationWrongURL(t *testing.T) {
	// Close the server when test finishes
	httpServer := makeTestHTTPServer()
	defer httpServer.Close()

	buffer := &bytes.Buffer{}

	app := NewApplication(buffer, httpServer.Client())

	testURL := "http://www.not-existing-host-anywhere"
	app.Run(10, []string{testURL})

	receivedAnswer := buffer.String()

	if !strings.HasPrefix(receivedAnswer, 	fmt.Sprintf("%s, error: ", testURL)) {
		t.Errorf("wrong result: expected error: %s", receivedAnswer)
	}
}
