package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiurl = kingpin.Flag("apiurl", "http API url to Openhab instance").Required().String()

	opsQueued = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "openhab",
			Subsystem: "temperature",
			Name:      "current",
			Help:      "Temerature",
		},
		[]string{
			// Which user has requested the operation?
			"user",
			// Of what type is the operation?
			"type",
		},
	)
)

//Items array of Item
type Items struct {
	Collection []Item
}

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

// func main() {
// 	ExampleCollector()
// }

func main() {
	HandleCollector()
}

// func main() {

// 	// prometheus.MustRegister(opsQueued)
// 	prometheus.MustRegister(opsQueued)
// 	kingpin.Parse()
// 	fmt.Printf("%s\n", *apiurl)

// 	itemName := make(chan string)
// 	go doAuthRequest("http://192.168.0.116:8080/rest/items", itemName)

// 	// Increase a value using compact (but order-sensitive!) WithLabelValues().
// 	rand.Seed(time.Now().UnixNano())
// 	for i := range itemName {
// 		fmt.Printf("%+v", i)
// 		opsQueued.WithLabelValues(i, "put").Add(rand.Float64() * 100)
// 	}
// 	// Increase a value with a map using WithLabels. More verbose, but order
// 	// doesn't matter anymore.
// 	opsQueued.With(prometheus.Labels{"type": "delete", "user": "alice"}).Inc()

// 	http.Handle("/metrics", promhttp.Handler())
// 	http.ListenAndServe(":2112", nil)
// }

func doAuthRequest(url string) []Item {
	body := doReq(url)
	var jsonBody []Item
	json.Unmarshal([]byte(body), &jsonBody)
	return jsonBody
	// parseToken(jsonBody.Token)
	// fmt.Printf("%+v", jsonBody)
	// for _, item := range jsonBody {
	// 	c <- item.Name
	// }
	// close(c)
}

func doReq(url string) []byte {
	urlReq := url

	req, err := http.NewRequest("GET", urlReq, nil)
	req.Header.Set("Content-Type", "application/json")
	// req.SetBasicAuth(config.AppConfig.Username, config.AppConfig.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
		log.Panicln(fmt.Sprint("Response Body:", string(body)))
	}

	return body
}
