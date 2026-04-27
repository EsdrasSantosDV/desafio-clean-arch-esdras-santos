package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/EsdrasSantosDV/desafio-clean-arch-esdras-santos/configs"
	"github.com/EsdrasSantosDV/desafio-clean-arch-esdras-santos/internal/event/handler"
	"github.com/EsdrasSantosDV/desafio-clean-arch-esdras-santos/internal/infra/graph"
	"github.com/EsdrasSantosDV/desafio-clean-arch-esdras-santos/internal/infra/grpc/pb"
	"github.com/EsdrasSantosDV/desafio-clean-arch-esdras-santos/internal/infra/grpc/service"
	"github.com/EsdrasSantosDV/desafio-clean-arch-esdras-santos/internal/infra/web/webserver"
	"github.com/EsdrasSantosDV/desafio-clean-arch-esdras-santos/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(cfg.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := waitForDB(db); err != nil {
		panic(err)
	}

	if err := runMigrations(db); err != nil {
		panic(err)
	}

	rabbitMQChannel := getRabbitMQChannel(cfg.RabbitMQUser, cfg.RabbitMQPassword, cfg.RabbitMQHost, cfg.RabbitMQPort)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	webserver := webserver.NewWebServer(cfg.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.Router.Post("/order", webOrderHandler.Create)
	webserver.Router.Get("/order", webOrderHandler.List)
	fmt.Println("Starting web server on port", cfg.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", cfg.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", cfg.GraphQLServerPort)
	http.ListenAndServe(":"+cfg.GraphQLServerPort, nil)
}

func waitForDB(db *sql.DB) error {
	var err error
	for attempt := 1; attempt <= 30; attempt++ {
		if err = db.Ping(); err == nil {
			return nil
		}
		fmt.Printf("Waiting for database (%d/30): %v\n", attempt, err)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("database did not become ready: %w", err)
}

func runMigrations(db *sql.DB) error {
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, file := range files {
		query, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(query)); err != nil {
			return fmt.Errorf("migration %s failed: %w", file, err)
		}
		fmt.Println("Applied migration", file)
	}

	return nil
}

func getRabbitMQChannel(user, password, host, port string) *amqp.Channel {
	var err error
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)

	for attempt := 1; attempt <= 30; attempt++ {
		conn, dialErr := amqp.Dial(dsn)
		if dialErr == nil {
			ch, channelErr := conn.Channel()
			if channelErr == nil {
				return ch
			}
			err = channelErr
			_ = conn.Close()
		} else {
			err = dialErr
		}

		fmt.Printf("Waiting for RabbitMQ (%d/30): %v\n", attempt, err)
		time.Sleep(2 * time.Second)
	}

	panic(fmt.Errorf("rabbitmq did not become ready: %w", err))
}
