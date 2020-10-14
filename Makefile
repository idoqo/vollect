restart:
	docker-compose stop api && docker-compose build api && docker-compose up --no-start api && docker-compose start api
