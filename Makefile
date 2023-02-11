runsimplebankdb:
	docker run --name simplebankdb -p 3307:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql:latest

connectdb:
	docker exec -it simplebankdb bash

migrateup:
	migrate -source "file://db/migration" -database "mysql://root:root@tcp(localhost:3307)/simple_bank" -verbose up

migratedown:
	migrate -source "file://db/migration" -database "mysql://root:root@tcp(localhost:3307)/simple_bank" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: runsimplebankdb connectdb migrateup migratedown sqlc test server