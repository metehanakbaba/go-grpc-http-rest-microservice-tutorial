protoc --proto_path=api/proto/api --proto_path=third_party --go_out=plugins=grpc:pkg/api/api todo-service.proto
protoc --proto_path=api/proto/api --proto_path=third_party --grpc-gateway_out=logtostderr=true:pkg/api/api todo-service.proto
protoc --proto_path=api/proto/api --proto_path=third_party --swagger_out=logtostderr=true:api/swagger/api todo-service.proto
