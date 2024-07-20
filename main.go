package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"unicode/utf8"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/server"
)

var (
	address          = flag.String("address", "localhost:8080", "Service address")
	maxTasks         = flag.Int("max", 1000, "Maximum concurrent tasks")
	receivedMessages int64
)

func CountString(ctx context.Context, args []string) (int, error) {
	if len(args) == 0 {
		return 0, jrpc2.Errorf(jrpc2.Code(-32602), "no string provided")
	}

	length := utf8.RuneCountInString(args[0])
	atomic.AddInt64(&receivedMessages, 1)

	return length, nil
}

func Status(ctx context.Context) (string, error) {
	if err := jrpc2.ServerFromContext(ctx).Notify(ctx, "pushback", []string{"hello, friend"}); err != nil {
		return "BAD", err
	}
	return "OK", nil
}

func Alert(ctx context.Context, a map[string]string) error {
	message, ok := a["message"]
	if !ok {
		return jrpc2.Errorf(jrpc2.Code(-32602), "Missing key")
	}
	log.Printf("[ALERT]: %s", message)
	return nil
}

func main() {
	flag.Parse()
	if *address == "" {
		log.Fatal("You must provide a network -address to listen on")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	mux := handler.ServiceMap{
		"StringOperations": handler.Map{
			"CountString": handler.New(CountString),
		},
		"Post": handler.Map{
			"Alert": handler.New(Alert),
		},
	}

	lst, err := net.Listen(jrpc2.Network(*address))
	if err != nil {
		log.Fatalln("Listen:", err)
	}
	defer lst.Close()

	log.Printf("Listening at %v...", lst.Addr())
	acc := server.NetAccepter(lst, channel.Line)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		server.Loop(ctx, acc, server.Static(mux), &server.LoopOptions{
			ServerOptions: &jrpc2.ServerOptions{
				Logger:      jrpc2.StdLogger(nil),
				Concurrency: *maxTasks,
				AllowPush:   true,
			},
		})
	}()

	<-sig
	log.Println("Server shutting down...")
	cancel()

	if err := lst.Close(); err != nil {
		log.Fatalf("Error closing listener: %v", err)
	}

	log.Printf("Messages received = %d", receivedMessages)
	log.Println("Server gracefully stopped")
}
