# How to build gRPC file
see: https://grpc.io/docs/languages/go/quickstart/

```
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative <Protocol Buffer file>.proto
```

## Requirements
```
$ go version
go version go1.19.2 darwin/amd64
$ protoc --version
libprotoc 3.21.12
```

## Tips
### data type for protoc
see: https://protobuf.dev/programming-guides/proto3/#scalar
| .proto | map to java |
| ---- | ---- |
| int32 | int |
| bool | boolean |
| string | String |

### Need to...
insert `go_package` option in each proto file like below..
```
option go_package = "ryo29wx/Caol-Ila/interfaces/recommend";
```