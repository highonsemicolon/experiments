using System.Xml.Serialization;
using System.Text.Json.Serialization;


[XmlRoot("Person")]
public class Person
{
    public string? Name { get; set; }

    [XmlElement("Address")]
    public List<Address> Addresses { get; set; } = new();
}
public class Address
{
    [XmlElement("Street")]
    [JsonIgnore]
    public string? StreetRaw { get; set; }

    [XmlElement("Road")]
    [JsonIgnore]
    public string? RoadRaw { get; set; }

    public string? City { get; set; }

    [XmlIgnore]
    [JsonPropertyName("Street")]
    public string? Street => StreetRaw ?? RoadRaw;
}