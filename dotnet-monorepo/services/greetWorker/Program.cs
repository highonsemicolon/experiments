using Greet.Worker;

using Greeter.V1;

using Microsoft.Extensions.Options;

using Platform.Observability;

var builder = Host.CreateApplicationBuilder(args);

builder.Services.AddOptions<AppSettings>().Bind(builder.Configuration).ValidateOnStart();
builder.Services.AddSingleton(sp => sp.GetRequiredService<IOptions<AppSettings>>().Value);

var settings = builder.Configuration.Get<AppSettings>()!;
builder.AddPlatformObservability(settings.Observability.ActivitySources);

builder.Services.AddGrpcClient<GreeterService.GreeterServiceClient>((options) => {
    options.Address = settings.Greeter.Url;
});

builder.Services.AddScoped<IMessageProcessor, MessageProcessor>();
builder.Services.AddHostedService<Worker>();

var host = builder.Build();
host.Run();
