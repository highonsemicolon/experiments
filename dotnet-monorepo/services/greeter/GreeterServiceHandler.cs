using System.Diagnostics;

using Greeter.Service.Logging;
using Greeter.V1;

using Grpc.Core;


namespace Handler;

public class GreeterServiceHandler : GreeterService.GreeterServiceBase {
    private readonly ILogger<GreeterServiceHandler> _logger;
    public GreeterServiceHandler(ILogger<GreeterServiceHandler> logger) {
        _logger = logger;
    }
    public override Task<SayHelloResponse> SayHello(
        SayHelloRequest request,
        ServerCallContext context) {
        var activity = Activity.Current;
        activity?.SetTag("user.name", request.Name);

        GreeterLogs.GreetingReceived(_logger, request.Name, context.Peer);
        var response = new SayHelloResponse {
            Message = $"Hello, {request.Name}!"
        };

        GreeterLogs.GreetingSent(_logger, request.Name);

        return Task.FromResult(response);
    }
}
