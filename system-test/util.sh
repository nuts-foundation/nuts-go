#!/usr/bin/env bash

function waitForDCService {
  SERVICE_NAME=$1
  printf "Waiting for Docker Compose service '%s' to become healthy" $SERVICE_NAME
  retry=0
  healthy=0
  while [ $retry -lt 30 ]; do
    status=$(docker inspect -f {{.State.Health.Status}} $(docker-compose ps -q $SERVICE_NAME))

    if [[ "$status" == "healthy" ]]; then
      healthy=1
      break
    fi

    printf "."
    sleep 0.5
    retry=$[$retry+1]
  done

  if [ $healthy -eq 0 ]; then
    echo "FAILED: Service took to long to start"
    exit 1
  fi
  echo ""
}
