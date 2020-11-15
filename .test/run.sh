#!/bin/bash

printf "Removing running containers...\n"
docker-compose -f .test/docker-compose.test.yml down > /dev/null 2>&1
printf "Done.\n"

printf "Starting mysql database...\n"
docker-compose -f .test/docker-compose.test.yml up -d mysql

printf "Building go-there container...\n"
if ! docker-compose -f .test/docker-compose.test.yml build go-there;
then
  printf "Failed!\n"
  printf "Stopping all services...\n"
  docker-compose down

  exit 1
fi
printf "Done.\n"

printf "Starting go-there container...\n"
docker-compose -f .test/docker-compose.test.yml up -d go-there

printf "Waiting a bit for the initialization to finish"
for i in {1..10}
do
  sleep 1
  printf "."
done
printf "\n"

printf "Building go-there-test container...\n"
if ! docker-compose -f .test/docker-compose.test.yml build go-there-test;
then
  printf "Failed!\n"
  printf "Stopping all services...\n"
  docker-compose down

  exit 1
fi
printf "Done.\n"

printf "Starting go-there-test container...\n"
if ! docker-compose -f .test/docker-compose.test.yml run go-there-test;
then
  printf "Failed!\n"
  printf "Stopping all services...\n"
  docker-compose down

  exit 1
fi
printf "Done.\n"

printf "Cleaning up...\n"
docker-compose -f .test/docker-compose.test.yml down
printf "Success!\n"