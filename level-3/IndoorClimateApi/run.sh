#!/bin/bash

API_PORT=8080
API_KEY=<replace>
ELASTIC_URL=<replace>
ELASTIC_API_KEY=<replace>

docker stop indoor-climate-level3
docker rm indoor-climate-level3

docker build --build-arg API_PORT=$API_PORT -t indoor-climate-level3 .

docker run --cpus 1 --env SensorApiKey=$API_KEY --env Elastic__Url=$ELASTIC_URL --env Elastic__ApiKey=$ELASTIC_API_KEY -p 127.0.0.1:92:8080/tcp --name indoor-climate-level3 --restart=on-failure:5 -d indoor-climate-level3