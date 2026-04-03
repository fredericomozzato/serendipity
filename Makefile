build:
	docker compose build

run:
	docker compose up

down:
	docker compose down

test:
	docker compose run --rm web bundle exec rspec
