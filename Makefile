db_run:
	docker compose up -d

db_stop:
	docker compose down

db_connect:
	docker compose exec db psql -h localhost -U serendipity -d serendipity -W
