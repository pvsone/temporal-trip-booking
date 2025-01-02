package workflows

import (
	"fmt"
	"temporal-trip-booking/activities"
	"temporal-trip-booking/app"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type BookTripInput struct {
	UserId   string `json:"userId"`
	FlightId string `json:"flightId"`
	HotelId  string `json:"hotelId"`
	CarId    string `json:"carId"`
}

func BookWorkflow(ctx workflow.Context, input BookTripInput) (output string, err error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Book workflow started", "userId", input.UserId)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    30 * time.Second,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Create saga to manage compensations
	var saga app.Saga
	defer func() {
		if err != nil {
			disconnectedCtx, _ := workflow.NewDisconnectedContext(ctx)
			saga.Compensate(disconnectedCtx)
		}
		err = nil
	}()

	// Book Flight
	var flight string
	err = workflow.ExecuteActivity(ctx, activities.BookFlight, input.FlightId).Get(ctx, &flight)
	if err != nil {
		logger.Warn("Flight booking failed", "error", err)
		return "Booking cancelled", err
	}
	saga.AddCompensation(activities.UndoBookFlight, input.FlightId)

	// Simulate delay
	logger.Info("Sleeping for 1 second...")
	workflow.Sleep(ctx, time.Second)

	// Book Hotel
	var hotel string
	err = workflow.ExecuteActivity(ctx, activities.BookHotel, input.HotelId).Get(ctx, &hotel)
	if err != nil {
		logger.Warn("Hotel booking failed", "error", err)
		return "Booking cancelled", err
	}
	saga.AddCompensation(activities.UndoBookHotel, input.HotelId)

	// Simulate delay
	logger.Info("Sleeping for 1 second...")
	workflow.Sleep(ctx, time.Second)

	// Book Car
	var car string
	err = workflow.ExecuteActivity(ctx, activities.BookCar, input.CarId).Get(ctx, &car)
	if err != nil {
		logger.Warn("Car booking failed", "error", err)
		return "Booking cancelled", err
	}
	saga.AddCompensation(activities.UndoBookCar, input.CarId)

	// Send Notification
	var notification string
	err = workflow.ExecuteActivity(ctx, activities.NotifyUser, input.UserId).Get(ctx, &notification)

	output = fmt.Sprintf("%s %s %s", flight, hotel, car)
	return output, nil
}
