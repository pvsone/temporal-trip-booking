package com.example.tripbooking.model;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class BookTripInput {

    private String userId;
    private String flightId;
    private String hotelId;
    private String carId;
}
