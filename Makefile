up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

logs_web:
	docker compose logs -f web

logs_db:
	docker compose logs -f serendipity_db

rebuild:
	docker compose up -d --build

web_connect:
	docker compose exec web bash

db_connect:
	docker compose exec serendipity_db psql -h localhost -U serendipity -d serendipity -W

db_migrate:
	docker compose exec web go run cmd/migrate/main.go -direction=up

db_rollback:
	docker compose exec web go run cmd/migrate/main.go -direction=down
