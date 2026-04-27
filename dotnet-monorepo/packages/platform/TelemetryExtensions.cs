using System.Diagnostics;
using System.Reflection;

using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

using OpenTelemetry.Exporter;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;

namespace Platform.Telemetry;

public static class TelemetryExtensions {
    public static IServiceCollection AddPlatformTelemetry(
        this IServiceCollection services,
        IHostEnvironment env) {

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
                    .AddSource("Grpc.Net.Client")
                    .AddAspNetCoreInstrumentation()
                    .AddHttpClientInstrumentation(opts => {
                        opts.FilterHttpRequestMessage = _ => true;
                    })
                    .AddConsoleExporter()
                    .AddOtlpExporter(); // configured via environment variables (e.g., in GKE)
            });

        return services;
    }
}



public interface IActivitySourceFactory {
    ActivitySource Create<T>();
}

// public class ActivitySourceFactory : IActivitySourceFactory
// {
//     private readonly string _serviceName;

//     public ActivitySourceFactory(IHostEnvironment env)
//     {
//         _serviceName = env.ApplicationName;
//     }

//     public ActivitySource Create<T>()
//     {
//         var name = $"{_serviceName}.{typeof(T).Name}";
//         return new ActivitySource(name);
//     }
// }


public class ActivitySourceFactory : IActivitySourceFactory {
    private readonly ActivitySource _source;

    public ActivitySourceFactory(IHostEnvironment env) {
        _source = new ActivitySource(env.ApplicationName);
    }

    public ActivitySource Create<T>() => _source;
}
