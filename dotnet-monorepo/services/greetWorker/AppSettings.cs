public sealed class AppSettings {
    public GreeterSettings Greeter { get; init; } = new();
    public ObservabilitySettings Observability { get; init; } = new();
}

public sealed class GreeterSettings {
    public Uri Url { get; init; } = default!;
}

public sealed class ObservabilitySettings {
    public string[] ActivitySources { get; init; } = default!;
}
