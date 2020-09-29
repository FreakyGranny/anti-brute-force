#!/bin/bash

set -e

trap "docker-compose -p integration-tests -f ./deployments/docker-compose.yaml -f ./deployments/docker-compose.tests.yaml down --rmi local" EXIT

docker-compose -p integration-tests -f ./deployments/docker-compose.yaml -f ./deployments/docker-compose.tests.yaml up -d --scale integration-tests=0
docker-compose -p integration-tests -f ./deployments/docker-compose.yaml -f ./deployments/docker-compose.tests.yaml run integration-tests
