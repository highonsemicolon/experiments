using Microsoft.AspNetCore.Builder;
using OpenTelemetry.Extensions.Hosting;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using OpenTelemetry.Instrumentation.Http; 
using Microsoft.Extensions.DependencyInjection;

namespace Platform.Telemetry;

public static class TelemetryExtensions
{
    public static WebApplicationBuilder AddPlatformTelemetry(
        this WebApplicationBuilder builder)
    {
        var serviceName = builder.Environment.ApplicationName;

        builder.Services.AddOpenTelemetry()
            .WithTracing(tracing =>
            {
                tracing
                    .SetResourceBuilder(
                        ResourceBuilder.CreateDefault()
                            .AddService(serviceName))
                    .AddAspNetCoreInstrumentation()
                    .AddHttpClientInstrumentation()
                    .AddOtlpExporter();
            });

        return builder;
    }
}