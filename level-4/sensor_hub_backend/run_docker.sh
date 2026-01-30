#!/bin/bash

docker stop sensor_hub
docker rm sensor_hub

docker build -t sensor_hub .

docker run -d --restart on-failure:5 --name sensor_hub --env-file .env -p 127.0.0.1:8282:8080/tcp sensor_hub
