package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type method string

// http methods
const (
	post   method = "POST"
	put    method = "PUT"
	delete method = "DELETE"
)

type query struct {
	m    method
	url  string
	body string
}

var client = &http.Client{}

func via(m string) method {
	switch m {
	case "post":
		return post
	case "put":
		return put
	case "delete":
		return delete
	}
	return post
}

func extract(uri string) (query, error) {
	colon := strings.Index(uri, ":")
	if colon == -1 {
		return query{}, errors.New("skip")
	}
	slash := strings.Index(uri[1:], "/") + 1
	question := strings.Index(uri, "?")
	var schema string
	if uri[slash+1:colon] == "http" {
		schema = "http://"
	} else {
		schema = "https://"
	}
	body, _ := url.QueryUnescape(uri[question+1:])
	return query{m: via(uri[1:slash]), url: schema + uri[colon+1:question], body: body}, nil
}

func handle(w http.ResponseWriter, r *http.Request) {
	q, err := extract(r.RequestURI)
	if err != nil {
		return
	}

	req, err := http.NewRequest(string(q.m), q.url, strings.NewReader(q.body))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("User-Agent", "insomnia/2020.3.3")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write(body)
}

func main() {
	http.HandleFunc("/", handle)

	err := http.ListenAndServe("0.0.0.0:1230", nil)
	if err != nil {
		fmt.Println("cannot start server")
	}
}
