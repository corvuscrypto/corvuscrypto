package main

import (
	"net/http"
	"net/url"
)

//Header is just a utility function that quickly preps the header template data
func Header(_title string) interface{} {
	return struct {
		Title  string
		Assets map[string]interface{}
	}{
		_title,
		Config.Assets,
	}
}

func Request(r *http.Request) interface{} {
	return struct {
		URL *url.URL
	}{
		r.URL,
	}
}

func BaseData(r *http.Request, title string) map[string]interface{} {
	data := make(map[string]interface{})
	data["Header"] = Header(title)
	data["Request"] = Request(r)
	return data
}
