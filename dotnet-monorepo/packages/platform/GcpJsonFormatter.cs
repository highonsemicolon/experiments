using System.Text.Json;

using Serilog.Events;
using Serilog.Formatting;

namespace Platform.Logging;

public class GcpJsonFormatter : ITextFormatter {
    private readonly string _projectId;

    public GcpJsonFormatter(string projectId) {
        _projectId = projectId;
    }

    public void Format(LogEvent logEvent, TextWriter output) {
        var log = new Dictionary<string, object?> {
            ["severity"] = MapLevel(logEvent.Level),
            ["message"] = logEvent.RenderMessage(),
            ["timestamp"] = logEvent.Timestamp.UtcDateTime.ToString("o")
        };

        // Copy all structured properties
        foreach (var p in logEvent.Properties) {
            log[p.Key] = Simplify(p.Value);
        }

        // Extract TraceId / SpanId from enriched properties
        var traceId = GetScalar(logEvent, "TraceId");
        var spanId = GetScalar(logEvent, "SpanId");

        if (!string.IsNullOrEmpty(traceId)) {
            log["traceId"] = traceId;
            log["logging.googleapis.com/trace"] =
                $"projects/{_projectId}/traces/{traceId}";
        }

        if (!string.IsNullOrEmpty(spanId)) {
            log["spanId"] = spanId;
            log["logging.googleapis.com/spanId"] = spanId;
        }

        output.WriteLine(JsonSerializer.Serialize(log));
    }

    private static string? GetScalar(LogEvent logEvent, string key) {
        if (logEvent.Properties.TryGetValue(key, out var value)) {
            if (value is ScalarValue scalar && scalar.Value != null)
                return scalar.Value.ToString();
        }
        return null;
    }

    private static object? Simplify(LogEventPropertyValue value) {
        return value switch {
            ScalarValue s => s.Value,
            SequenceValue seq => seq.Elements.Select(Simplify).ToArray(),
            StructureValue str => str.Properties.ToDictionary(p => p.Name, p => Simplify(p.Value)),
            DictionaryValue dict => dict.Elements.ToDictionary(
                k => k.Key.Value?.ToString() ?? "",
                v => Simplify(v.Value)
            ),
            _ => value.ToString()
        };
    }

    private static string MapLevel(LogEventLevel level) => level switch {
        LogEventLevel.Verbose => "DEBUG",
        LogEventLevel.Debug => "DEBUG",
        LogEventLevel.Information => "INFO",
        LogEventLevel.Warning => "WARNING",
        LogEventLevel.Error => "ERROR",
        LogEventLevel.Fatal => "CRITICAL",
        _ => "DEFAULT"
    };
}
