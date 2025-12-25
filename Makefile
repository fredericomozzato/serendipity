up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

logs-web:
	docker compose logs -f web

logs-db:
	docker compose logs -f db

rebuild:
	docker compose up -d --build

db-connect:
	docker compose exec db psql -h localhost -U serendipity -d serendipity -W
