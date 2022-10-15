package main

import (
    "context"
    "log"
    "time"
    "math/rand"
	"google.golang.org/grpc/credentials/insecure"

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

func SendMessage(c TCPHandshake.HandshakeClient, p Pack) (*TCPHandshake.TCPPack, error) {
    r, err := c.ConnSend(context.Background(), &TCPHandshake.TCPPack{
        SeqNum: p.SeqNum,
        AckNum: p.AckNum,
        Message: p.Message,
        Status: &TCPHandshake.Flags{
            SYN: p.Status.SYN,
            ACK: p.Status.ACK,
            FIN: p.Status.FIN,
        },
    })

    if err != nil {
        log.Fatalf("could not handshake: %v", err)
    }

    return r, err
}

func Shake(c TCPHandshake.HandshakeClient) {
    log.Printf("Establishing Simulated TCP connection with server...")

    syn := Pack{SeqNum: rand.Uint32(), Message: "SYN", Status: Flags{SYN: true}}
    log.Printf("Sending message to server:\n\t%+v\n", syn)

    r, err := SendMessage(c,syn)
    if err != nil { return }
    log.Printf("Recieved message from server:\n\t%+v\n", r)

    ack := Pack{AckNum: r.SeqNum+1, Message: "ACK", Status: Flags{ACK: true}}
    log.Printf("sending message to server:\n\t%+v\n", ack)

    r, err = SendMessage(c,ack)
    if err != nil { return }

    log.Printf("Simulated TCP handshake successfully connected...")
}

func main() {
    rand.Seed(time.Now().UnixNano())
    conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    c := TCPHandshake.NewHandshakeClient(conn)

	Shake(c)
	err = conn.Close()
	if err != nil {
		return
	}
}
