using System.Diagnostics;
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
        var activity = Activity.Current;

        var log = new Dictionary<string, object?> {
            ["severity"] = MapLevel(logEvent.Level),
            ["message"] = logEvent.RenderMessage(),
            ["timestamp"] = logEvent.Timestamp.UtcDateTime.ToString("o")
        };

        foreach (var p in logEvent.Properties) {
            log[p.Key] = p.Value.ToString().Trim('"');
        }

        if (activity != null) {
            var traceId = activity.TraceId.ToString();
            var spanId = activity.SpanId.ToString();

            log["traceId"] = traceId;
            log["spanId"] = spanId;

            log["logging.googleapis.com/trace"] =
                $"projects/{_projectId}/traces/{traceId}";

            log["logging.googleapis.com/spanId"] = spanId;
        }

        output.WriteLine(JsonSerializer.Serialize(log));
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
