using Handler;
using Microsoft.AspNetCore.Server.Kestrel.Core;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddGrpc();
builder.Services.AddGrpcReflection();

builder.WebHost.ConfigureKestrel(options => {
    options.ConfigureEndpointDefaults(o => {
        o.Protocols = Microsoft.AspNetCore.Server.Kestrel.Core.HttpProtocols.Http2;
    });
});

var app = builder.Build();
app.MapGrpcService<GreeterServiceHandler>();

if (app.Environment.IsDevelopment()) {
    app.MapGrpcReflectionService();
}

app.Run();