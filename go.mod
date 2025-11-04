module orderservice

go 1.24.0

toolchain go1.24.9

require (
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250929231259-57b25ae835d4 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250929231259-57b25ae835d4 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.5.1 // indirect
)

tool (
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
)
