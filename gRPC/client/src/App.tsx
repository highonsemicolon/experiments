import React, { useState } from "react"
import { GreeterClient } from "./proto/GreeterServiceClientPb"
import { HelloRequest } from "./proto/greeter_pb"

const greetClient = async (name: string) => {
  const EnvoyURL = "http://localhost:8000"
  const client = new GreeterClient(EnvoyURL)
  const request = new HelloRequest()
  request.setName(name)
  const response = await client.sayHello(request, {})
  console.log(response)
  const div = document.getElementById("response")
  if (div) div.innerText = response.getMessage()
}

function App() {
  const [name, setName] = useState("")
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
      <button onClick={onClickGreet}>greet</button>
      {name && <div id="response"></div>}
    </div>
  )
}

export default App
