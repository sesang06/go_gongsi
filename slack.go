package main

import (
	"fmt"
	"github.com/slack-go/slack"
)

func SendMessage(message string) {
	api := slack.New("xoxb-1772754878533-1772775182021-29BDNhblX4L7WrKupXIoJEwb")
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// slack.New("YOUR_TOKEN_HERE", slack.OptionDebug(true))

	_, _, err := api.PostMessage("test", slack.MsgOptionText(message, false))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}


//func main() {
//	SendMessage("HEELO WOLRD!")
//}