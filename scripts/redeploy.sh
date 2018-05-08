#!/usr/bin/env bash
set -e
curl -sL --fail -H "token: ${DEPLOY_TOKEN}" 'http://goalkeeper.win.tue.nl:42424/redeploy'
