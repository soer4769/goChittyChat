package main

import (
    "context"
    "io"
    "log"
    "os"
    "bufio"
    "fmt"

    "github.com/gochittychat/proto"
    "google.golang.org/grpc"
)

var client goChittyChat.ChatServiceClient
var usrId int32 = -1

func inPost(done chan bool, inStream goChittyChat.ChatService_ConnectClient) {
    for {
        resp, err := inStream.Recv()
        if err == io.EOF {
            done <- true
            return
        }

        if err != nil {
            log.Fatalf("cannot receive %v", err)
        }

        if usrId < 0 {
            usrId = resp.Id
            continue
        }

        log.Println(resp.Message)
    }
}

func outPost(done chan bool, outStream goChittyChat.ChatService_MessagesClient) {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        usrIn := goChittyChat.Post{Id: usrId, Message: scanner.Text()}
        outStream.Send(&usrIn)

        if usrIn.Message == "exit" {
            client.Disconnect(context.Background(), &usrIn)
            return
        }
    }

    if err := scanner.Err(); err != nil {
       log.Printf("ERROR: Failed to scan input: %v", err)
    }
}

func main() {
    // dial server
    conn, err := grpc.Dial("localhost:50005", grpc.WithInsecure())
    log.Printf("Client connecting to server...")
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }

    // setup stream
    client = goChittyChat.NewChatServiceClient(conn)
    inStream, inErr := client.Connect(context.Background(), &goChittyChat.ConPost{Username: "test", Lobby: ""})
    outStream, outErr := client.Messages(context.Background())
    if inErr != nil && outErr != nil {
        log.Fatalf("open stream error: %v; %v", inErr, outErr)
    }
    log.Println("Client connected successfully...")
    fmt.Println("-------------------------")

    // running goroutines streams
    done := make(chan bool)
    go inPost(done, inStream)
    go outPost(done, outStream)
    <-done

    // closes server connection
    fmt.Println("-------------------------")
    log.Printf("Closing server connection...")
    conn.Close()
    log.Printf("server connection ended...")
}
