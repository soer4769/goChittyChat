package main

import (
    "context"
    "io"
    "log"
    "os"
    "bufio"
    "strings"
    "strconv"

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

        log.Println(strings.Replace(resp.Message, "#L#", strconv.FormatInt(lamportTime,10), 1))
    }
}

func outPost(done chan bool, outStream goChittyChat.ChatService_MessagesClient) {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        textIn := scanner.Text()
        if len(textIn) > 128 {
            log.Println("Failed to send message: more than 128 characters")
            continue
        }

        lamportTime++
        usrIn := goChittyChat.Post{Id: usrId, Message: textIn, Lamport: lamportTime}
        outStream.Send(&usrIn)

        if usrIn.Message == "--exit" {
            client.Disconnect(context.Background(), &usrIn)
            return
        }
    }

    if err := scanner.Err(); err != nil {
       log.Printf("Failed to scan input: %v", err)
    }
}

func main() {
    // dial server
    conn, err := grpc.Dial("localhost:50005", grpc.WithInsecure())
    log.Printf("Client connecting to server...")
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    client = goChittyChat.NewChatServiceClient(conn)

    // setup streams
    inStream, inErr := client.Connect(context.Background(), &goChittyChat.Post{Lamport: lamportTime})
    if inErr != nil {
        log.Fatalf("Failed to open connection stream: %v; %v", inErr)
    }
    outStream, outErr := client.Messages(context.Background())
    if outErr != nil {
        log.Fatalf("Failed to open message stream: %v", outErr)
    }
    log.Println("Client connected successfully...")

    // running goroutines streams
    done := make(chan bool)
    go inPost(done, inStream)
    go outPost(done, outStream)
    <-done

    // closes server connection
    log.Printf("Closing server connection...")
    conn.Close()
    log.Printf("Server connection ended...")
}
