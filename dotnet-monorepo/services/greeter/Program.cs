using Handler;

using Microsoft.AspNetCore.Server.Kestrel.Core;

using Platform.Observability;

var builder = WebApplication.CreateBuilder(args);
builder.AddPlatformObservability();

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
