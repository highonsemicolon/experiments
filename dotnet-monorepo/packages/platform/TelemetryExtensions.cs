using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;

namespace Platform.Telemetry;

public static class TelemetryExtensions
{
    public static IServiceCollection AddPlatformTelemetry(
        this IServiceCollection services,
        IHostEnvironment env)
    {
        services.AddOpenTelemetry()
            .ConfigureResource(resource =>
            {
                resource.AddService(
                    serviceName: env.ApplicationName,
                    serviceVersion: "1.0.0",
                    serviceInstanceId: Environment.MachineName);
            })
            .WithTracing(tracing =>
            {
                tracing
                    .AddAspNetCoreInstrumentation()
                    .AddHttpClientInstrumentation()
                    .SetSampler(new AlwaysOnSampler())
                    .AddOtlpExporter(); // later overridden in GKE
            });

        return services;
    }
}