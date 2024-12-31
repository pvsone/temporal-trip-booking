package activities

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

func BookCar(ctx context.Context, carId string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Booking car", "carId", carId)

	time.Sleep(1 * time.Second)

	if containsIgnoreCase(carId, "flaky") {
		if activity.GetInfo(ctx).Attempt < 6 {
			return "", errors.New("car booking service is currently unavailable")
		}
	}

	result := fmt.Sprintf("Booked car: %s", carId)
	return result, nil
}

func BookHotel(ctx context.Context, hotelId string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Booking hotel", "hotelId", hotelId)

	time.Sleep(1 * time.Second)

	if containsIgnoreCase(hotelId, "buggy") {
		// a logical bug in the code
		error := true
		if error {
			return "", errors.New("error due to bug in code")
		}
	}

	result := fmt.Sprintf("Booked hotel: %s", hotelId)
	return result, nil
}

func BookFlight(ctx context.Context, flightId string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Booking flight", "flightId", flightId)

	time.Sleep(1 * time.Second)

	if containsIgnoreCase(flightId, "invalid") {
		return "", temporal.NewNonRetryableApplicationError(
			fmt.Sprintf("Flight '%s' is not valid", flightId),
			"InvalidFlight", nil)
	}

	result := fmt.Sprintf("Booked flight: %s", flightId)
	return result, nil
}

func UndoBookCar(ctx context.Context, carId string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Undo car booking", "carId", carId)

	time.Sleep(1 * time.Second)

	return nil
}

func UndoBookHotel(ctx context.Context, hotelId string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Undo hotel booking", "hotelId", hotelId)

	time.Sleep(1 * time.Second)

	return nil
}

func UndoBookFlight(ctx context.Context, flightId string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Undo flight booking", "flightId", flightId)

	time.Sleep(1 * time.Second)

	return nil
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
