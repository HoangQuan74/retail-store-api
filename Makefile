# ───── Dev (local) ─────
run-api:
	go run cmd/api/main.go

run-socket:
	go run cmd/socket/main.go

run-consumer:
	go run cmd/consumer/main.go

build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/socket cmd/socket/main.go
	go build -o bin/consumer cmd/consumer/main.go

sqlc:
	sqlc generate

migrate-up:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/retail_store?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/retail_store?sslmode=disable" -verbose down

# ───── Docker ─────
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker build --build-arg SERVICE=api -t retail-store-api:latest .
	docker build --build-arg SERVICE=socket -t retail-store-socket:latest .
	docker build --build-arg SERVICE=consumer -t retail-store-consumer:latest .

docker-infra:
	docker-compose up -d postgres redis nats elasticsearch logstash kibana

# ───── K8s ─────
k8s-dev:
	kubectl apply -k deploy/k8s/overlays/dev

k8s-prod:
	kubectl apply -k deploy/k8s/overlays/prod

k8s-delete-dev:
	kubectl delete -k deploy/k8s/overlays/dev

k8s-delete-prod:
	kubectl delete -k deploy/k8s/overlays/prod

.PHONY: run-api run-socket run-consumer build sqlc migrate-up migrate-down \
	docker-up docker-down docker-build docker-infra \
	k8s-dev k8s-prod k8s-delete-dev k8s-delete-prod
