package net

import (
	"fmt"
	"log"
	"os"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/vocdoni/go-dvote/types"
)

type PubSubHandle struct {
	c *types.Connection
	s *shell.PubSubSubscription
}

func PsSubscribe(topic string) *shell.PubSubSubscription {
	sh := shell.NewShell("localhost:5001")
	sub, err := sh.PubSubSubscribe(topic)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	return sub
}

func PsPublish(topic, data string) error {
	sh := shell.NewShell("localhost:5001")
	err := sh.PubSubPublish(topic, data)
	if err != nil {
		return err
	}
	return nil
}

func (p *PubSubHandle) Init(c *types.Connection) error {
	p.c = c
	p.s = PsSubscribe(p.c.Topic)
	return nil
}

func (p *PubSubHandle) Listen(reciever chan<- types.Message) {
	var psMessage *shell.Message
	var msg types.Message
	var err error
	for {
		psMessage, err = p.s.Next()
		if err != nil {
			log.Printf("PubSub recieve error: %s", err)
		}
		ctx := new(types.PubSubContext)
		ctx.Topic = p.c.Topic
		ctx.PeerAddress = psMessage.From.String()
		msg.Data = psMessage.Data
		msg.TimeStamp = time.Now()
		msg.Context = ctx

		reciever <- msg
	}
}

func (p *PubSubHandle) Send(msg types.Message) {
	err := PsPublish(p.c.Topic, string(msg.Data))
	if err != nil {
		log.Printf("PubSub send error: %s", err)
	}
}
