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

var userId int32

type user struct {
    id int32
    username string
    lobby string
    chanCom chan goChittyChat.Post
    chanDone chan bool
}

type server struct {
    goChittyChat.UnimplementedChatServiceServer
    users []*user
}

func (s *server) FindUser(id int32) *user {
    for _, u := range s.users {
        if u.id == id {
            return u
        }
    }
    return nil
}

func (s *server) RemoveUser(id int32) {
    for i, u := range s.users {
        if u.id == id {
            s.users = append(s.users[:i], s.users[i+1:]...)
            return
        }
    }
}

func (s *server) BroadcastPost(usr *user, msg *goChittyChat.Post) {
    for _, u := range s.users {
        if u != usr  && u.lobby == usr.lobby {
            u.chanCom <- *msg
        }
    }
}

func (s *server) Connect(in *goChittyChat.ConPost, srv goChittyChat.ChatService_ConnectServer) error {
    userId++
    userData := user{userId, in.Username, "default", make(chan goChittyChat.Post), make(chan bool)}
    s.users = append(s.users, &userData)

    s.BroadcastPost(&userData, &goChittyChat.Post{Id: userData.id, Message: fmt.Sprintf("client %v (%v) connected to server...", userData.id, userData.username)})
    srv.Send(&goChittyChat.Post{Id: userData.id, Message: "Server connection established..."})
    log.Printf("Client %v (%v) connected to server...", userData.id, userData.username)

    for {
        select {
            case m := <-userData.chanCom:
                srv.Send(&m)
            case <-userData.chanDone:
                s.RemoveUser(userData.id)
                return nil
        }
    }
}

func (s *server) Disconnect(context context.Context, in *goChittyChat.Post) (out *goChittyChat.Empty, err error) {
    usr := s.FindUser(in.Id)
    log.Printf("client %v (%v) disconnected...", usr.id, usr.username)
    s.BroadcastPost(usr, &goChittyChat.Post{Id: usr.id, Message: fmt.Sprintf("client %v (%v) disconnected...", usr.id, usr.username)})
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
            log.Fatalf("ERROR: Cannot receive %v", err)
            return nil
        }

        if resp.Message == "exit" {
            return nil
        }

        usr := s.FindUser(resp.Id)

        if strings.Contains(resp.Message,"--change-lobby ") {
            log.Printf("Client %v (%v) changed lobby from '%v' to %v", resp.Id, usr.username, usr.lobby, resp.Message[15:])
            s.BroadcastPost(usr, &goChittyChat.Post{Id: resp.Id, Message: fmt.Sprintf("client %v (%v) changed lobby to %v", usr.username, resp.Id, resp.Message[15:])})
            usr.lobby = resp.Message[15:]
            s.BroadcastPost(usr, &goChittyChat.Post{Id: resp.Id, Message: fmt.Sprintf("client %v (%v) joined the lobby...", usr.username, resp.Id)})
            continue
        }

        log.Printf("Client %v (%v) sent message: %v\n", resp.Id, usr.username, resp.Message)
        s.BroadcastPost(usr, &goChittyChat.Post{Id: resp.Id, Message: fmt.Sprintf("%v (%v): %v", usr.username, resp.Id, resp.Message)})
    }
}

func main() {
    // create listener
    lis, err := net.Listen("tcp", "localhost:50005")
    if err != nil {
        log.Fatalf("ERROR: Failed to listen: %v", err)
    }

    // create grpc server
    s := grpc.NewServer()
	goChittyChat.RegisterChatServiceServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())

    // launch server
    if err := s.Serve(lis); err != nil {
        log.Fatalf("ERROR: Failed to serve: %v", err)
    }
}
