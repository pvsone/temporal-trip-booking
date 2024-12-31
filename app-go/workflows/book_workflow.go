package workflows

import (
	"fmt"
	"temporal-sagas/activities"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/multierr"
)

type BookTripInput struct {
	BookUserId   string `json:"book_user_id"`
	BookCarId    string `json:"book_car_id"`
	BookHotelId  string `json:"book_hotel_id"`
	BookFlightId string `json:"book_flight_id"`
}

func BookWorkflow(ctx workflow.Context, input BookTripInput) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Book workflow started")

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    1 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    30 * time.Second,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Book car
	var car string
	err := workflow.ExecuteActivity(ctx, activities.BookCar, input.BookCarId).Get(ctx, &car)
	if err != nil {
		logger.Warn("Car booking failed", "error", err)
		return "Booking cancelled", nil
	}

	// Add compensation function to undo car booking
	defer func() {
		if err != nil {
			errCompensation := workflow.ExecuteActivity(ctx, activities.UndoBookCar, input.BookCarId).Get(ctx, nil)
			err = multierr.Append(err, errCompensation)
		}
	}()

	// Dramatic pause
	logger.Info("Sleeping for 1 second...")
	workflow.Sleep(ctx, 1*time.Second)

	// Book hotel
	var hotel string
	err = workflow.ExecuteActivity(ctx, activities.BookHotel, input.BookHotelId).Get(ctx, &hotel)
	if err != nil {
		logger.Warn("Hotel booking failed", "error", err)
		return "Booking cancelled", nil
	}

	// Add compensation function to undo hotel booking
	defer func() {
		if err != nil {
			errCompensation := workflow.ExecuteActivity(ctx, activities.UndoBookHotel, input.BookHotelId).Get(ctx, nil)
			err = multierr.Append(err, errCompensation)
		}
	}()

	// Dramatic pause
	logger.Info("Sleeping for 1 second...")
	workflow.Sleep(ctx, 1*time.Second)

	// Book flight
	var flight string
	err = workflow.ExecuteActivity(ctx, activities.BookFlight, input.BookFlightId).Get(ctx, &flight)
	if err != nil {
		logger.Warn("Flight booking failed", "error", err)
		return "Booking cancelled", nil
	}

	output := fmt.Sprintf("%s %s %s", car, hotel, flight)
	return output, nil
}
