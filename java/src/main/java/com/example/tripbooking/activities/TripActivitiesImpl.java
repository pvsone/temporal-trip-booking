package com.example.tripbooking.activities;

import io.temporal.activity.Activity;
import io.temporal.failure.ApplicationFailure;
import lombok.extern.slf4j.Slf4j;

import java.util.concurrent.TimeUnit;

@Slf4j
public class TripActivitiesImpl implements TripActivities {

    private static void sleep(long ms) {
        try {
            TimeUnit.MILLISECONDS.sleep(ms);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }

    @Override
    public String bookFlight(String flightId) {
        log.info("Booking flight: {}", flightId);

        sleep(1000);

        if (containsIgnoreCase(flightId, "flaky")) {
            // a transient error, which will be retried
            if (Activity.getExecutionContext().getInfo().getAttempt() < 6) {
                throw new RuntimeException("flight booking service is currently unavailable");
            }
        }

        return String.format("Booked flight: %s", flightId);
    }

    @Override
    public String bookHotel(String hotelId) {
        log.info("Booking hotel: {}", hotelId);

        sleep(1000);

        if (containsIgnoreCase(hotelId, "buggy")) {
            // a logical bug in the code
            boolean error = true;
            if (error) {
                throw new RuntimeException("error due to bug in code");
            }
        }

        return String.format("Booked hotel: %s", hotelId);
    }

    @Override
    public String bookCar(String carId) {
        log.info("Booking car: {}", carId);

        sleep(1000);

        if (containsIgnoreCase(carId, "invalid")) {
            // a business error, which cannot be retried
            throw ApplicationFailure.newNonRetryableFailure(String.format("Car '%s' is invalid", carId), "InvalidCar");
        }

        return String.format("Booked car: %s", carId);
    }

    @Override
    public String notifyUser(String userId) {
        log.info("Notifying user: {}", userId);

        sleep(1000);

        return String.format("Notified user: %s", userId);
    }

    @Override
    public String undoBookFlight(String flightId) {
        log.info("Undo flight booking: {}", flightId);

        sleep(1000);

        return String.format("Unbooked flight: %s", flightId);
    }

    @Override
    public String undoBookHotel(String hotelId) {
        log.info("Undo hotel booking: {}", hotelId);

        sleep(1000);

        return String.format("Unbooked hotel: %s", hotelId);
    }

    @Override
    public String undoBookCar(String carId) {
        log.info("Undo car booking: {}", carId);

        sleep(1000);

        return String.format("Unbooked car: %s", carId);
    }

    private boolean containsIgnoreCase(String s, String substr) {
        return s.toLowerCase().contains(substr.toLowerCase());
    }
}
