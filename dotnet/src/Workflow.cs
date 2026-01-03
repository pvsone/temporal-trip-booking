using Microsoft.Extensions.Logging;
using Temporalio.Common;
using Temporalio.Workflows;

namespace TemporalTripBooking;

[Workflow("BookWorkflow")]
public class BookWorkflow
{
    [WorkflowRun]
    public async Task<string> Execute(Shared input)
    {
        Workflow.Logger.LogInformation("Book workflow started, user_id = {}", input.UserId);
        var options = new ActivityOptions() { 
            StartToCloseTimeout = TimeSpan.FromSeconds(5),
            RetryPolicy = new() {
                InitialInterval = TimeSpan.FromSeconds(1),
                MaximumInterval = TimeSpan.FromSeconds(30),
                BackoffCoefficient = 2
            }
        };
        
        var compensations = new Stack<Func<Task>>();

        try
        {
            // Book Flight
            var flight = await Workflow.ExecuteActivityAsync(() =>
                TripActivities.BookFlight(input), options);
            compensations.Push(async () => 
                await Workflow.ExecuteActivityAsync(() => TripActivities.UndoBookFlight(input), options));
            
            Workflow.Logger.LogInformation("Sleeping for 1 second...");
            await Workflow.DelayAsync(TimeSpan.FromSeconds(1));

            // Book Hotel
            var hotel = await Workflow.ExecuteActivityAsync(() =>
                TripActivities.BookHotel(input), options);
            compensations.Push(async () => 
                await Workflow.ExecuteActivityAsync(() => TripActivities.UndoBookHotel(input), options));
            
            Workflow.Logger.LogInformation("Sleeping for 1 second...");
            await Workflow.DelayAsync(TimeSpan.FromSeconds(1));

            // Book Car
            var car = await Workflow.ExecuteActivityAsync(() =>
                TripActivities.BookCar(input), options);
            compensations.Push(async () => 
                await Workflow.ExecuteActivityAsync(() => TripActivities.UndoBookCar(input), options));

            // Notify User
            await Workflow.ExecuteActivityAsync(() =>
                TripActivities.NotifyUser(input), options);

            return $"{flight} {hotel} {car}";
        }
        catch (Exception ex)
        {
            Workflow.Logger.LogError(ex, "Failed to book trip");
            while (compensations.TryPop(out var compensation))
            {
                try
                {
                    await compensation();
                }
                catch (Exception compEx)
                {
                    Workflow.Logger.LogWarning(compEx, "Compensation failed");
                }
            }

            return "Booking cancelled";
        }
    }
}
