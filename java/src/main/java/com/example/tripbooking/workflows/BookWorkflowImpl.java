package com.example.tripbooking.workflows;

import com.example.tripbooking.activities.TripActivities;
import com.example.tripbooking.shared.BookTripInput;
import io.temporal.activity.ActivityOptions;
import io.temporal.common.RetryOptions;
import io.temporal.failure.ActivityFailure;
import io.temporal.workflow.Saga;
import io.temporal.workflow.Workflow;
import java.time.Duration;
import org.slf4j.Logger;

public class BookWorkflowImpl implements BookWorkflow {

    private static final Logger log = Workflow.getLogger(BookWorkflowImpl.class);

    private final TripActivities activities = Workflow.newActivityStub(
        TripActivities.class,
        ActivityOptions.newBuilder()
            .setStartToCloseTimeout(Duration.ofSeconds(5))
            .setRetryOptions(
                RetryOptions.newBuilder()
                    .setInitialInterval(Duration.ofSeconds(1))
                    .setBackoffCoefficient(2.0)
                    .setMaximumInterval(Duration.ofSeconds(30))
                    .build()
            )
            .build()
    );

    private static void sleep(long ms) {
        log.info("Simulate delay for {} milliseconds...", ms);
        Workflow.sleep(ms);
    }

    @Override
    public String bookTrip(BookTripInput input) {
        log.info("Book workflow started for userId: {}", input.getUserId());

        // Create saga to manage compensations
        Saga saga = new Saga(new Saga.Options.Builder().setParallelCompensation(false).build());

        try {
            // Book Flight
            String flight = activities.bookFlight(input.getFlightId());
            saga.addCompensation(activities::undoBookFlight, input.getFlightId());
            sleep(1000);

            // Book Hotel
            String hotel = activities.bookHotel(input.getHotelId());
            saga.addCompensation(activities::undoBookHotel, input.getHotelId());
            sleep(1000);

            // Book Car
            String car = activities.bookCar(input.getCarId());
            saga.addCompensation(activities::undoBookCar, input.getCarId());

            // Send Notification
            activities.notifyUser(input.getUserId());

            return String.format("%s %s %s", flight, hotel, car);
        } catch (ActivityFailure af) {
            log.warn("Booking failed, compensating...", af);
            saga.compensate();
            return "Booking canceled";
        }
    }
}
