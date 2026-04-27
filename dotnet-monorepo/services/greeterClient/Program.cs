using System;
using System.Threading.Tasks;
using Grpc.Net.Client;
using Greeter.V1;
using System.Diagnostics;

using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

using Platform.Observability;
using Serilog;

using Platform.Telemetry;

namespace Greeter.Client;

class Program
{
    static async Task Main(string[] args)
    {
        var builder = Host.CreateApplicationBuilder(args);

        builder.AddPlatformObservability();
        // builder.Services.AddSingleton(new ActivitySource(builder.Environment.ApplicationName));

        var host = builder.Build();
        await host.StartAsync();

        var activityFactory = host.Services.GetRequiredService<IActivitySourceFactory>();
        var activitySource = activityFactory.Create<Program>();
        using var activity = activitySource.StartActivity("SayHello", ActivityKind.Client);


        using var channel = GrpcChannel.ForAddress("http://localhost:8080");

        var client = new GreeterService.GreeterServiceClient(channel);

        var request = new SayHelloRequest
        {
            Name = "Alice"
        };

        
        var response = await client.SayHelloAsync(request);

        Log.Information("Server response: {Message}", response.Message);

        await host.StopAsync(TimeSpan.FromSeconds(5));
        Log.CloseAndFlush();
    }
}