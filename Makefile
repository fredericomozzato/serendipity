db_run:
	docker compose up -d

db_stop:
	docker compose down

db_root_connect:
	docker compose exec db mysql -u root -p

db_web_connect:
	docker compose exec db mysql -u web -p
