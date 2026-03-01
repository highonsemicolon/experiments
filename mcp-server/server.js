import express from "express";
import { z } from "zod";
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StreamableHTTPServerTransport } from "@modelcontextprotocol/sdk/server/streamableHttp.js";

const app = express();
app.use(express.json());

function createMcpServer() {
  const mcpServer = new McpServer({
    name: "New MCP server",
    description: "A simple MCP server example",
    version: "1.0.0",
  });

  mcpServer.registerTool(
    "getWeather",
    { inputSchema: { city: z.string() } },
    async ({ city }) => ({
      content: [{ type: "text", text: `Weather in ${city} is 25°C and sunny.` }],
    })
  );

  return mcpServer;
}

app.post("/mcp", async (req, res) => {
  const transport = new StreamableHTTPServerTransport({
    sessionIdGenerator: undefined,
  });
  const mcpServer = createMcpServer();
  await mcpServer.connect(transport);
  await transport.handleRequest(req, res, req.body);
});

app.get("/mcp", async (req, res) => {
  const transport = new StreamableHTTPServerTransport({
    sessionIdGenerator: undefined,
  });
  const mcpServer = createMcpServer();
  await mcpServer.connect(transport);
  await transport.handleRequest(req, res);
});

app.delete("/mcp", (req, res) => {
  res.status(200).end();
});

app.listen(3001, () => {
  console.log("MCP server running at http://localhost:3001/mcp");
});