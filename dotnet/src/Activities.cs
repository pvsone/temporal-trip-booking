using Microsoft.Extensions.Logging;
using Temporalio.Activities;
using Temporalio.Exceptions;

namespace TemporalTripBooking;

public class TripActivities
{
    [Activity]
    public static async Task<string> BookFlight(BookTripInput input)
    {
        var logger = ActivityExecutionContext.Current.Logger;
        logger.LogInformation("Booking flight: {FlightId}", input.FlightId);

        await Task.Delay(TimeSpan.FromSeconds(1));

        if (input.FlightId.Contains("flaky", StringComparison.OrdinalIgnoreCase) &&
            ActivityExecutionContext.Current.Info.Attempt < 6)
            throw new Exception("Flight booking service is currently unavailable");

        return $"Booked flight: {input.FlightId}";
    }

    [Activity]
    public static async Task<string> BookHotel(BookTripInput input)
    {
        var logger = ActivityExecutionContext.Current.Logger;
        logger.LogInformation("Booking hotel: {HotelId}", input.HotelId);

        await Task.Delay(TimeSpan.FromSeconds(1));

        if (input.HotelId.Contains("buggy", StringComparison.OrdinalIgnoreCase))
        {
            var error = true;
            if (error) throw new Exception("Error due to bug in code");
        }

        return $"Booked hotel: {input.HotelId}";
    }

    [Activity]
    public static async Task<string> BookCar(BookTripInput input)
    {
        var logger = ActivityExecutionContext.Current.Logger;
        logger.LogInformation("Booking car: {CarId}", input.CarId);

        await Task.Delay(TimeSpan.FromSeconds(1));

        if (input.CarId.Contains("invalid", StringComparison.OrdinalIgnoreCase))
            throw new ApplicationFailureException(
                $"Car {input.CarId} is invalid", "InvalidCar", true);

        return $"Booked car: {input.CarId}";
    }

    [Activity]
    public static async Task<string> NotifyUser(BookTripInput input)
    {
        var logger = ActivityExecutionContext.Current.Logger;
        logger.LogInformation("Notifying user: {UserId}", input.UserId);

        await Task.Delay(TimeSpan.FromSeconds(1));

        return $"Notified user: {input.UserId}";
    }

    [Activity]
    public static async Task<string> UndoBookFlight(BookTripInput input)
    {
        var logger = ActivityExecutionContext.Current.Logger;
        logger.LogInformation("Undo flight booking: {FlightId}", input.FlightId);

        await Task.Delay(TimeSpan.FromSeconds(1));

        return $"Unbooked flight: {input.FlightId}";
    }

    [Activity]
    public static async Task<string> UndoBookHotel(BookTripInput input)
    {
        var logger = ActivityExecutionContext.Current.Logger;
        logger.LogInformation("Undo hotel booking: {HotelId}", input.HotelId);

        await Task.Delay(TimeSpan.FromSeconds(1));

        return $"Unbooked hotel: {input.HotelId}";
    }

    [Activity]
    public static async Task<string> UndoBookCar(BookTripInput input)
    {
        var logger = ActivityExecutionContext.Current.Logger;
        logger.LogInformation("Undo car booking: {CarId}", input.CarId);

        await Task.Delay(TimeSpan.FromSeconds(1));

        return $"Unbooked car: {input.CarId}";
    }
}