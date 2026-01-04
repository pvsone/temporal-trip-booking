package main

import (
	"log"
	"log/slog"
	"os"
	"temporal-trip-booking/app"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/envconfig"
	tlog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	log.Printf("⚙️ Using TEMPORAL_PROFILE: '%s'", os.Getenv("TEMPORAL_PROFILE"))
	clientOptions := envconfig.MustLoadDefaultClientOptions()
	clientOptions.Logger = tlog.NewStructuredLogger(logger)

	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	log.Printf("✅ Client connected to '%v' in namespace '%v'", clientOptions.HostPort, clientOptions.Namespace)
	defer temporalClient.Close()

	w := worker.New(temporalClient, "trip-task-queue", worker.Options{})

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
