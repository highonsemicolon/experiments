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
