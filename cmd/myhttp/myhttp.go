package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/sergey-kurenkov/test-task-http-md5/internal/httpmd5"
)

func main() {
	parallelFlag := flag.Int("parallel", 10, "number of parallel requests")
	flag.Parse()

	app := httpmd5.NewApplication(os.Stdout, &http.Client{})
	app.Run(*parallelFlag, flag.Args())
}
