package main

import "encoding/json"

type WsMessage struct {
	token  string
	url    string
	fields map[string]string
}

func fromJson(encoded []byte) (WsMessage, error) {
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal(encoded, &objmap)
	if err != nil {
		return WsMessage{}, err
	}
	var token string
	err = json.Unmarshal(*objmap["token"], &token)
	if err != nil {
		return WsMessage{}, err
	}
	var url string
	err = json.Unmarshal(*objmap["url"], &url)
	if err != nil {
		return WsMessage{}, err
	}
	fields := make(map[string]string)
	err = json.Unmarshal(*objmap["fields"], &fields)
	if err != nil {
		return WsMessage{}, err
	}

	return WsMessage{token, url, fields}, nil
}
