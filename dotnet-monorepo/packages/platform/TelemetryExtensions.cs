using System.Reflection;
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
                var version = Assembly.GetEntryAssembly()?.GetCustomAttribute<AssemblyInformationalVersionAttribute>()?.InformationalVersion ?? "unknown";
                resource.AddService(
                    serviceName: env.ApplicationName,
                    serviceVersion: version,
                    serviceInstanceId: Environment.MachineName);
            })
            .WithTracing(tracing =>
            {
                tracing
                    .AddAspNetCoreInstrumentation()
                    .AddHttpClientInstrumentation()
                    .AddOtlpExporter(); // configured via environment variables (e.g., in GKE)
            });

        return services;
    }
}