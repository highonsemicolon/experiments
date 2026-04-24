using Grpc.Core;
using Greeter.V1;

namespace Handler;

public class GreeterServiceHandler : GreeterService.GreeterServiceBase
{
    private readonly ILogger<GreeterServiceHandler> _logger;
    public GreeterServiceHandler(ILogger<GreeterServiceHandler> logger) {
        _logger = logger;
    }
    public override Task<SayHelloResponse> SayHello(
        SayHelloRequest request,
        ServerCallContext context)
    {
        _logger.LogInformation(
            "Received greeting request for {Name} from {Peer}",
            request.Name,
            context.Peer);

        var response = new SayHelloResponse
        {
            Message = $"Hello, {request.Name}!"
        };

        _logger.LogInformation(
            "Sending greeting response for {Name} to {Peer}",
            request.Name,
            context.Peer);

        return Task.FromResult(response);
    }
}