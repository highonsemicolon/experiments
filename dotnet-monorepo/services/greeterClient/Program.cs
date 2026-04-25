using System;
using System.Threading.Tasks;
using Grpc.Net.Client;
using Greeter.V1;

class Program
{
    static async Task Main(string[] args)
    {        
        using var channel = GrpcChannel.ForAddress("http://localhost:8080");

        var client = new GreeterService.GreeterServiceClient(channel);

        var request = new SayHelloRequest
        {
            Name = "Alice"
        };

        var response = await client.SayHelloAsync(request);

        Console.WriteLine("Server response: " + response.Message);
    }
}