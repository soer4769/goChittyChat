package main

import (
    "context"
    "log"
    "net"
    "io"
    "fmt"
    "strings"

	"github.com/gochittychat/proto"
	"google.golang.org/grpc"
)

type user struct {
    id int32
    lobby string
    chanCom chan goChittyChat.Post
    chanDone chan bool
}

type server struct {
    goChittyChat.UnimplementedChatServiceServer
    users []*user
}

func (s *server) Broadcast(msg *goChittyChat.Post) {
    for i, u := range s.users {
        if int32(i) != msg.Id && u != nil && u.lobby == s.users[msg.Id].lobby {
            u.chanCom <- *msg
        }
    }
}

func (s *server) Connect(in *goChittyChat.Post, srv goChittyChat.ChatService_ConnectServer) error {
    userData := user{int32(len(s.users)), "default", make(chan goChittyChat.Post), make(chan bool)}
    userMsg := fmt.Sprintf("Participant %v joined Chitty-Chat at Lamport time #L#", userData.id)
    s.users = append(s.users, &userData)

    s.Broadcast(&goChittyChat.Post{Id: userData.id, Lamport: in.Lamport, Message: userMsg})
    srv.Send(&goChittyChat.Post{Id: userData.id, Lamport: in.Lamport, Message: userMsg})
    srv.Send(&goChittyChat.Post{Message: "================================\n                        Welcome to Chitty Chat!\n                   \n                     Commands:\n                     --help            Display Help\n                     --exit            Leave Server\n                     --change-lobby X  Change Lobby\n                    ================================"})

    for {
        select {
            case m := <-userData.chanCom:
                srv.Send(&m)
            case <-userData.chanDone:
                s.users[userData.id] = nil
                return nil
        }
    }
}

func (s *server) Disconnect(context context.Context, in *goChittyChat.Post) (out *goChittyChat.Empty, err error) {
    usr := s.users[in.Id]
    msg := fmt.Sprintf("Participant %v left Chitty-Chat at Lamport time #L#", usr.id)

    log.Printf("Client %v disconnected from the server...", usr.id)
    s.Broadcast(&goChittyChat.Post{Id: usr.id, Lamport: in.Lamport, Message: msg})
    usr.chanDone <- true
    return &goChittyChat.Empty{}, nil
}

func (s *server) Messages(srv goChittyChat.ChatService_MessagesServer) error {
    for {
        resp, err := srv.Recv()

        if err == io.EOF {
            return nil
        }

        if err != nil {
            log.Fatalf("Failed to receive %v", err)
            return nil
        }

        if resp.Message == "--exit" {
            return nil
        }

        usr := s.users[resp.Id]

	    if resp.Message == "--help" {
	        usr.chanCom <-goChittyChat.Post{Message: "================================\n                     Commands:\n                     --help            Display Help\n                     --exit            Leave Server\n                     --change-lobby X  Change Lobby\n                    ================================"}
	        continue
	    }

        if strings.Contains(resp.Message,"--change-lobby ") {
            log.Printf("Client %v changed lobby from '%v' to '%v'", resp.Id, usr.lobby, resp.Message[15:])
            s.Broadcast(&goChittyChat.Post{Id: resp.Id, Message: fmt.Sprintf("Participant %v changed lobby to %v", resp.Id, resp.Message[15:])})
            usr.lobby = resp.Message[15:]
            s.Broadcast(&goChittyChat.Post{Id: resp.Id, Message: fmt.Sprintf("Participant %v joined the lobby...", resp.Id)})
            continue
        }

        log.Printf("Client %v wrote the message: %v", resp.Id, resp.Message)
        s.Broadcast(&goChittyChat.Post{Id: resp.Id, Lamport: resp.Lamport, Message: fmt.Sprintf("Participant %v wrote at Lamport time #L#: %v", resp.Id, resp.Message)})
    }
}

func main() {
    // create listener
    lis, err := net.Listen("tcp", "localhost:50005")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    // create grpc server
    s := grpc.NewServer()
	goChittyChat.RegisterChatServiceServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())

    // launch server
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
