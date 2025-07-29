using System.Text.Json;
using System.Xml.Serialization;

var options = new JsonSerializerOptions { WriteIndented = true, PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower };

// string xmlPath = "sample.xml";

// // Deserialize XML into a Person object
// var serializer = new XmlSerializer(typeof(Person));
// using var reader = new StreamReader(xmlPath);
// var person = (Person)serializer.Deserialize(reader)!;

// // Convert the object to JSON using System.Text.Json
// var json = JsonSerializer.Serialize(person, options);

// Console.WriteLine(json);

string xmlPath = "b2mml.xml";

// Deserialize XML into a schedule object
var serializer = new XmlSerializer(typeof(ProductionSchedule));
using var reader = new StreamReader(xmlPath);
var schedule = (ProductionSchedule)serializer.Deserialize(reader)!;

// Convert the object to JSON using System.Text.Json
// var json = JsonSerializer.Serialize(schedule, options);

var map = Mapper.Map(schedule);
var json = JsonSerializer.Serialize(map, options);

Console.WriteLine(json);

