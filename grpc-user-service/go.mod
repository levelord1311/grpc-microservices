module github.com/levelord1311/grpc-microservices/grpc-user-service

go 1.20

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Masterminds/squirrel v1.5.3
	github.com/jackc/pgx/v5 v5.3.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/user-service-api v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/pressly/goose/v3 v3.9.0
	github.com/stretchr/testify v1.8.2
	google.golang.org/grpc v1.53.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.9.1 // indirect
	github.com/gammazero/deque v0.2.1 // indirect
	github.com/gammazero/workerpool v1.1.3 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	golang.org/x/tools v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/protobuf v1.29.1 // indirect
)

replace github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/user-service-api => ./pkg/user-service-api
