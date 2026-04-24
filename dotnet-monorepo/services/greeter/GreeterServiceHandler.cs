using Grpc.Core;
using Greeter.V1;

namespace Handler;

public class GreeterServiceHandler : GreeterService.GreeterServiceBase
{
    public override Task<SayHelloResponse> SayHello(
        SayHelloRequest request,
        ServerCallContext context)
    {
        var response = new SayHelloResponse
        {
            Message = $"Hello, {request.Name}!"
        };

        return Task.FromResult(response);
    }
}