package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

func BookFlight(ctx context.Context, flightId string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Booking flight", "flightId", flightId)

	time.Sleep(1 * time.Second)

	if containsIgnoreCase(flightId, "flaky") {
		if activity.GetInfo(ctx).Attempt < 6 {
			return "", errors.New("flight booking service is currently unavailable")
		}
	}

	return fmt.Sprintf("Booked flight: %s", flightId), nil
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

	return fmt.Sprintf("Booked hotel: %s", hotelId), nil
}

func BookCar(ctx context.Context, carId string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Booking car", "carId", carId)

	time.Sleep(1 * time.Second)

	if containsIgnoreCase(carId, "invalid") {
		return "", temporal.NewNonRetryableApplicationError(
			fmt.Sprintf("Car '%s' is invalid", carId),
			"InvalidCar",
			nil,
		)
	}

	return fmt.Sprintf("Booked car: %s", carId), nil
}

func NotifyUser(ctx context.Context, userId string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Notifying user", "userId", userId)

	time.Sleep(1 * time.Second)

	return fmt.Sprintf("Notified user: %s", userId), nil
}

func UndoBookFlight(ctx context.Context, flightId string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Undo flight booking", "flightId", flightId)

	time.Sleep(1 * time.Second)

	return fmt.Sprintf("Unbooked flight: %s", flightId), nil
}

func UndoBookHotel(ctx context.Context, hotelId string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Undo hotel booking", "hotelId", hotelId)

	time.Sleep(1 * time.Second)

	return fmt.Sprintf("Unbooked hotel: %s", hotelId), nil
}

func UndoBookCar(ctx context.Context, carId string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Undo car booking", "carId", carId)

	time.Sleep(1 * time.Second)

	return fmt.Sprintf("Unbooked car: %s", carId), nil
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
