package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/ddliu/go-httpclient"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

type HTTPPlugin struct {
	URL url.URL
}

func (h HTTPPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	func() {
		defer func(start time.Time) {
			stat["msec"] = uint64(time.Since(start) / time.Millisecond)
		}(time.Now())

		_, err := httpclient.Get(h.URL.String(), map[string]string{})
		if err != nil {
			log.Fatal(err)
		}
	}()

	return stat, nil
}

func (h HTTPPlugin) GraphDefinition() map[string](mp.Graphs) {
	return map[string](mp.Graphs){
		"http.response_time": mp.Graphs{
			Label: "HTTP Response",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				mp.Metrics{
					Name:  "msec",
					Label: "Milliseconds",
				},
			},
		},
	}
}

func main() {
	optURL := flag.String("url", "http://localhost/", "URL")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	u, err := url.Parse(*optURL)
	if err != nil {
		log.Fatal(err)
	}

	if *optTempfile == "" {
		sh := sha1.New()
		io.WriteString(sh, u.String())
		*optTempfile = fmt.Sprintf("/tmp/mackerel-plugin-http-%x", sh.Sum(nil))
	}

	var h HTTPPlugin
	h.URL = *u

	helper := mp.NewMackerelPlugin(h)
	helper.Tempfile = *optTempfile

	helper.Run()
}
