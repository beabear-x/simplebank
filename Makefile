runsimplebankdb:
	docker run --name simplebankdb --network bank_network -p 3307:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql:latest

connectdb:
	docker exec -it simplebankdb bash

migrateup:
	migrate -source "file://db/migration" -database "mysql://root:root@tcp(localhost:3307)/simple_bank" -verbose up

migrateup1:
	migrate -source "file://db/migration" -database "mysql://root:root@tcp(localhost:3307)/simple_bank" -verbose up 1

migratedown:
	migrate -source "file://db/migration" -database "mysql://root:root@tcp(localhost:3307)/simple_bank" -verbose down

migratedown1:
	migrate -source "file://db/migration" -database "mysql://root:root@tcp(localhost:3307)/simple_bank" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/beabear/simplebank/db/sqlc Store

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto

evans:
	evans -r repl

.PHONY: runsimplebankdb connectdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock proto evans