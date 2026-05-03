using System.Diagnostics;

using Greeter.V1;

using Microsoft.Extensions.Options;

using Platform.Telemetry;

namespace Greet.Worker;

public class Worker : BackgroundService {
    private readonly ILogger<Worker> _logger;
    private readonly GreeterService.GreeterServiceClient _greeterClient;
    private readonly ActivitySource _activitySource;
    private readonly AppSettings _settings;


    public Worker(ILogger<Worker> logger, IOptions<AppSettings> options, GreeterService.GreeterServiceClient greeterClient, IActivitySourceFactory factory) {
        _logger = logger;
        _greeterClient = greeterClient;
        _activitySource = factory.Create<Worker>();
        _settings = options.Value;
    }

    protected override async Task ExecuteAsync(CancellationToken stoppingToken) {
        _logger.LogInformation("Greeter URL: {Url}", _settings.Greeter.Url);

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
