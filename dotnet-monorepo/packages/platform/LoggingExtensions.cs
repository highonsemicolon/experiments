using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.Hosting;
using Serilog;
using Serilog.Events;
using Serilog.Formatting.Compact;

namespace Platform.Logging;

public static class LoggingExtensions
{
    public static WebApplicationBuilder AddPlatformLogging(
        this WebApplicationBuilder builder)
    {
        var env = builder.Environment;

        var logger = new LoggerConfiguration()
            // Base level
            .MinimumLevel.Information()

            // Framework noise control
            .MinimumLevel.Override("Microsoft", LogEventLevel.Warning)
            .MinimumLevel.Override("System", LogEventLevel.Warning)

            // ASP.NET request pipeline noise
            .MinimumLevel.Override(
                "Microsoft.AspNetCore.Hosting.Diagnostics",
                LogEventLevel.Warning)

            .MinimumLevel.Override(
                "Microsoft.AspNetCore.Routing.EndpointMiddleware",
                LogEventLevel.Warning)

            // gRPC reflection / internal probe noise
            .MinimumLevel.Override(
                "Grpc.AspNetCore.Server.Internal.ServerCallHandlerFactory",
                LogEventLevel.Warning)

            // Structured enrichment
            .Enrich.FromLogContext()
            .Enrich.WithProperty("service", env.ApplicationName)
            .Enrich.WithProperty("environment", env.EnvironmentName);

        if (env.IsDevelopment())
        {
            // Human-friendly local logs
            logger = logger.WriteTo.Console(
                outputTemplate:
                "[{Timestamp:HH:mm:ss} {Level:u3}] {Message:lj} " +
                "(trace:{TraceId} span:{SpanId}){NewLine}{Exception}");
        }
        else
        {
            // GKE / Cloud Logging friendly JSON
            logger = logger.WriteTo.Console(new RenderedCompactJsonFormatter());
        }

        Log.Logger = logger.CreateLogger();

        builder.Host.UseSerilog();

        return builder;
    }
}