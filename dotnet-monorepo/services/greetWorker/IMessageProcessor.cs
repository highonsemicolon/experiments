namespace Greet.Worker;

public interface IMessageProcessor {
    Task ProcessAsync(string message, CancellationToken cancellationToken);
}
