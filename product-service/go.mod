module github.com/best-microservice/product-service

go 1.24.2

require (
	github.com/best-microservice/common/protos/product v0.0.0-20231016123456-abcdef123456
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	golang.org/x/net v0.40.0
	google.golang.org/grpc v1.74.2
)

require (
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/best-microservice/common/protos/product => ../common/protos/product
