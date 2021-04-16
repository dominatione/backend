#!/usr/bin/env bash

GO_MODULE="github.com/dominati-one/backend"

set -e
set -u
set -o pipefail

function generate_golang() {
  package="${1}"
  file="${2}"

  echo "Generating golang sources for ${file}@${package} ..."
  protoc --go_opt=module="${GO_MODULE}" --plugin="protoc-gen-go-grpc" --experimental_allow_proto3_optional --go_out=plugins=grpc:. "api/protoc/${package}/${file}.proto"
  protoc-go-inject-tag -input="pkg/protocol/${package}/${file}.pb.go"
}

function verify() {
  program="${1}"

  if ! command -v ${program} &> /dev/null
  then
      echo "ERROR: '${program}' binary could not be found"
      exit 1
  fi
}

function main() {
  echo "This script generates protobuf sources."

  verify "protoc"
  verify "protoc-gen-go-grpc"
  verify "protoc-go-inject-tag"

  generate_golang "entity" "planet"
  generate_golang "entity" "seed"

  generate_golang "component" "planet"
  generate_golang "component" "seed"
  generate_golang "component" "area_position"
  generate_golang "component" "area_tile"
  generate_golang "component" "area"
  generate_golang "component" "possession"

  generate_golang "blockchain" "block"
  generate_golang "blockchain" "event"
  generate_golang "blockchain" "event_create_planet"
  generate_golang "blockchain" "event_create_player"

  generate_golang "gameapi" "game_api_service"
  generate_golang "gameapi" "query_param_area_position"
  generate_golang "gameapi" "query_param_possession"
  generate_golang "gameapi" "create_planet_request_message"
  generate_golang "gameapi" "create_planet_response_message"
  generate_golang "gameapi" "get_planet_request"
  generate_golang "gameapi" "get_planet_response"
  generate_golang "gameapi" "get_planets_request"
  generate_golang "gameapi" "get_planets_response"
  generate_golang "gameapi" "get_seeds_request"
  generate_golang "gameapi" "get_seeds_response"
  generate_golang "gameapi" "get_area_tiles_request"
  generate_golang "gameapi" "get_area_tiles_response"

  echo -e "Done!"
}

main "${@:-}"