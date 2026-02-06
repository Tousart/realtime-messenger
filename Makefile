up:
	docker compose up --build

down:
	docker compose down -v

test-messenger:
	cd messenger && go test ./... -v