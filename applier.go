package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fireworq/fireworq/model"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type applier interface {
	apply(buf []byte) error
}

type queueApplier struct {
	endpoint string
}

func newQueueApplier(host string) *queueApplier {
	endpoint := fmt.Sprintf("%s/queue", host)
	return &queueApplier{endpoint}
}

func (q *queueApplier) apply(buf []byte) error {
	var queue model.Queue
	if err := json.Unmarshal(buf, &queue); err != nil {
		return err
	}

	r := bytes.NewReader(buf)

	url := fmt.Sprintf("%s/%s", q.endpoint, queue.Name)

	return put(url, r)
}

type routingApplier struct {
	endpoint string
}

func newRoutingApplier(host string) *routingApplier {
	endpoint := fmt.Sprintf("%s/routing", host)
	return &routingApplier{endpoint}
}

func (r *routingApplier) apply(buf []byte) error {
	var routing model.Routing
	if err := json.Unmarshal(buf, &routing); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s", r.endpoint, routing.JobCategory)

	reader := bytes.NewReader(buf)
	return put(url, reader)
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

func apply(fileName string, a applier) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		buf := scanner.Bytes()
		if err := a.apply(buf); err != nil {
			return err
		}
	}

	return nil
}

func applyQueue(fileName, host string) error {
	a := newQueueApplier(host)
	return apply(fileName, a)
}

func applyRouting(fileName, host string) error {
	a := newRoutingApplier(host)
	return apply(fileName, a)
}
