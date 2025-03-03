proto-gen-auth:
	protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           auth-service/proto/*.proto
proto-gen-user:
		protoc --go_out=./user-service/proto/pb --go_opt=paths=source_relative \
		  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
          ./user-service/proto/*.proto

.PHONY: proto-gen-auth proto-gen-user