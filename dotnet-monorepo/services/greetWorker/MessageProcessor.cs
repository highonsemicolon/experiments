using System.Diagnostics;

using Greeter.V1;

using Platform.Telemetry;

namespace Greet.Worker;

public sealed class MessageProcessor : IMessageProcessor {
    private readonly ILogger<MessageProcessor> _logger;
    private readonly GreeterService.GreeterServiceClient _greeterClient;
    private readonly ActivitySource _activitySource;

    public MessageProcessor(
        ILogger<MessageProcessor> logger,
        GreeterService.GreeterServiceClient greeterClient,
        IActivitySourceFactory factory) {
        _logger = logger;
        _greeterClient = greeterClient;
        _activitySource = factory.Create<MessageProcessor>();
    }

    public async Task ProcessAsync(string message, CancellationToken cancellationToken) {
        using var activity = _activitySource.StartActivity("greet-processing", ActivityKind.Internal);

        _logger.LogInformation("Processing message {Message}", message);

        var name = message.Trim().ToUpperInvariant();

        _logger.LogInformation("Normalized name {Name}", name);

        var reply = await _greeterClient.SayHelloAsync(new SayHelloRequest { Name = name }, cancellationToken: cancellationToken);

        _logger.LogInformation("Greeter response {Message}", reply.Message);
    }
}
