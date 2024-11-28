# GRPC server

## GRPC dependencies
### go get google.golang.org/grpc
### go get google.golang.org/protobuf
### google.golang.org/genproto/googleapis/rpc

## Generate proto buf file 
### protoc --go_out=. --go-grpc_out=. ./grpc/*.proto; 


## execute cmd line
### go run main.go

## check your port is running or not 
### lsof -i :8080 