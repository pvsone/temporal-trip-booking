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
	BookUserId   string `json:"book_user_id"`
	BookCarId    string `json:"book_car_id"`
	BookHotelId  string `json:"book_hotel_id"`
	BookFlightId string `json:"book_flight_id"`
}

func BookWorkflow(ctx workflow.Context, input BookTripInput) (output string, err error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Book workflow started", "userId", input.BookUserId)

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
	err = workflow.ExecuteActivity(ctx, activities.BookFlight, input.BookFlightId).Get(ctx, &flight)
	if err != nil {
		logger.Warn("Flight booking failed", "error", err)
		return "Booking cancelled", err
	}
	saga.AddCompensation(activities.UndoBookFlight, input.BookFlightId)

	// Simulate delay
	logger.Info("Sleeping for 1 second...")
	workflow.Sleep(ctx, time.Second)

	// Book Hotel
	var hotel string
	err = workflow.ExecuteActivity(ctx, activities.BookHotel, input.BookHotelId).Get(ctx, &hotel)
	if err != nil {
		logger.Warn("Hotel booking failed", "error", err)
		return "Booking cancelled", err
	}
	saga.AddCompensation(activities.UndoBookHotel, input.BookHotelId)

	// Simulate delay
	logger.Info("Sleeping for 1 second...")
	workflow.Sleep(ctx, time.Second)

	// Book Car
	var car string
	err = workflow.ExecuteActivity(ctx, activities.BookCar, input.BookCarId).Get(ctx, &car)
	if err != nil {
		logger.Warn("Car booking failed", "error", err)
		return "Booking cancelled", err
	}
	saga.AddCompensation(activities.UndoBookCar, input.BookCarId)

	// Send Notification
	var notification string
	err = workflow.ExecuteActivity(ctx, activities.NotifyUser, input.BookUserId).Get(ctx, &notification)

	output = fmt.Sprintf("%s %s %s", flight, hotel, car)
	return output, nil
}
