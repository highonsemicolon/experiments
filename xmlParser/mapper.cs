public class Mapper
{
    public static MongoProductionSchedule Map(ProductionSchedule schedule)
    {
        var mongoModel = new MongoProductionSchedule
        {
            ScheduleId = schedule.ID,
            PublishedDate = schedule.PublishedDate,
            Requests = schedule.ProductionRequests.Select(r => new MongoProductionRequest
            {
                RequestId = r.ID,
                StartTime = r.RequestedStartTime,
                Quantity = double.Parse(r.Quantity.QuantityString),
                Unit = r.Quantity.UnitOfMeasure
            }).ToList()
        };

        return mongoModel;
    }

}