.PHONY: all

all:
	docker-compose down
	docker-compose up -d
	# docker-compose exec go /bin/sh