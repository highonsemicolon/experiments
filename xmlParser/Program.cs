using System.Text.Json;
using System.Xml.Serialization;


string xmlPath = "sample.xml";
var options = new JsonSerializerOptions { WriteIndented = true };

// Deserialize XML into a Person object
var serializer = new XmlSerializer(typeof(Person));
using var reader = new StreamReader(xmlPath);
var person = (Person)serializer.Deserialize(reader)!;

// Convert the object to JSON using System.Text.Json
var json = JsonSerializer.Serialize(person, options);

Console.WriteLine(json);

