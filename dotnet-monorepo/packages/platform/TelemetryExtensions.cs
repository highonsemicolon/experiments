using System.Diagnostics;
using System.Reflection;

using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

using OpenTelemetry.Resources;
using OpenTelemetry.Trace;

namespace Platform.Telemetry;

public static class TelemetryExtensions {
    public static IServiceCollection AddPlatformTelemetry(
        this IServiceCollection services,
        IHostEnvironment env, params string[] activitySources) {

        services.AddSingleton<IActivitySourceFactory, ActivitySourceFactory>();

        services.AddOpenTelemetry()
            .ConfigureResource(resource => {
                var version = Assembly.GetEntryAssembly()?.GetCustomAttribute<AssemblyInformationalVersionAttribute>()?.InformationalVersion ?? "unknown";
                resource.AddService(
                    serviceName: env.ApplicationName,
                    serviceVersion: version,
                    serviceInstanceId: Environment.MachineName);
            })
            .WithTracing(tracing => {
                tracing
                    .AddSource(env.ApplicationName)
                    .AddAspNetCoreInstrumentation()
                    .AddHttpClientInstrumentation(opts => {
                        opts.FilterHttpRequestMessage = _ => true;
                    })
                    .AddOtlpExporter(); // configured via environment variables (e.g., in GKE)
                foreach (var source in activitySources.Distinct())
                    tracing.AddSource(source);
            });

        return services;
    }
}



public interface IActivitySourceFactory {
    ActivitySource Create<T>();
}

public class ActivitySourceFactory : IActivitySourceFactory {
    private readonly ActivitySource _source;

    public ActivitySourceFactory(IHostEnvironment env) {
        _source = new ActivitySource(env.ApplicationName);
    }

    public ActivitySource Create<T>() => _source;
}
