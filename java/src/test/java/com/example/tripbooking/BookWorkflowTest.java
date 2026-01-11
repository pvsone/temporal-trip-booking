package com.example.tripbooking;

import com.example.tripbooking.activities.TripActivities;
import com.example.tripbooking.activities.TripActivitiesImpl;
import com.example.tripbooking.shared.BookTripInput;
import com.example.tripbooking.workflows.BookWorkflow;
import com.example.tripbooking.workflows.BookWorkflowImpl;
import io.temporal.testing.TestWorkflowEnvironment;
import io.temporal.testing.TestWorkflowExtension;
import io.temporal.worker.Worker;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.RegisterExtension;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.Mockito.*;

public class BookWorkflowTest {

    @RegisterExtension
    public static final TestWorkflowExtension testWorkflowExtension =
            TestWorkflowExtension.newBuilder()
                    .registerWorkflowImplementationTypes(BookWorkflowImpl.class)
                    .setDoNotStart(true)
                    .build();

    @Test
    public void testActivitiesImpl(
            TestWorkflowEnvironment testEnv, Worker worker, BookWorkflow workflow) {
        worker.registerActivitiesImplementations(new TripActivitiesImpl());
        testEnv.start();

        String result = workflow.bookTrip(new BookTripInput("Test User", "Flight 123", "Hotel 456", "Car 789"));
        assertEquals("Booked flight: Flight 123 Booked hotel: Hotel 456 Booked car: Car 789", result);
    }

    @Test
    public void testMockedActivities(
            TestWorkflowEnvironment testEnv, Worker worker, BookWorkflow workflow) {
        TripActivities activities =
                mock(TripActivities.class, withSettings().withoutAnnotations());
        when(activities.bookFlight("Flight 123")).thenReturn("Booked flight: Flight 123");
        when(activities.bookHotel("Hotel 456")).thenReturn("Booked hotel: Hotel 456");
        when(activities.bookCar("Car 789")).thenReturn("Booked car: Car 789");
        when(activities.notifyUser("Test User")).thenReturn("Notified user: Test User");
        worker.registerActivitiesImplementations(activities);
        testEnv.start();

        String result = workflow.bookTrip(new BookTripInput("Test User", "Flight 123", "Hotel 456", "Car 789"));
        assertEquals("Booked flight: Flight 123 Booked hotel: Hotel 456 Booked car: Car 789", result);
    }
}
