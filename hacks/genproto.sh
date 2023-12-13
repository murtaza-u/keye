#!/usr/bin/env bash

exec protoc -I . \
    --go_out=paths=source_relative:./internal/pb \
    --go-grpc_out=paths=source_relative:./internal/pb \
    ./keye.proto
