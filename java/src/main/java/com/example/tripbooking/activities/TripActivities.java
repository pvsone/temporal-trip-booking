package com.example.tripbooking.activities;

import io.temporal.activity.ActivityInterface;
import io.temporal.activity.ActivityMethod;

@ActivityInterface
public interface TripActivities {
    @ActivityMethod
    String bookFlight(String flightId);

    @ActivityMethod
    String bookHotel(String hotelId);

    @ActivityMethod
    String bookCar(String carId);

    @ActivityMethod
    String notifyUser(String userId);

    @ActivityMethod
    String undoBookFlight(String flightId);

    @ActivityMethod
    String undoBookHotel(String hotelId);

    @ActivityMethod
    String undoBookCar(String carId);
}
