using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Platform.Logging;
using Platform.Telemetry;

namespace Platform.Hosting;

public static class HostingExtensions {
    public static WebApplicationBuilder AddPlatformHost(
        this WebApplicationBuilder builder) {
        builder.AddPlatformLogging();
        builder.Services.AddPlatformTelemetry(builder.Environment);

        return builder;
    }
}
