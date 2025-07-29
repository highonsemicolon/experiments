using System;
using System.Collections.Generic;
using System.IO;
using System.Linq; 
using Microsoft.CodeAnalysis;
using Microsoft.CodeAnalysis.CSharp;
using Microsoft.CodeAnalysis.CSharp.Syntax;
using System.Text.Json;

string rootPath = args.Length > 0 ? args[0] : "../tmp";
string outputPath = args.Length > 1 ? args[1] : "output.json";

var endpoints = new List<object>();

foreach (var csFile in Directory.EnumerateFiles(rootPath, "*.cs", SearchOption.AllDirectories))
{
    var code = File.ReadAllText(csFile);
    var tree = CSharpSyntaxTree.ParseText(code);
    var root = tree.GetRoot();

    var classNodes = root.DescendantNodes().OfType<ClassDeclarationSyntax>()
        .Where(c => c.BaseList?.Types.Any(t => t.ToString().Contains("Controller")) == true);

    foreach (var classNode in classNodes)
    {
        var classRoute = classNode.AttributeLists
            .SelectMany(a => a.Attributes)
            .Where(a => a.Name.ToString().Contains("Route"))
            .Select(a => a.ArgumentList?.Arguments.FirstOrDefault()?.ToString().Trim('"'))
            .FirstOrDefault() ?? "[controller]";

        foreach (var method in classNode.Members.OfType<MethodDeclarationSyntax>())
        {
            var httpAttr = method.AttributeLists
                .SelectMany(a => a.Attributes)
                .FirstOrDefault(a => a.Name.ToString().StartsWith("Http"));

            if (httpAttr == null) continue;

            var httpMethod = httpAttr.Name.ToString().Replace("Http", "").ToUpper();
            var methodRoute = httpAttr.ArgumentList?.Arguments.FirstOrDefault()?.ToString().Trim('"') ?? "";

            endpoints.Add(new
            {
                File = csFile,
                Controller = classNode.Identifier.ToString(),
                Method = method.Identifier.ToString(),
                HttpMethod = httpMethod,
                Route = CombineRoutes(classRoute, methodRoute)
            });
        }
    }
}

File.WriteAllText(outputPath, JsonSerializer.Serialize(endpoints, new JsonSerializerOptions { WriteIndented = true }));
Console.WriteLine($"API extraction complete. Saved to {outputPath}");

string CombineRoutes(string classRoute, string methodRoute)
{
    if (string.IsNullOrEmpty(methodRoute)) return classRoute;
    return classRoute.TrimEnd('/') + "/" + methodRoute.TrimStart('/');
}
