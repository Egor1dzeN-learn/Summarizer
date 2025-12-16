#!/usr/bin/env -S bash -euo pipefail

invoke_protoc() {
  uv run --with 'grpcio-tools==1.76.0' -m grpc_tools.protoc "$@"
}

generate_py_protos() {
  local -r project_root="${1}"
  local -r input="${2}"
  local -r output="${3}"

  rm -rf "${output}"
  mkdir -p "${output}"
  cp -r "${input}"/* "${output}"

  invoke_protoc \
    -I"${project_root}" \
    --python_out="${project_root}" \
    --pyi_out="${project_root}" \
    --grpc_python_out="${project_root}" \
    "${output}"/*.proto

  rm -f "${output}"/*.proto
}

generate_go_protos() {
  local -r project_root="${1}"
  local -r input="${2}"
  local -r output="${3}"

  export PATH="$PATH:$(go env GOPATH)/bin"
  if ! [[ -x "$(command -v protoc-gen-go)" ]]; then
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  fi

  rm -rf "${output}"
  mkdir -p "${output}"
  cp -r "${input}"/* "${output}"

  invoke_protoc \
    -I"${project_root}" \
    --go_out="${project_root}" \
    --go_opt=paths=source_relative \
    --go-grpc_out="${project_root}" \
    --go-grpc_opt=paths=source_relative \
    "${output}"/*.proto

  rm -f "${output}"/*.proto
}

declare -r rootdir="$(git rev-parse --show-toplevel)"

generate_py_protos \
  "${rootdir}/backend/worker/src" \
  "${rootdir}/backend/worker/protos" \
  "${rootdir}/backend/worker/src/worker/generated/protos"
generate_go_protos \
  "${rootdir}/backend" \
  "${rootdir}/backend/worker/protos" \
  "${rootdir}/backend/http/generated/protos"
