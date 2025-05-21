module consumer-service

go 1.24

require (
	github.com/Laniakea00/e-commerce/proto v0.0.0
	github.com/nats-io/nats.go v1.42.0
	google.golang.org/grpc v1.71.1
)

require (
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/Laniakea00/e-commerce/proto => ../proto
