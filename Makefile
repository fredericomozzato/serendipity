db_run:
	docker compose up -d

db_stop:
	docker compose down

db_connect:
	docker compose exec db mysql -p
