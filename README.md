# PRISM-Backend

## gRPC code generation example
```
protoc \
--go_out=auth/generated \
--go_opt=paths=source_relative \
--go-grpc_out=auth/generated \
--go-grpc_opt=paths=source_relative \
auth.proto
```