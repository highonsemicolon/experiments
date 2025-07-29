public class MongoProductionSchedule
{
    public string? ScheduleId { get; set; }
    public DateTime? PublishedDate { get; set; }
    public List<MongoProductionRequest>? Requests { get; set; }
}

public class MongoProductionRequest
{
    public string? RequestId { get; set; }
    public DateTime? StartTime { get; set; }
    public double? Quantity { get; set; }
    public string? Unit { get; set; }
}
