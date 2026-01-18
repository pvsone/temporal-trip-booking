using Microsoft.Extensions.Logging;
using Temporalio.Client;
using Temporalio.Common.EnvConfig;
using Temporalio.Worker;
using TripBooking;

Console.WriteLine("⚙️ Using TEMPORAL_PROFILE: '{0}'", Environment.GetEnvironmentVariable("TEMPORAL_PROFILE"));
var connectOptions = ClientEnvConfig.LoadClientConnectOptions();
connectOptions.LoggerFactory = LoggerFactory.Create(builder =>
    builder.AddSimpleConsole(options => options.TimestampFormat = "[HH:mm:ss] ")
        .SetMinimumLevel(LogLevel.Information));
var client = await TemporalClient.ConnectAsync(connectOptions);
Console.WriteLine("✅ Client connected to '{0}' in namespace '{1}'", connectOptions.TargetHost,
    connectOptions.Namespace);


// Cancellation token to shutdown worker on ctrl+c
using var tokenSource = new CancellationTokenSource();
Console.CancelKeyPress += (_, eventArgs) =>
{
    tokenSource.Cancel();
    eventArgs.Cancel = true;
};

var activities = new TripActivities();

using var worker = new TemporalWorker(
    client,
    new TemporalWorkerOptions("trip-task-queue")
        .AddAllActivities(activities)
        .AddWorkflow<BookWorkflow>()
);

// Run worker until canceled
Console.WriteLine("Running worker...");
try
{
    await worker.ExecuteAsync(tokenSource.Token);
}
catch (OperationCanceledException)
{
    Console.WriteLine("Worker canceled");
}