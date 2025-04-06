# gRPC server example

### Prerequisites:

proto

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# https://github.com/google/protobuf/releases
# it's best to install it via releases
# place the binary in the PATH
#   `cp protoc ~/.local/bin`
# place the 'well known' types under '/usr/local/include/'
#   `sudo cp -r google /usr/local/include/`
protoc --go_out=. --go-grpc_out=. todo.proto

grpcurl -d '{"text": "vnj"}' -plaintext -proto todo.proto localhost:50051 todoPackage.Todo.createTodo
grpcurl -plaintext -proto todo.proto localhost:50051 todoPackage.Todo.readTodos
```
