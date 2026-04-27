using Greet.Worker;

using Greeter.V1;

using Platform.Observability;

var builder = Host.CreateApplicationBuilder(args);

builder.AddPlatformObservability();

builder.Services.AddGrpcClient<GreeterService.GreeterServiceClient>(options => {
    options.Address = new Uri("http://localhost:8080");
});

builder.Services.AddHostedService<Worker>();

var host = builder.Build();
host.Run();
