using System.Diagnostics;

using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.DependencyInjection;

namespace Platform.Observability;

public static class CorrelationExtensions {
    public static IServiceCollection AddPlatformCorrelation(
        this IServiceCollection services) {

        services.AddHttpContextAccessor();

        services.AddSingleton<ICorrelationContext, CorrelationContext>();

        return services;
    }
}

public interface ICorrelationContext {
    string TraceId { get; }
    string SpanId { get; }
}


public class CorrelationContext : ICorrelationContext {
    public string TraceId =>
        Activity.Current?.TraceId.ToString()
        ?? System.Guid.NewGuid().ToString("N");

    public string SpanId =>
        Activity.Current?.SpanId.ToString()
        ?? string.Empty;
}
