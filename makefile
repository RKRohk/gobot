export GO111MODULE=on
export PATH="$PATH:$(go env GOPATH)/bin"v

setpath:
	export GO111MODULE=on
	export PATH="$PATH:$(go env GOPATH)/bin"v

protofiles: setpath
	protoc -I./  --go-grpc_out=./ --go_out=./ messageservice.proto   