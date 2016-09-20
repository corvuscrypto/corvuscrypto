package main

func Header(_title string) interface{} {
	return struct {
		Title  string
		Assets map[string]interface{}
	}{
		_title,
		Config.Assets,
	}
}
