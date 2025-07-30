using System;
using ApiExtractor;

string rootPath = args.Length > 0 ? args[0] : "../tmp";
string outputPath = args.Length > 1 ? args[1] : "output.json";

var extractor = new EndpointDiscoveryService(rootPath, outputPath);
extractor.Run();
Console.WriteLine($"API extraction complete. Saved to {outputPath}");
