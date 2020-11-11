#!/bin/bash

docker-compose -f .test/docker-compose.test.yml down

docker-compose -f .test/docker-compose.test.yml up -d mysql

docker-compose -f .test/docker-compose.test.yml build go-there
docker-compose -f .test/docker-compose.test.yml up -d go-there

sleep 10

docker-compose -f .test/docker-compose.test.yml build go-there-test
docker-compose -f .test/docker-compose.test.yml run go-there-test

docker-compose -f .test/docker-compose.test.yml down