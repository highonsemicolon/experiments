using System.Diagnostics;

using Greeter.V1;

using Platform.Telemetry;

namespace Greet.Worker;

public class Worker : BackgroundService {
    private readonly ILogger<Worker> _logger;
    private readonly GreeterService.GreeterServiceClient _greeterClient;
    private readonly ActivitySource _activitySource;


    public Worker(ILogger<Worker> logger, GreeterService.GreeterServiceClient greeterClient, IActivitySourceFactory factory) {
        _logger = logger;
        _greeterClient = greeterClient;
        _activitySource = factory.Create<Worker>();
    }

    protected override async Task ExecuteAsync(CancellationToken stoppingToken) {
        while (!stoppingToken.IsCancellationRequested) {
            using var activity = _activitySource.StartActivity("message-loop", ActivityKind.Consumer);
            await SendGreeting(stoppingToken);
            await Task.Delay(1000, stoppingToken);
        }
    }

    private async Task SendGreeting(CancellationToken stoppingToken) {
        var reply = await _greeterClient.SayHelloAsync(new SayHelloRequest { Name = "World" }, cancellationToken: stoppingToken);
        _logger.LogInformation("Greeting: {Greeting}", reply.Message);
    }
}
