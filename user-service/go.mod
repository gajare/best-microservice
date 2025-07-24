module github.com/best-microservice/user-service

go 1.24.2

require (
	github.com/best-microservice/common/protos/user v0.0.0-20241002120000-abcdef123456
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.40.0
	google.golang.org/grpc v1.74.2
)

require (
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/best-microservice/common/protos => ../common/protos

replace github.com/best-microservice/common/protos/user => ../common/protos/user
