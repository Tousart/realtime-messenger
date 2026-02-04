package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"router/internal/server"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	serverAddr    string
	rabbitMQAddr  string
	messagesQueue string
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := Config{}

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("no env params, using os environments")
		cfg.serverAddr = os.Getenv("SERVER_ADDR")
		cfg.rabbitMQAddr = os.Getenv("RABBITMQ_ADDR_DOCKER")
		cfg.messagesQueue = os.Getenv("MESSAGES_QUEUE")
	} else {
		cfg.serverAddr = os.Getenv("SERVER_ADDR")
		cfg.rabbitMQAddr = os.Getenv("RABBITMQ_ADDR_LOCALHOST")
		cfg.messagesQueue = os.Getenv("MESSAGES_QUEUE")
	}

	go consumeMessages(ctx)

	r := chi.NewRouter()

	// create and run server
	srv := server.NewServer(cfg.serverAddr, r)
	srv.CreateAndRunServer(ctx)

	srv.Wg.Wait()
}

func consumeMessages(ctx context.Context) {
	// connection
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}
	defer conn.Close()

	// create channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Ошибка создания канала:", err)
	}
	defer ch.Close()

	// queue
	q, err := ch.QueueDeclare(
		"messages", // имя очереди
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // аргументы
	)
	if err != nil {
		log.Fatal("Ошибка объявления очереди:", err)
	}

	// register consumer
	msgs, err := ch.Consume(
		q.Name, // имя очереди
		"",     // consumer tag
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal("Ошибка регистрации консьюмера:", err)
	}

	for {
		select {
		case msg := <-msgs:
			log.Printf("Получено сообщение: %s", msg.Body)
		case <-ctx.Done():
			return
		}
	}
}
