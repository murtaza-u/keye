#!/usr/bin/env bash

exec protoc -I . \
    --go_out=paths=source_relative:. \
    --go-grpc_out=paths=source_relative:. \
    keye.proto
