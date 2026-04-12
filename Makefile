up:
	docker compose up --build

down:
	docker compose down

test-messenger:
	cd messenger && go test ./... -v

psql:
	docker exec -it postgresql psql -U postgres -d messenger_db