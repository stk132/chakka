package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

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
