import { GreeterClient } from './grpc/GreeterServiceClientPb';
import { HelloRequest } from './grpc/greeter_pb';

const client = new GreeterClient('http://localhost:8081', null, null);

const request = new HelloRequest();
request.setName('Alice');

client.sayHello(request, {}, (err, response) => {
  if (err) {
    console.error('Error:', err.message);
  } else {
    console.log('Response:', response.getMessage());
  }
});
