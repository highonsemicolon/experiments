#!/bin/bash
npx protoc \
  --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts \
  --ts_out=./src/proto \
  --proto_path=../proto \
  ../proto/*.proto
