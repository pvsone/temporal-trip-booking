using Temporalio.Activities;
using Temporalio.Client;
using Temporalio.Testing;
using Temporalio.Worker;
using TripBooking;

namespace Tests;

public class WorkflowTests
{
    [Fact]
    public async Task BookWorkflow_Succeeds()
    {
        // Use real activities
        var activities = new TripActivities();

        await using var env = await WorkflowEnvironment.StartTimeSkippingAsync();

        using var worker = new TemporalWorker(
            env.Client,
            new TemporalWorkerOptions($"task-queue-{Guid.NewGuid()}")
                .AddWorkflow<BookWorkflow>()
                .AddAllActivities(activities));

        var input = new BookTripInput("123", "FL123", "HT123", "CR123");

        // Run the worker only for the life of the code within
        await worker.ExecuteAsync(async () =>
        {
            var result = await env.Client.ExecuteWorkflowAsync(
                (BookWorkflow wf) => wf.Execute(input),
                new WorkflowOptions($"wf-{Guid.NewGuid()}", worker.Options.TaskQueue!));
            Assert.Equal("Booked flight: FL123 Booked hotel: HT123 Booked car: CR123", result);
        });
    }

    [Fact]
    public async Task BookWorkflow_MockActivities_Succeeds()
    {
        // Use mock activities
        [Activity("BookFlight")]
        string MockBookFlight(BookTripInput input) => $"Booked flight: {input.FlightId}";

        [Activity("BookHotel")]
        string MockBookHotel(BookTripInput input) => $"Booked hotel: {input.HotelId}";

        [Activity("BookCar")]
        string MockBookCar(BookTripInput input) => $"Booked car: {input.CarId}";

        [Activity("NotifyUser")]
        string MockNotifyUser(BookTripInput input) => $"Notified user: {input.UserId}";

        await using var env = await WorkflowEnvironment.StartTimeSkippingAsync();

        using var worker = new TemporalWorker(
            env.Client,
            new TemporalWorkerOptions($"task-queue-{Guid.NewGuid()}")
                .AddWorkflow<BookWorkflow>()
                .AddActivity(MockBookFlight)
                .AddActivity(MockBookHotel)
                .AddActivity(MockBookCar)
                .AddActivity(MockNotifyUser));

        var input = new BookTripInput("123", "FL123", "HT123", "CR123");

        // Run the worker only for the life of the code within
        await worker.ExecuteAsync(async () =>
        {
            var result = await env.Client.ExecuteWorkflowAsync(
                (BookWorkflow wf) => wf.Execute(input),
                new WorkflowOptions($"wf-{Guid.NewGuid()}", worker.Options.TaskQueue!));
            Assert.Equal("Booked flight: FL123 Booked hotel: HT123 Booked car: CR123", result);
        });
    }
}