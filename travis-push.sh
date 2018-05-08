#!/usr/bin/env bash
set -xe

if [ ${TRAVIS_TAG} ]; then
    docker login -u="${DOCKER_USERNAME}" -p="${DOCKER_PASSWORD}"
    docker push rvolosatovs/srr:latest
    TAG=${TRAVIS_TAG#"v"}
    for name in rvolosatovs/srr:${TAG} rvolosatovs/srr:${TAG:0:1} rvolosatovs/srr:${TAG:0:3} rvolosatovs/srr:${TAG:0:5}; do
        docker tag rvolosatovs/srr:latest ${name} && docker push ${name}
    done
fi
