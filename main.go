package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"bufio"
	"log"
	"strings"
	"net/http"
	"github.com/jbogarin/go-cisco-spark/ciscospark"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	sparkClient := ciscospark.NewClient(client)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Auth Token: ")
	token, _ := reader.ReadString('\n')
	token = strings.TrimSuffix(token, "\n")
	sparkClient.Authorization = "Bearer " + token

	myPersonID := "722bb271-d7ca-4bce-a9e3-471e4412fa77"

	// POST messages - Text Message
	message := &ciscospark.MessageRequest{
		Text:   "This is a text message",
		ToPersonID: myPersonID,
	}
	newTextMessage, _, err := sparkClient.Messages.Post(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST:", newTextMessage.ID, newTextMessage.Text, newTextMessage.Created)
}
