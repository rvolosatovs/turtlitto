#!/usr/bin/env bash
set -e

echo "${DOCKER_PASSWORD}" | docker login -u="${DOCKER_USERNAME}" --password-stdin
docker push rvolosatovs/srr:latest
if ! [ "${TRAVIS_TAG}" ]; then
  exit 0
fi

TAG=${TRAVIS_TAG#"v"}
for name in rvolosatovs/srr:${TAG} rvolosatovs/srr:${TAG:0:1} rvolosatovs/srr:${TAG:0:3} rvolosatovs/srr:${TAG:0:5}; do
  docker tag rvolosatovs/srr:latest ${name} && docker push ${name}
done
