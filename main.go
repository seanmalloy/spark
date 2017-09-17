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
	//
	// https://blog.komand.com/build-a-simple-cli-tool-with-golang
	//

	helpCommand := flag.NewFlagSet("help", flag.ExitOnError)

	msgCommand := flag.NewFlagSet("msg", flag.ExitOnError)
	msgPersonOpt := msgCommand.String("p", "", "send message to a person")
	msgSpaceOpt := msgCommand.String("s", "", "send message to a space")
	msgFileOpt := msgCommand.String("f", "", "send a file attachment")

	// verify that a sub command has been provided
	if len(os.Args) < 2 {
		fmt.Println("msg command is required")
		os.Exit(1)
	}

	// parse CLI options for each subcommand
	switch os.Args[1] {
	case "help":
		helpCommand.Parse(os.Args[2:])
	case "msg":
		msgCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if helpCommand.Parsed() {

		if os.Args[2] == "msg" {
			msgCommand.PrintDefaults()
		}
	}

	if msgCommand.Parsed() {
		//
		// msgPersonOpt
		// msgSpaceOpt
		if *msgPersonOpt == "" && *msgSpaceOpt == "" {
			// neither -p or -s were given
			msgCommand.PrintDefaults()
			os.Exit(1)
		}
		if *msgPersonOpt != "" && *msgSpaceOpt != "" {
			// -p and -s were both given
			msgCommand.PrintDefaults()
			os.Exit(1)
		}

		// -f is optional
		//
		// msgFileOpt
		if *msgFileOpt == "" {
			if *msgPersonOpt != "" {
				// send message to a person
				fmt.Println("send message to a person")
			}

			if *msgSpaceOpt != "" {
				// send message to a space
				fmt.Println("send message to a space")
			}
		} else {
			if *msgPersonOpt != "" {
				// send file to a person
				fmt.Println("send file to a person")
			}

			if *msgSpaceOpt != "" {
				// send file to a space
				fmt.Println("send file to a space")
			}
		}

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
