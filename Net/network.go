package Net

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var pageTimeout = 3
var UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) snap Chromium/70.0.3538.77 Chrome/70.0.3538.77 Safari/537.36"
var client = http.Client{Timeout: time.Duration(pageTimeout) * time.Second}

var Forbidden403 error = errors.New("Got status code 403 forbidden.")
var Notfound404 error = errors.New("Got status code 404 not found.")

func SetTimeout(t int) {
	pageTimeout = t
	client.Timeout = time.Duration(t) * time.Second
}

func FetchData(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Host", GetDomain(url))
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	if UserAgent != "" {
		req.Header.Add("User-Agent", UserAgent)
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	rep, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rep.Body.Close()
	if rep.StatusCode == 403 {
		return nil, Forbidden403
	}
	if rep.StatusCode == 404 {
		return nil, Notfound404
	}
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func FetchContext(url string) (string, error) {
	data, err := FetchData(url)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Download(url string, filename string) error {
	data, err := FetchData(url)
	if err != nil {
		return err
	}
	fhandle, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fhandle.Close()
	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		return err
	}
	return nil
}
