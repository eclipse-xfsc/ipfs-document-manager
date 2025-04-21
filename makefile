include $(PWD)/.env

docker-compose-run:
	docker compose -f deployment/docker/docker-compose.yml -p $(PROJECT_NAME) --env-file=.env rm --force --stop
	docker-compose -f deployment/docker/docker-compose.yml -p $(PROJECT_NAME) --env-file=.env build --no-cache
	docker compose -f deployment/docker/docker-compose.yml -p $(PROJECT_NAME) --env-file=.env up --detach

restart:
	docker compose -f deployment/docker/docker-compose.yml -p $(PROJECT_NAME) --env-file=.env rm --force --stop server
	docker-compose -f deployment/docker/docker-compose.yml -p $(PROJECT_NAME) --env-file=.env build --no-cache server
	docker compose -f deployment/docker/docker-compose.yml -p $(PROJECT_NAME) --env-file=.env up --detach server