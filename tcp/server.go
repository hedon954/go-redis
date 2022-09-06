package tcp

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go-redis/interface/tcp"
	"go-redis/lib/logger"
)

type Config struct {
	Address string
}

// ListenAndServeWithSignal binds port and handle requests, blocking until receive stop signal
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	logger.Info("start to listen", cfg.Address)

	closeChan := make(chan struct{})
	go handleSystemSignal(closeChan)

	ListenAndServe(listener, handler, closeChan)

	return nil
}

// ListenAndServe listens and serves connection
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {

	// system force exit
	go func() {
		<-closeChan
		logger.Info("shutting down")
		_ = listener.Close()
		_ = handler.Close()
	}()

	// close listener and handler
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	// records how many connections right now
	waitDone := sync.WaitGroup{}
	ctx := context.Background()

	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accepted connection", conn.RemoteAddr())
		// add a new connection
		waitDone.Add(1)
		go func() {
			defer waitDone.Done()
			handler.Handle(ctx, conn)
		}()

	}

	// wait for all connections to close
	waitDone.Done()
}

// handleSystemSignal handles the operator system signal
func handleSystemSignal(closeChan chan struct{}) {
	sigChan := make(chan os.Signal)

	// see https://blog.csdn.net/secretii/article/details/118342752
	// SIGHUP hand up
	// SIGQUIT error
	// SIGTERM kill
	// SIGINT ctrl + c
	// SIGKILL kill -9
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	sig := <-sigChan

	switch sig {
	case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
		closeChan <- struct{}{}
	}
}
