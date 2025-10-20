FROM    alpine:latest AS build
WORKDIR /app
RUN     apk add --no-cache protoc
COPY    ./proto/greeter.proto .
RUN     protoc --include_imports --include_source_info \
            --descriptor_set_out=greeter.pb greeter.proto

FROM    envoyproxy/envoy:v1.22.0
COPY    --from=build /app/greeter.pb /tmp/