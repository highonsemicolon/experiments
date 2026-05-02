FROM bufbuild/buf:latest AS proto-builder
COPY ./proto /proto
COPY buf.* ./
RUN buf build -o descriptor.pb

FROM envoyproxy/envoy:v1.33-latest
COPY --from=proto-builder /descriptor.pb /etc/envoy/descriptor.pb
