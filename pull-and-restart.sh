#!/bin/sh

git pull && \
docker compose up -d --build && \
docker compose exec -it nginx nginx -s reload && \
docker compose logs -f
