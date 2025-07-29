using Google.Cloud.Diagnostics.Common;

namespace dotnet_gcp_logging;

public class Worker : BackgroundService
{
    private readonly ILogger<Worker> _logger;
    private readonly IManagedTracer _tracer;


    public Worker(ILogger<Worker> logger, IManagedTracer tracer)
    {
        _logger = logger;
        _tracer = tracer;
    }

    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        while (!stoppingToken.IsCancellationRequested)
        {
            await _tracer.RunInSpan(async() =>
            {
                _logger.LogInformation("Worker running at: {time}", DateTimeOffset.Now);
                await Task.Delay(1000, stoppingToken);
            }, "myWorker");
        }
    }
}
