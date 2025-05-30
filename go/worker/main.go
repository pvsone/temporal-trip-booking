package main

import (
	"log"
	"log/slog"
	"os"

	"temporal-trip-booking/activities"
	"temporal-trip-booking/workflows"

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

	w.RegisterWorkflow(workflows.BookWorkflow)
	w.RegisterActivity(activities.BookFlight)
	w.RegisterActivity(activities.BookHotel)
	w.RegisterActivity(activities.BookCar)
	w.RegisterActivity(activities.NotifyUser)
	w.RegisterActivity(activities.UndoBookFlight)
	w.RegisterActivity(activities.UndoBookHotel)
	w.RegisterActivity(activities.UndoBookCar)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
