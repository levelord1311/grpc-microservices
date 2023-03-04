module github.com/levelord1311/grpc-microservices/grpc-user-service

go 1.20

require (
	github.com/jmoiron/sqlx v1.3.5
	github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/user-service-api v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.53.0
)

require (
	github.com/envoyproxy/protoc-gen-validate v0.9.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/user-service-api => ./pkg/user-service-api
