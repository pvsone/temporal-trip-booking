package main

import (
	"log"
	"log/slog"
	"os"
	"temporal-trip-booking/app"

	"go.temporal.io/sdk/client"
	tlog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

const TASK_QUEUE = "trip-task-queue"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	clientOptions := client.Options{
		HostPort:  "localhost:7233",
		Namespace: "default",
		Logger:    tlog.NewStructuredLogger(logger),
	}

	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporalClient.Close()

	w := worker.New(temporalClient, TASK_QUEUE, worker.Options{})

	w.RegisterWorkflow(app.BookWorkflow)
	w.RegisterActivity(app.BookFlight)
	w.RegisterActivity(app.BookHotel)
	w.RegisterActivity(app.BookCar)
	w.RegisterActivity(app.NotifyUser)
	w.RegisterActivity(app.UndoBookFlight)
	w.RegisterActivity(app.UndoBookHotel)
	w.RegisterActivity(app.UndoBookCar)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
