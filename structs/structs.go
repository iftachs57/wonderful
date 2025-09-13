package structs

import(
	"encoding/json"
)

type Config struct {
	URL_Header        map[string]string `json:"url_header"`
	Host_URL          string `json:"host_url"`
	Scheme_URL        string `json:"scheme_url"`
	Path_URL          string `json:"path_url"`
	Model			  string `json:"model"`
	EnableCompression bool   `json:"enablecompression"`
}

type OutBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type OutMessage struct {
	Type    string     `json:"type"`   
	Role    string     `json:"role"`    
	Content []OutBlock `json:"content"` 
}

type EvtConversationItemCreate struct {
	Type string     `json:"type"` 
	Item OutMessage `json:"item"`
}

type EvtResponseCreate struct {
	Type     string      `json:"type"`  
	Response interface{} `json:"response"` 
}

type Envelope struct {
	Type  string          `json:"type"`
	Delta json.RawMessage `json:"delta,omitempty"`
	Error *struct {
		Type    string `json:"type"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Param   any    `json:"param"`
	} `json:"error,omitempty"`
}