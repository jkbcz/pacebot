#!/bin/bash

set -o allexport; source .env; set +o allexport

FILTER='{"contributionDate": {"_gte": "2025-12-31T23:00:00"}}'
FILTER_ENC=$(echo $FILTER | jq -Rr '@uri')

curl -X 'GET' \
  "https://api.bcc.no/contributions/Contributions?env=prod&limit=2000&filter=$FILTER_ENC&sort=contributionDate" \
  -H 'accept: application/json' \
  -H "authorization: Bearer $TOKEN" \
  | jq '.data | map(.contributionDate |= . + "Z")' > contributions.json


FILTER='{"amount": {"_gte": 0}, "clubSeasonId": {"_eq": 12}}'
FILTER_ENC=$(echo $FILTER | jq -Rr '@uri')

curl -X 'GET' \
  "https://api.bcc.no/contributions/Targets?env=prod&filter=$FILTER_ENC" \
  -H 'accept: application/json' \
  -H "authorization: Bearer $TOKEN" \
  | jq '.data' > targets.json

curl -X 'GET' \
  "https://api.bcc.no/v2/persons?limit=1000" \
  -H 'accept: application/json' \
  -H "authorization: Bearer $TOKEN" \
  | jq '.data' > persons.json