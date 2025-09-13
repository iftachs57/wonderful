package main

import (
	"fmt"
	base_actions "wonderful/actions"
	base_connection "wonderful/connection"
)

//Call the basic connection to the server.
//Then, the Session_read_write_actions are entered for sending and receiving input and output from the session.

func main() {
	c, _, err := base_connection.URL_Dailer()
	if err != nil {
		fmt.Printf("error appaned: %v", err)
	}
	base_actions.Session_Read_Write_Actions(c)
}