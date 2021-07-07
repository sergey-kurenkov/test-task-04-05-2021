package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/sergey-kurenkov/test-task-http-md5/internal/httpmd5"
)

func main() {
	parallelFlag := flag.Int("parallel", 10, "number of parallel requests")
	timeout := flag.Duration("timeout", time.Second, "maximum timeout")
	flag.Parse()

	app := httpmd5.NewApplication(os.Stdout, &http.Client{})
	app.Run(*parallelFlag, *timeout, flag.Args())
}
