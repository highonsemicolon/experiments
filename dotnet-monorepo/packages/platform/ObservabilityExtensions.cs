using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.Hosting;

using Platform.Logging;
using Platform.Telemetry;

namespace Platform.Observability;

public static class ObservabilityExtensions {
    public static WebApplicationBuilder AddPlatformObservability(
        this WebApplicationBuilder builder, params string[] activitySources) {
        builder.AddPlatformLogging();

        builder.Services.AddPlatformTelemetry(builder.Environment, activitySources);
        builder.Services.AddPlatformCorrelation();

        return builder;
    }

    public static HostApplicationBuilder AddPlatformObservability(
        this HostApplicationBuilder builder, params string[] activitySources) {
        builder.AddPlatformLogging();

        builder.Services.AddPlatformTelemetry(builder.Environment, activitySources);
        builder.Services.AddPlatformCorrelation();

        return builder;
    }
}
