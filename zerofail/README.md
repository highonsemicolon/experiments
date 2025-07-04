'''
helm install my-mongo -n mongodb --create-namespace -f mongo-values.yaml bitnami/mongodb
'''

'''
protoc \
  --go_out=. --go_opt=paths=source_relative,Mproto/record.proto=record-service/proto \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative,Mproto/record.proto=record-service/proto \
  proto/record.proto
'''

## Golang in k8s
'''
kubectl run insert-job --image=golang --restart=Never -i --tty --namespace=mongodb -- bash
'''

'''
kubectl cp ./ mongodb/insert-job:/app
'''

'''
kubectl run grpc-client --rm -it --image=tarrunkhosla/grpcio:v1 --restart=Never
'''

'''
kubectl cp ./tests default/grpc-client:/go
grpcurl -plaintext -d @ 10.244.0.25:8080 record.RecordService/UpsertRecords < tests/insert.json
'''

## Mongosh in k8s
'''
kubectl run mongo-client --rm -it --image=mongo --restart=Never -- bash
'''

'''
mongosh mongodb://admin:admin@my-mongo-mongodb-headless.mongodb.svc.cluster.local:27017/?authSource=admin
'''
