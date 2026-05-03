using System.Diagnostics;

using Microsoft.Extensions.Options;

using Platform.Telemetry;

namespace Greet.Worker;

public class Worker : BackgroundService {
    private readonly ILogger<Worker> _logger;
    private readonly ActivitySource _activitySource;
    private readonly AppSettings _settings;
    private readonly IServiceScopeFactory _scopeFactory;

    public Worker(ILogger<Worker> logger, IOptions<AppSettings> options, IActivitySourceFactory factory, IServiceScopeFactory scopeFactory) {
        _logger = logger;
        _activitySource = factory.Create<Worker>();
        _settings = options.Value;
        _scopeFactory = scopeFactory;
    }

    protected override async Task ExecuteAsync(CancellationToken stoppingToken) {
        while (!stoppingToken.IsCancellationRequested) {
            using var activity = _activitySource.StartActivity("message-loop", ActivityKind.Consumer);

            using var scope = _scopeFactory.CreateScope();
            var processor = scope.ServiceProvider.GetRequiredService<IMessageProcessor>();
            await processor.ProcessAsync("World", stoppingToken);

            await Task.Delay(1000, stoppingToken);
        }
    }
}
