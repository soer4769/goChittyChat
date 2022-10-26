package main

import (
	"bufio"
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gochittychat/proto"
	"google.golang.org/grpc"
)

var client goChittyChat.ChatServiceClient
var usrId int32 = -1
var lamportTime int64 = 0

func inPost(done chan bool, inStream goChittyChat.ChatService_ConnectClient) {
	for {
		resp, err := inStream.Recv()

		if err == io.EOF {
			done <- true
			return
		}

		if err != nil {
			log.Fatalf("Failed to receive %v", err)
		}

		if usrId < 0 {
			usrId = resp.Id
		}

		if usrId != resp.Id {
			if resp.Lamport > lamportTime {
				lamportTime = resp.Lamport + 1
			} else {
				lamportTime++
			}
		}

		log.Println(strings.Replace(resp.Message, "#L#", strconv.FormatInt(lamportTime, 10), 1))
	}
}

func outPost(outStream goChittyChat.ChatService_MessagesClient) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		textIn := scanner.Text()
		if len(textIn) > 128 {
			log.Println("Failed to send message: more than 128 characters")
			continue
		}

		lamportTime++
		usrIn := goChittyChat.Post{Id: usrId, Message: textIn, Lamport: lamportTime}

		err := outStream.Send(&usrIn)
		if err != nil {
			return
		}

		if usrIn.Message == "--exit" {

			_, err := client.Disconnect(context.Background(), &usrIn)
			if err != nil {
				return
			}
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Failed to scan input: %v", err)
	}
}

func main() {
	// dial server
	conn, connErr := grpc.Dial("localhost:50005", grpc.WithTransportCredentials(insecure.NewCredentials()))
	log.Printf("Client connecting to server...")
	if connErr != nil {
		log.Fatalf("Failed to connect: %v", connErr)
	}
	client = goChittyChat.NewChatServiceClient(conn)

	// setup streams

	// InStream | all incoming messages
	inStream, inErr := client.Connect(context.Background(), &goChittyChat.Post{Lamport: lamportTime})
	if inErr != nil {
		log.Fatalf("Failed to open connection stream: %v", inErr)
	}

	// outStream | all outgoing messages
	outStream, outErr := client.Messages(context.Background())
	if outErr != nil {
		log.Fatalf("Failed to open message stream: %v", outErr)
	}
	log.Println("Client connected successfully...")

	// running goroutines streams
	done := make(chan bool)
	go inPost(done, inStream)
	go outPost(outStream)
	<-done

	// closes server connection

	err := conn.Close()

	if err != nil {
		log.Fatalf("Error when closing connection: %v", err)
	}

	log.Printf("Server connection ended...")
}
