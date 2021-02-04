package main

import (
	"encoding/json"
	"fmt"
	"github.com/fireworq/fireworq/model"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type marshaler interface {
	marshal() ([]byte, error)
}

type myQueue model.Queue
type myRouting model.Routing

func (m myQueue) marshal() ([]byte, error) {
	return json.Marshal(&m)
}

func (m myRouting) marshal() ([]byte, error) {
	return json.Marshal(&m)
}

type retriever interface {
	endpoint(host string) string
	dataMarshal(buf []byte) ([]marshaler, error)
}

type queueRetriever struct {
}

func newQueueRetriever() *queueRetriever {
	return &queueRetriever{}
}

func (q *queueRetriever) endpoint(host string) string {
	return fmt.Sprintf("%s/queues", host)
}

func (q *queueRetriever) dataMarshal(buf []byte) ([]marshaler, error) {
	var mq []myQueue
	if err := json.Unmarshal(buf, &mq); err != nil {
		return []marshaler{}, err
	}

	ret := make([]marshaler, len(mq))
	for i,v := range mq {
		var itf marshaler = v
		ret[i] = itf
	}

	return ret, nil
}

type routingRetriever struct {
}

func newRoutingRetriever() *routingRetriever {
	return &routingRetriever{}
}

func (r *routingRetriever) endpoint(host string) string {
	return fmt.Sprintf("%s/routings", host)
}

func (r *routingRetriever) dataMarshal(buf []byte) ([]marshaler, error) {
	var mr []myRouting
	if err := json.Unmarshal(buf, &mr); err != nil {
		return []marshaler{}, err
	}

	ret := make([]marshaler, len(mr))
	for i,v := range mr {
		var itf marshaler = v
		ret[i] = itf
	}

	return ret, nil
}

func main() {
	host := "http://localhost:18080"
	if err := saveQueueData("./queues.json", host); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := saveRoutingData("./routings.json", host); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func fileSave(fileName string, marshalers []marshaler) error {
	f, err := os.OpenFile(fileName, os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	defer f.Close()

	for _, v := range marshalers {
		buf, err := v.marshal()
		if err != nil {
			return err
		}
		fmt.Fprintln(f, string(buf))
	}

	return nil
}

func getData(url string, unmarshalFunc func([]byte) ([]marshaler, error)) ([]marshaler, error) {
	res, err := http.Get(url)
	if err != nil {
		return []marshaler{}, err
	}

	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []marshaler{}, err
	}

	return unmarshalFunc(buf)
}

func saveData(fileName, host string, r retriever) error {
	url := r.endpoint(host)
	marshalers, err := getData(url, r.dataMarshal)
	if err != nil {
		return err
	}
	return fileSave(fileName, marshalers)
}

func saveQueueData(fileName, host string) error {
	r := newQueueRetriever()
	return saveData(fileName, host, r)
}

func saveRoutingData(fileName, host string) error {
	r := newRoutingRetriever()
	return saveData(fileName, host, r)
}

func put(url string, r io.Reader) error {
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodPut, url, r)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if _, err :=io.Copy(ioutil.Discard, res.Body); err != nil {
		return err
	}

	return nil
}