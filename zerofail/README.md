'''
helm install my-mongo -n mongodb --create-namespace -f mongo-values.yaml bitnami/mongodb
'''

'''
kubectl run insert-job --image=golang --restart=Never -i --tty --namespace=mongodb -- bash
'''

'''
kubectl cp ./ mongodb/insert-job:/app
'''
