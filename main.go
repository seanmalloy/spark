// Copyright (c) 2017, Sean Malloy
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
// CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER
// IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT
// OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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

	flag.Usage = func() {
		os.Stderr.WriteString("Commands:\n")
		os.Stderr.WriteString("  spark help [COMMAND]\n")
		os.Stderr.WriteString("  spark msg [OPTION...] MESSAGE\n")
	}

	flag.Parse()

	helpCommand := flag.NewFlagSet("help", flag.ExitOnError)

	helpCommand.Usage = func() {
		os.Stderr.WriteString("Print help messages\n")
		os.Stderr.WriteString("\n")
		os.Stderr.WriteString("SYNOPSIS\n")
		os.Stderr.WriteString("  spark help COMMAND\n")
		os.Stderr.WriteString("\n")
		os.Stderr.WriteString("EXAMPLES\n")
		os.Stderr.WriteString("  $ spark help msg\n")
	}

	msgCommand := flag.NewFlagSet("msg", flag.ExitOnError)
	msgPersonOpt := msgCommand.String("p", "", "send message to `PERSON`")
	msgSpaceOpt := msgCommand.String("s", "", "send message to `SPACE`")
	msgFileOpt := msgCommand.String("f", "", "send `FILE` as an attachment")

	msgCommand.Usage = func() {
		os.Stderr.WriteString("Send message to a person or a space\n")
		os.Stderr.WriteString("\n")
		os.Stderr.WriteString("SYNOPSIS\n")
		os.Stderr.WriteString("  spark msg -p PERSON [-f FILE] MESSAGE\n")
		os.Stderr.WriteString("  spark msg -s SPACE [-f FILE] MESSAGE\n")
		os.Stderr.WriteString("\n")
		os.Stderr.WriteString("OPTIONS\n")
		msgCommand.PrintDefaults()
		os.Stderr.WriteString("\n")
		os.Stderr.WriteString("EXAMPLES\n")
		os.Stderr.WriteString("  $ spark msg -p joe.smith@example.com 'hello world'\n")
	}

	// verify that a sub command has been provided
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	// parse CLI options for each subcommand
	switch os.Args[1] {
	case "help":
		if len(os.Args) <= 2 {
			flag.Usage()
			os.Exit(0)
		}
		helpCommand.Parse(os.Args[2:])
	case "msg":
		msgCommand.Parse(os.Args[2:])
	default:
		flag.Usage()
		os.Exit(1)
	}

	if helpCommand.Parsed() {
		if os.Args[2] == "msg" {
			msgCommand.Usage()
		} else {
			flag.Usage()
			os.Exit(1)
		}
	}

	if msgCommand.Parsed() {
		if msgCommand.Arg(0) == "" {
			// missing message
			os.Stderr.WriteString("Error: missing MESSAGE\n")
			os.Stderr.WriteString("\n")
			msgCommand.Usage()
			os.Exit(1)
		}

		if *msgPersonOpt == "" && *msgSpaceOpt == "" {
			// neither -p or -s were given
			msgCommand.Usage()
			os.Exit(1)
		}
		if *msgPersonOpt != "" && *msgSpaceOpt != "" {
			// -p and -s were both given
			msgCommand.Usage()
			os.Exit(1)
		}

		// -f is optional
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
