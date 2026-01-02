using Microsoft.Extensions.Logging;
using Temporalio.Common;
using Temporalio.Workflows;
using TemporalTripBooking.Model;

namespace TemporalTripBooking;

[Workflow("BookWorkflow")]
public class BookWorkflow
{
    [WorkflowRun]
    public async Task<string> Execute(BookTripInput input)
    {
        Workflow.Logger.LogInformation("Book workflow started, user_id = {}", input.UserId);

        var compensations = new Stack<Func<Task>>();

        try
        {
            compensations.Push(async () => 
                await Workflow.ExecuteActivityAsync((TripActivities act) =>
                    act.UndoBookFlight(input), TripActivities.ActivityOpts));

            var flight = await Workflow.ExecuteActivityAsync((TripActivities act) =>
                act.BookFlight(input), TripActivities.ActivityOpts);

            Workflow.Logger.LogInformation("Sleeping for 1 second...");
            await Workflow.DelayAsync(TimeSpan.FromSeconds(1));

            compensations.Push(async () => 
                await Workflow.ExecuteActivityAsync((TripActivities act) =>
                    act.UndoBookHotel(input), TripActivities.ActivityOpts));

            var hotel = await Workflow.ExecuteActivityAsync((TripActivities act) =>
                act.BookHotel(input), TripActivities.ActivityOpts);

            Workflow.Logger.LogInformation("Sleeping for 1 second...");
            await Workflow.DelayAsync(TimeSpan.FromSeconds(1));

            compensations.Push(async () => 
                await Workflow.ExecuteActivityAsync((TripActivities act) =>
                    act.UndoBookCar(input), TripActivities.ActivityOpts));

            var car = await Workflow.ExecuteActivityAsync((TripActivities act) =>
                act.BookCar(input), TripActivities.ActivityOpts);

            await Workflow.ExecuteActivityAsync((TripActivities act) =>
                act.NotifyUser(input), TripActivities.ActivityOpts);

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
