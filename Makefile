up:
	docker compose up --build

down:
	docker compose down

test-messenger:
	cd messenger && go test ./... -v