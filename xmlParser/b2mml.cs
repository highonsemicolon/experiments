using System;
using System.Collections.Generic;
using System.Xml.Serialization;

[XmlRoot("ProductionSchedule", Namespace = "http://www.wbf.org/xml/B2MML-V0600")]
public class ProductionSchedule
{
    public string? ID { get; set; }
    public string? Description { get; set; }
    public DateTime PublishedDate { get; set; }

    [XmlElement("ProductionRequest")]
    public List<ProductionRequest>? ProductionRequests { get; set; }
}

public class ProductionRequest
{
    public string? ID { get; set; }
    public string? Description { get; set; }
    public DateTime RequestedStartTime { get; set; }
    public DateTime RequestedEndTime { get; set; }
    public string? ProductProductionRuleID { get; set; }
    public Quantity? Quantity { get; set; }
}

public class Quantity
{
    public string? QuantityString { get; set; }
    public string? DataType { get; set; }
    public string? UnitOfMeasure { get; set; }
}
