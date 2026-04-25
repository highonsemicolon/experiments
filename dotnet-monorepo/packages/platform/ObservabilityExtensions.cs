using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

using Platform.Logging;
using Platform.Telemetry;

namespace Platform.Observability;

public static class ObservabilityExtensions {
    public static WebApplicationBuilder AddPlatformObservability(
        this WebApplicationBuilder builder) {
        builder.AddPlatformLogging();
        builder.Services.AddPlatformTelemetry(builder.Environment);
        builder.Services.AddPlatformCorrelation();

        return builder;
    }
}
