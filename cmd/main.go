package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IskenT/MultiGameServices/cmd/providers"
	"github.com/IskenT/MultiGameServices/configs"

	grpcrepository "github.com/IskenT/MultiGameServices/internal/infrastructure/repository/grpcrepo"
	repository "github.com/IskenT/MultiGameServices/internal/infrastructure/repository/restrepo"
	handler "github.com/IskenT/MultiGameServices/internal/infrastructure/transport/http"
	grpcservice "github.com/IskenT/MultiGameServices/internal/infrastructure/usecase/grpcservice"
	restservice "github.com/IskenT/MultiGameServices/internal/infrastructure/usecase/restservice"

	broadcast "github.com/IskenT/MultiGameServices/internal/delivery/websocket"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var wg sync.WaitGroup
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	cnf, err := configs.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := providers.ProvideConsoleLogger(cnf)
	if err != nil {
		return fmt.Errorf("failed to provide console logger: %w", err)
	}

	cache, err := providers.ProvideCache()
	if err != nil {
		return fmt.Errorf("failed to provide cache: %w", err)
	}

	var errc = make(chan error, 1)
	restWalletRepository := repository.NewWalletRepository(ctx, cache, logger)
	restWalletService := restservice.NewWalletService(ctx, restWalletRepository, logger)
	restWalletHandler := handler.NewWalletService(restWalletService, logger)

	echoServer := providers.ProvideHTTPServer(cnf, restWalletHandler, logger, errc)

	grpcWalletRepository := grpcrepository.NewWalletRepository(ctx, cache, logger)
	grpcWalletService := grpcservice.NewWalletService(ctx, grpcWalletRepository, logger)

	grpcServer := providers.ProvideGRPCServer(cnf, grpcWalletService, logger, errc)

	var wgDone = make(chan bool)
	wg.Add(3)

	go func() {
		defer wg.Done()
		echoServer.Start()
	}()

	go func() {
		defer wg.Done()
		grpcServer.Start()
	}()

	go func() {
		defer wg.Done()

		addr := flag.String("addr", ":8081", "http service address")
		hub := broadcast.NewHub()
		go hub.Run()
		http.HandleFunc("/", serveHome)
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			broadcast.ServeWs(hub, w, r)
		})
		server := &http.Server{
			Addr:              *addr,
			ReadHeaderTimeout: 3 * time.Second,
		}
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	select {
	case <-wgDone:
		break
	case err := <-errc:
		return err
	}

	go func() {
		wg.Wait()
		close(wgDone)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		// context cancelled
	case sig := <-c:
		// signal received
		log.Printf("received signal: %v\n", sig)
		cancel()
	}

	grpcServer.Stop(ctx)
	log.Println("gRPC server stopped gracefully")

	echoServer.Stop(ctx)
	log.Println("echo server stopped gracefully")

	return nil
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}
