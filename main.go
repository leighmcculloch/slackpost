package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	var webhook, username, iconURL, channel, text string
	flag.StringVar(&webhook, "w", "", "slack webhook url")
	flag.StringVar(&username, "u", "", "username to post as")
	flag.StringVar(&iconURL, "i", "", "icon to post as")
	flag.StringVar(&channel, "c", "", "slack channel to post to")
	flag.StringVar(&text, "t", "", "text to post")
	flag.Parse()

	if webhook == "" || username == "" || iconURL == "" || channel == "" || text == "" {
		fmt.Fprintln(os.Stderr, "not all options provided")
		flag.Usage()
		os.Exit(1)
	}

	type Request struct {
		Username string `json:"username"`
		IconURL  string `json:"icon_url"`
		Channel  string `json:"channel"`
		Text     string `json:"text"`
	}

	r := Request{
		Username: username,
		IconURL:  iconURL,
		Channel:  channel,
		Text:     text,
	}

	reqBytes, err := json.Marshal(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error encoding request:", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", webhook, bytes.NewReader(reqBytes))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error preparing request:", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error performing request:", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading response:", err)
		os.Exit(1)
	}

	fmt.Println(string(resBytes))
}
