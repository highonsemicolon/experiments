import React, { useState } from "react"
import { GreeterClient } from "./proto/greeter.client"
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport"

function App() {
  const [name, setName] = useState("")
  const [response, setResponse] = useState("")

  const greetClient = async (name: string) => {
    const EnvoyURL = "http://localhost:8000"
    const transport = new GrpcWebFetchTransport({ baseUrl: EnvoyURL })
    const client = new GreeterClient(transport)

    try {
      const res = await client.sayHello({ name })
      console.log(res.response)
      setResponse(res.response.message)
    } catch (err) {
      console.error("Error calling gRPC service:", err)
      setResponse("Error: " + (err as Error).message)
    }
  }

  const onClickGreet = () => {
    if (name) greetClient(name)
  }

  return (
    <div className="App">
      <input
        type="text"
        value={name}
        onChange={(e) => setName(e.target.value)}
      />
      <button onClick={onClickGreet}>Greet</button>
      <div id="response">{response}</div>
    </div>
  )
}

export default App
