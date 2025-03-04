proto-gen-auth:
	protoc --proto_path=auth-service/proto --go_out=auth-service/proto/pb --go_opt=paths=source_relative \
           --go-grpc_out=auth-service/proto/pb --go-grpc_opt=paths=source_relative \
           auth-service/proto/*.proto
proto-gen-user:
		protoc --proto_path=user-service/proto --go_out=user-service/proto/pb --go_opt=paths=source_relative \
		  --go-grpc_out=user-service/proto/pb --go-grpc_opt=paths=source_relative \
          user-service/proto/*.proto

.PHONY: proto-gen-auth proto-gen-user