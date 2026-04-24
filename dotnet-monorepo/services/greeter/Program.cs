using Handler;
using Platform.Logging;
using Platform.Telemetry;
using Microsoft.AspNetCore.Server.Kestrel.Core;

var builder = WebApplication.CreateBuilder(args);
builder.AddPlatformLogging();
builder.AddPlatformTelemetry();

builder.Services.AddGrpc();
builder.Services.AddGrpcReflection();

builder.WebHost.ConfigureKestrel(options => {
    options.ConfigureEndpointDefaults(o => {
        o.Protocols = HttpProtocols.Http2;
    });
});

var app = builder.Build();
app.MapGrpcService<GreeterServiceHandler>();

if (app.Environment.IsDevelopment()) {
    app.MapGrpcReflectionService();
}

app.Run();