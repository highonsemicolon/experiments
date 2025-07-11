using dotnet_gcp_logging;

Host.CreateDefaultBuilder(args)
    .ConfigureServices(services => {
        services.AddHostedService<Worker>();
    })
    .Build()
    .Run();
