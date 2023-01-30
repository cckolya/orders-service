all: go-update pg-start nats-stream service-start

service-start:
	go run cmd/orders-service.go

go-update:
	go mod tidy

producer-start:
	go run cmd/orders-producer.go

pg-start:
	docker run --rm \
    --detach \
    --publish 5432:5432 \
    --env POSTGRES_DB=postgres \
    --env POSTGRES_USER=postgres \
    --env POSTGRES_PASSWORD=postgres \
    postgres

nats-stream:
	docker run --rm --detach -p 4222:4222 -p 8223:8223 nats-streaming -p 4222 -m 8223

wrk-test:
	wrk -c1 -t1 -d5s -s ./script.lua --latency http://localhost:8080
