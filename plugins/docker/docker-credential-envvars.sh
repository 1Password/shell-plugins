#!/bin/bash

# Credentials helper that returns values from environment variables:
# * DOCKER_CREDS_USERNAME: the username
# * DOCKER_CREDS_PASSWORD: the password / token
# https://docs.docker.com/engine/reference/commandline/login/
die() {
  echo "$@" 1>&2
  exit 1
}

case "$1" in
  get)
    read -r HOST
    if [ "$HOST" = "$REG" ]; then
      printf '{"ServerURL":"%s","Username":"%q","Secret":"%q"}\n' \
        "$HOST" "$DOCKER_CREDS_USR" "$DOCKER_CREDS_PSW"
    else
      die "No credentials available for $HOST"
    fi
    ;;
  *)
    die "Unsupported operation"
    ;;
esac