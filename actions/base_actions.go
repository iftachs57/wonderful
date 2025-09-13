package base_actions

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"bufio"
	"strings"
	"sync"
	"wonderful/configuration"
	"wonderful/structs"

	"github.com/gorilla/websocket"
)

//Getting the configuration file 

func Get_Configuration() (structs.Config, error) {
	file, err := os.Open(consts.Conf_Location)
	if err != nil {
		log.Fatal(err)
		return structs.Config{}, err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	var conf structs.Config
	if err := json.Unmarshal(bytes, &conf); err != nil {
		log.Fatal(err)
	}
	return conf, nil
}

//Takes the session we created and enables sending and receiving messages

func Session_Read_Write_Actions(c *websocket.Conn){
	conf, err := Get_Configuration() 
	if err != nil {
		log.Fatalln("error: ",err) 
	} 
	defer c.Close() 
	answer_returned := true 
	_ = c.SetReadDeadline(time.Now().Add(consts.PongWait)) 
	c.SetPongHandler(func(string) error {
		_ = c.SetReadDeadline(time.Now().Add(consts.PongWait)) 
		return nil 
		}) 
		done := make(chan struct{}) 
		input := make(chan string) 
		go func(){
			r := bufio.NewReader(os.Stdin) 
			for { if answer_returned{
				fmt.Print("\nEnter text: ")
				line, err := r.ReadString('\n')
				fmt.Println("") 
				answer_returned = false 
				if err != nil {  
					close(input) 
					return 
				} 
				line = strings.TrimSpace(line) 
				if line != "" {
					input <- line 
				} 
			} 
		} 
	}() 
	go func() {
		defer close(done) 
		for { 
			_, msg, err := c.ReadMessage() 
			if err != nil {
				log.Println("read:", err) 
				return 
			} 
			var env structs.Envelope 
			if err := json.Unmarshal(msg, &env); err == nil && env.Type == consts.Respons_done {
				var res map[string]any 
				if err := json.Unmarshal(msg, &res); err != nil { 
					fmt.Printf("%s\n", msg) 
					break 
				} 
				t, _ := res["type"].(string) 
				if t != consts.Respons_done{
					break 
				} 
				resp, _ := res["response"].(map[string]any) 
				output, _ := resp["output"].([]any) 
				for _, item := range output {
					im, _ := item.(map[string]any) 
					contents, _ := im["content"].([]any) 
					for _, c := range contents {
						cm, _ := c.(map[string]any) 
						if tr, ok := cm["transcript"].(string); ok && tr != "" {
							fmt.Printf("%s: %s \n",conf.Host_URL,tr) 
							answer_returned = true 
						}  
					} 
				}  
			continue 
			} 
		} 
		}() 
	var mu sync.Mutex 
	writeJSON := func(v any) error {
			mu.Lock() 
			defer mu.Unlock() 
			_ = c.SetWriteDeadline(time.Now().Add(consts.WriteWait)) 
			return c.WriteJSON(v) 
		} 
		writePing := func() error {
			mu.Lock() 
			defer mu.Unlock() 
			_ = c.SetWriteDeadline(time.Now().Add(consts.WriteWait)) 
			return c.WriteMessage(websocket.PingMessage, nil) 
		} 
	ticker := time.NewTicker(consts.PingPeriod) 
	defer ticker.Stop() 
	for {
	select {
		case <-done: 
			return 
		case <-ticker.C: 
			if err := writePing(); err != nil {
				log.Println("ping:", err) 
				return 
			} 
		case line, ok := <-input: 
			if !ok {
				return 
			} 
			msg := structs.EvtConversationItemCreate{
				Type: consts.Conversation_item_create,
				Item: structs.OutMessage{ 
					Type: "message",
					Role: "user", 
					Content: []structs.OutBlock{ 
						{ 
							Type: "input_text", 
							Text: line, 
						}, 
					}, 
				}, 
			} 
			if err := writeJSON(msg); err != nil { 
				log.Println("write message:", err) 
				return 
			} 
			if err := writeJSON(structs.EvtResponseCreate{ 
				Type: consts.Respons_create, 
				Response: struct{}{}, 
			}); err != nil { 
				log.Println("write response.create: ", err) 
				return 
			} 
		} 
	} 
}