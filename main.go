package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/jbogarin/go-cisco-spark/ciscospark"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// START: figure out flag sets
	// helpOpt := flag.Bool("h", false, "print help message and exit")
	// flag.Parse()

	// if *helpOpt {
	//		flag.PrintDefaults()
	//		os.Exit(0)
	//}

	//
	// https://blog.komand.com/build-a-simple-cli-tool-with-golang
	//

	msgCommand := flag.NewFlagSet("msg", flag.ExitOnError)

	msgPersonOpt := msgCommand.String("p", "", "send message to a person")

	// verify that a sub command has been provided
	if len(os.Args) < 2 {
		fmt.Println("msg command is required")
		os.Exit(1)
	}

	// parse CLI options for each subcommand
	switch os.Args[1] {
	case "msg":
		msgCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if msgCommand.Parsed() {
		// TODO: sort out required and optional options
		fmt.Printf("msgPersonOpt: %s\n", *msgPersonOpt)
	}

	os.Exit(0)

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
		Text:       "This is a text message",
		ToPersonID: myPersonID,
	}
	newTextMessage, _, err := sparkClient.Messages.Post(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST:", newTextMessage.ID, newTextMessage.Text, newTextMessage.Created)
}
