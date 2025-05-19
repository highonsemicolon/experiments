declare -a services=(
  "service/controlcenter-bootstrap-lb 8080:443"
  "service/kafka-bootstrap-lb 9092:9092"
  "service/schemaregistry-bootstrap-lb 8081:443"
)

pids=()
for svc in "${services[@]}"; do
  read name ports <<< "$svc"
  kubectl -n confluent port-forward $name $ports &
  pids+=($!)
done

trap "kill ${pids[@]}" SIGINT SIGTERM
wait
