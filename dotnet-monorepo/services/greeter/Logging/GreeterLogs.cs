using Microsoft.Extensions.Logging;

namespace Greeter.Service.Logging;

internal static partial class GreeterLogs
{
    [LoggerMessage(
        EventId = 1001,
        Level = LogLevel.Information,
        Message = "Received greeting request for {Name} from {Peer}")]
    public static partial void GreetingReceived(
        ILogger logger,
        string name,
        string peer);

    [LoggerMessage(
        EventId = 1002,
        Level = LogLevel.Information,
        Message = "Sending greeting response for {Name}")]
    public static partial void GreetingSent(
        ILogger logger,
        string name);
}
