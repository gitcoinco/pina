.PHONY: docker-build docker-run docker-kill docker-stop docker-logs docker-deploy-contracts docker-all

IMAGE_NAME=pina
CONTAINER_NAME=pina

docker-all: docker-kill docker-build docker-run

docker-build:
		docker build . -t $(IMAGE_NAME) --no-cache --progress=plain

docker-run:
		docker run --name $(CONTAINER_NAME) --rm -d -p 127.0.0.1:8000:8000/tcp $(IMAGE_NAME)

docker-kill:
		-docker kill $(CONTAINER_NAME)

docker-stop:
		docker stop $(CONTAINER_NAME)

docker-logs:
		docker logs -f $(CONTAINER_NAME)

