package main

import (
    "context"
    "log"
    "net"
    "time"
    "math/rand"

	"github.com/gonetwork/proto"
	"google.golang.org/grpc"
)

type Flags struct {
    SYN, ACK, FIN bool
}

type Pack struct {
    SeqNum, AckNum uint32
    Message string
    Status Flags
}

type server struct {
	TCPHandshake.UnimplementedHandshakeServer
}

func (s *server) ConnSend(ctx context.Context, in *TCPHandshake.TCPPack) (*TCPHandshake.TCPPack, error) {
    if in.Status.SYN {
        log.Printf("New client trying to establish simulated TCP connection...")
        log.Printf("Recieved message from client:\n\t%+v\n", in)

        ack := Pack{SeqNum: rand.Uint32(), AckNum: in.SeqNum+1, Message: "SYN+ACK", Status: Flags{SYN: true, ACK: true}}
        log.Printf("Sending response to client:\n\t%+v\n", ack)

        return &TCPHandshake.TCPPack{
                 SeqNum: ack.SeqNum,
                 AckNum: ack.AckNum,
                 Message: ack.Message,
                 Status: &TCPHandshake.Flags{
                   SYN: ack.Status.SYN,
                   ACK: ack.Status.ACK,
                 },
               }, nil
    }

    if in.Status.ACK {
        log.Printf("Established simulated TCP connection succefully with client...")
    }

    return &TCPHandshake.TCPPack{}, nil
}

func main() {
    rand.Seed(time.Now().UnixNano())
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()

	TCPHandshake.RegisterHandshakeServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
