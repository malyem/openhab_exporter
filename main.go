package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiurl        = kingpin.Flag("apiurl", "http API url to Openhab instance").Required().String()
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9266").String()
)

//Item of items
type Item struct {
	Link       string   `json:"link"`
	State      string   `json:"state"`
	Type       string   `json:"type"`
	Name       string   `json:"name"`
	Label      string   `json:"label"`
	Tags       []string `json:"tags"`
	GroupNames []string `json:"groupNames"`
}

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting openhab_exporter", "version", version.Info())
	level.Info(logger).Log("build_context", version.BuildContext())

	handleCollector(logger)

	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}

func getRestItems() (i []Item, err error) {
	u, err := url.Parse(*apiurl)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "rest/items")

	body, err := getRest(u.String())
	if err != nil {
		return nil, err
	}

	var items []Item
	json.Unmarshal([]byte(body), &items)
	return items, nil
}

func getRest(endpoint string) (resp []byte, err error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Expected http response code 200 got %d, check `apiurl`", response.StatusCode)
	}

	return body, nil
}
