using System;
using System.Collections.Generic;
using System.IO;
using System.Linq; 
using Microsoft.CodeAnalysis;
using Microsoft.CodeAnalysis.CSharp;
using Microsoft.CodeAnalysis.CSharp.Syntax;
using System.Text.Json;
using System.Collections.Concurrent;
using System.Threading.Tasks;

namespace ApiExtractor
{
    public class EndpointExtractor(string rootPath, string outputPath)
    {
        private static readonly JsonSerializerOptions JsonOptions = new()
        {
            WriteIndented = true,
            PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
            // Encoder = System.Text.Encodings.Web.JavaScriptEncoder.UnsafeRelaxedJsonEscaping
        };
        private readonly ConcurrentBag<object> _allEndpoints = [];

        public void Run()
        {
            var csFiles = Directory.EnumerateFiles(rootPath, "*.cs", SearchOption.AllDirectories);
            Parallel.ForEach(csFiles, ProcessFile);

            var json = JsonSerializer.Serialize(_allEndpoints, JsonOptions);
            File.WriteAllText(outputPath, json);
        }

        private void ProcessFile(string csFile)
        {
            var code = File.ReadAllText(csFile);
            var tree = CSharpSyntaxTree.ParseText(code);
            var root = tree.GetRoot();
            var endpoints = new List<object>();

            var classNodes = root.DescendantNodes().OfType<ClassDeclarationSyntax>()
                .Where(c => c.BaseList?.Types.Any(t => t.ToString().Contains("Controller")) == true);

            foreach (var classNode in classNodes)
            {
                var classRoute = classNode.AttributeLists
                    .SelectMany(a => a.Attributes)
                    .FirstOrDefault(a => a.Name.ToString().Contains("Route"))
                    ?.ArgumentList?.Arguments.FirstOrDefault()?.ToString().Trim('"') ?? "[controller]";

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
                        Type = "Controller",
                        Controller = classNode.Identifier.ToString(),
                        Method = method.Identifier.ToString(),
                        HttpMethod = httpMethod,
                        Route = CombineRoutes(classRoute, methodRoute)
                    });
                }
            }

            var invocations = root.DescendantNodes().OfType<InvocationExpressionSyntax>();
            foreach (var call in invocations)
            {
                if (call.Expression is not MemberAccessExpressionSyntax expression) continue;

                var methodName = expression.Name.Identifier.Text;
                if (!methodName.StartsWith("Map", StringComparison.OrdinalIgnoreCase)) continue;

                var httpMethod = methodName.Replace("Map", "").ToUpper(); // GET, POST, etc.

                var argsList = call.ArgumentList?.Arguments;
                if (argsList == null || argsList.Value.Count == 0) continue;

                var routeArg = argsList.Value[0].ToString().Trim('"');
                endpoints.Add(new
                {
                    File = csFile,
                    Type = "Minimal",
                    HttpMethod = httpMethod,
                    Route = routeArg,
                    Handler = argsList.Value.Count > 1 ? argsList.Value[1].ToString() : "Lambda/Delegate"
                });
            }

            foreach (var ep in endpoints) _allEndpoints.Add(ep);
        }

        private static string CombineRoutes(string classRoute, string methodRoute)
        {
            if (string.IsNullOrEmpty(methodRoute)) return classRoute;
            return classRoute.TrimEnd('/') + "/" + methodRoute.TrimStart('/');
        }
    }
}
