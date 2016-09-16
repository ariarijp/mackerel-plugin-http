package main

import (
	"flag"
	"log"
	"time"

	"github.com/ddliu/go-httpclient"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

type HttpPlugin struct {
	URL string
}

func (h HttpPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	func() {
		defer func(start time.Time) {
			stat["msec"] = uint64(time.Since(start) / time.Millisecond)
		}(time.Now())

		_, err := httpclient.Get(h.URL, map[string]string{})
		if err != nil {
			log.Fatal(err)
		}
	}()

	return stat, nil
}

func (h HttpPlugin) GraphDefinition() map[string](mp.Graphs) {
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
	optUrl := flag.String("url", "http://localhost/", "URL")
	optTempfile := flag.String("tempfile", "/tmp/mackerel-plugin-http", "Temp file name")
	flag.Parse()

	var h HttpPlugin
	h.URL = *optUrl

	helper := mp.NewMackerelPlugin(h)
	helper.Tempfile = *optTempfile

	helper.Run()
}
