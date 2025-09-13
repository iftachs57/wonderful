package base_connection

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"wonderful/actions"
	"wonderful/configuration"

	"github.com/gorilla/websocket"
)

//Basic connection to the server
//Creates a session using the headers from the conf.json file 

func URL_Dailer()(*websocket.Conn, *http.Response, error){
	conf, err := base_actions.Get_Configuration()
	if err != nil {
		log.Fatalf("configuration/configuration_location error: %v ",err)
		return nil, nil, err
	}
	u := url.URL{Scheme: conf.Scheme_URL, Host: conf.Host_URL, Path: conf.Path_URL, RawQuery: consts.URL_RawQuery + url.QueryEscape(conf.Model),}
	dialer := websocket.Dialer{
		HandshakeTimeout:  consts.HandshakeTimeout,
		EnableCompression: conf.EnableCompression,
		ReadBufferSize:    consts.ReadBufferSize,
		WriteBufferSize:   consts.WriteBufferSize,
	}
	headers := http.Header{}
	for header := range conf.URL_Header{
		headers.Add(header,conf.URL_Header[header])
	}
	ctx := context.TODO()
	conn, resp, err := dialer.DialContext(ctx, u.String(), headers)
	if err != nil {
		log.Fatalf("dial error: %v (HTTP status: %v)", err, resp)
		return conn, resp, err
	}
	return conn, resp, nil
}
