using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Hosting;

using Serilog;
using Serilog.Enrichers.Span;
using Serilog.Events;

namespace Platform.Logging;

public static class LoggingExtensions {
    public static WebApplicationBuilder AddPlatformLogging(
        this WebApplicationBuilder builder) {
        var logger = BuildLogger(
            builder.Environment,
            builder.Configuration);

        Log.Logger = logger;

        builder.Host.UseSerilog();

        return builder;
    }

    public static HostApplicationBuilder AddPlatformLogging(
        this HostApplicationBuilder builder) {
        var logger = BuildLogger(
            builder.Environment,
            builder.Configuration);

        Log.Logger = logger;

        builder.Services.AddSerilog(logger);

        return builder;
    }

    private static ILogger BuildLogger(
        IHostEnvironment env,
        IConfiguration config) {
        var projectId = config["GCP_PROJECT_ID"] ?? "local";

        var logger = new LoggerConfiguration()
            .MinimumLevel.Information()
            .MinimumLevel.Override("Microsoft", LogEventLevel.Warning)
            .MinimumLevel.Override("System", LogEventLevel.Warning)
            .Enrich.FromLogContext()
            .Enrich.WithSpan()
            .Enrich.WithProperty("service", env.ApplicationName)
            .Enrich.WithProperty("environment", env.EnvironmentName);

        if (env.IsDevelopment()) {
            logger = logger.WriteTo.Console(
                outputTemplate:
                "[{Timestamp:HH:mm:ss} {Level:u3}] {Message:lj} " +
                "(trace:{TraceId} span:{SpanId}){NewLine}{Exception}");
        }
        else {
            logger = logger.WriteTo.Console(
                new GcpJsonFormatter(projectId));
        }

        return logger.CreateLogger();
    }
}
