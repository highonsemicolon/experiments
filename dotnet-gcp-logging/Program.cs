using dotnet_gcp_logging;
using Google.Cloud.Logging.Console;

Host.CreateDefaultBuilder(args)
    .ConfigureLogging(logging => {
        logging.ClearProviders();
        logging.AddGoogleCloudConsole();
    })
    .ConfigureServices(services => {
        services.AddHostedService<Worker>();
    })
    .Build()
    .Run();
