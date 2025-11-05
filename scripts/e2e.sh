#!/usr/bin/env bash
set -euo pipefail

# End-to-end test harness for the microservices stack.

usage() {
  cat <<'EOF'
Usage: e2e.sh -e EMAIL -p PASSWORD -k FIREBASE_API_KEY [options]

Options:
  -t TENANT_ID        Tenant identifier (default: tenant-e2e-<timestamp>)
  -T TENANT_NAME      Tenant display name (default: Tenant <TENANT_ID>)
  -n PROJECT_NAME     Docker Compose project name override
  -f COMPOSE_FILE     Path to docker-compose.yml (default: ../docker-compose.yml)
  -b API_BASE_URL     Base URL for the API gateway (default: http://localhost)
  -w TIMEOUT_SECONDS  Max seconds to wait for containers to become healthy (default: 420)
  -i POLL_SECONDS     Poll interval while waiting for readiness (default: 5)
  -K                  Keep the docker stack running after the script finishes
  -h                  Show this help and exit
EOF
}

log_info() {
  printf '[INFO] %s\n' "$*"
}

log_warn() {
  printf '[WARN] %s\n' "$*" >&2
}

log_error() {
  printf '[ERROR] %s\n' "$*" >&2
}

TARGET_SUBNET="10.200.0.0/24"
NETWORK_GATEWAY="10.200.0.1"
CONTAINER_NAMES=(api_gateway consul tenant_api tenant_db product_api product_db shipping_api shipping_db users_api users_db rabbitmq)
CONSUL_SERVICES=(users-service products-service tenant-service shipping-service)
BASE_PROJECT_NAME="proyecto-microservicios"
BASE_PROJECT_NETWORK="${BASE_PROJECT_NAME}_micro-network"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    log_error "Required command '$1' not found in PATH"
    exit 1
  fi
}

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
REPO_ROOT=$(cd "$SCRIPT_DIR/.." && pwd)
COMPOSE_FILE_DEFAULT="$REPO_ROOT/docker-compose.yml"

EMAIL=""
PASSWORD=""
API_KEY="${FIREBASE_API_KEY:-}"
TENANT_ID="acmes"
TENANT_NAME=""
PROJECT_NAME="pm-e2e-$(date +%s)"
WAIT_TIMEOUT=420
POLL_INTERVAL=5
API_BASE="${API_GATEWAY_URL:-http://localhost}"
COMPOSE_FILE="$COMPOSE_FILE_DEFAULT"
KEEP_STACK=0
CONSUL_API_URL="${E2E_CONSUL_API_URL:-http://localhost:8500}"
KONG_ADMIN_URL="${E2E_KONG_ADMIN_URL:-http://localhost:8090}"

while getopts ":e:p:k:t:n:f:w:i:b:T:Kh" opt; do
  case "$opt" in
    e) EMAIL="$OPTARG" ;;
    p) PASSWORD="$OPTARG" ;;
    k) API_KEY="$OPTARG" ;;
    t) TENANT_ID="$OPTARG" ;;
    n) PROJECT_NAME="$OPTARG" ;;
    f) COMPOSE_FILE="$OPTARG" ;;
    w)
      if [[ "$OPTARG" =~ ^[0-9]+$ ]]; then
        WAIT_TIMEOUT=$OPTARG
      else
        log_error "Invalid timeout value: $OPTARG"
        exit 1
      fi
      ;;
    i)
      if [[ "$OPTARG" =~ ^[0-9]+$ && "$OPTARG" -gt 0 ]]; then
        POLL_INTERVAL=$OPTARG
      else
        log_error "Invalid poll interval: $OPTARG"
        exit 1
      fi
      ;;
    b) API_BASE="$OPTARG" ;;
    T) TENANT_NAME="$OPTARG" ;;
    K) KEEP_STACK=1 ;;
    h)
      usage
      exit 0
      ;;
    :)
      log_error "Option -$OPTARG requires an argument"
      usage
      exit 1
      ;;
    \?)
      log_error "Unknown option -$OPTARG"
      usage
      exit 1
      ;;
  esac
done
shift $((OPTIND - 1))

if [[ -z "$EMAIL" || -z "$PASSWORD" ]]; then
  log_error "Email and password are required"
  usage
  exit 1
fi

if [[ -z "$API_KEY" ]]; then
  log_error "Firebase API key not provided (use -k or FIREBASE_API_KEY env var)"
  exit 1
fi

if [[ -z "$TENANT_NAME" ]]; then
  TENANT_NAME="Tenant ${TENANT_ID}"
fi

if [[ "$COMPOSE_FILE" != /* ]]; then
  COMPOSE_FILE="$(cd "$(dirname "$COMPOSE_FILE")" && pwd)/$(basename "$COMPOSE_FILE")"
fi

COMPOSE_DIR=$(dirname "$COMPOSE_FILE")

compose() {
  docker compose --project-name "$PROJECT_NAME" --project-directory "$COMPOSE_DIR" -f "$COMPOSE_FILE" "$@"
}

ensure_network_capacity() {
  local networks
  mapfile -t networks < <(docker network ls -q)
  for net in "${networks[@]}"; do
    [[ -z "$net" ]] && continue
    local subnets
    subnets=$(docker network inspect "$net" --format '{{range .IPAM.Config}}{{if .Subnet}}{{printf "%s " .Subnet}}{{end}}{{end}}' 2>/dev/null || true)
    [[ -z "$subnets" ]] && continue
    if [[ " $subnets " == *" $TARGET_SUBNET "* ]]; then
      local name
      name=$(docker network inspect "$net" --format '{{.Name}}' 2>/dev/null || echo "$net")
      local attached
      attached=$(docker network inspect "$net" --format '{{len .Containers}}' 2>/dev/null || echo 0)
      if [[ "$attached" == "0" ]]; then
        log_info "Removing orphan network '$name' with subnet $TARGET_SUBNET"
        docker network rm "$net" >/dev/null 2>&1 || true
      else
        log_error "Subnet $TARGET_SUBNET already in use by network '$name' with $attached attached container(s)."
        log_error "Stop that stack or adjust the subnet before running this script."
        exit 1
      fi
    fi
  done
}

ensure_container_names_available() {
  for name in "${CONTAINER_NAMES[@]}"; do
    local info
    info=$(docker ps -a --filter "name=^/${name}$" --format '{{.ID}} {{.Status}}' 2>/dev/null || true)
    [[ -z "$info" ]] && continue
    local id status
    id=$(cut -d' ' -f1 <<<"$info")
    status=$(cut -d' ' -f2- <<<"$info")
    if [[ "$status" == Up* ]]; then
      log_error "Container name '$name' already in use by running container $id."
      log_error "Stop the existing stack (e.g. docker compose down) or remove the container before running this script."
      exit 1
    fi
    log_warn "Removing leftover container '$name' (status: $status)"
    docker rm "$id" >/dev/null 2>&1 || true
  done
}

check_consul_service_passing() {
  local service="$1"
  local json
  json=$(curl -sS --fail "${CONSUL_API_URL}/v1/health/checks/${service}" 2>/dev/null || true)
  [[ -z "$json" ]] && return 1
  if CONSUL_RESPONSE="$json" python3 - "$service" <<'PY'
import json
import os
import sys

data = os.environ.get("CONSUL_RESPONSE", "")
if not data:
    sys.exit(1)
try:
    checks = json.loads(data)
except json.JSONDecodeError:
    sys.exit(1)

for check in checks:
    if check.get("Status") == "passing":
        sys.exit(0)

sys.exit(1)
PY
  then
    return 0
  fi
  return 1
}

wait_for_consul_leader() {
  local attempts=0
  local max_attempts=60
  while (( attempts < max_attempts )); do
    attempts=$((attempts + 1))
    local response
    response=$(curl -sS --fail "${CONSUL_API_URL}/v1/status/leader" 2>/dev/null || true)
    if [[ -n "$response" && "$response" != '""' ]]; then
      log_info "Consul leader elected: $response"
      return 0
    fi
    log_warn "Waiting for Consul leader (attempt ${attempts}/${max_attempts})"
    sleep "$POLL_INTERVAL"
  done
  log_error "Consul leader not available after ${max_attempts} attempts"
  return 1
}

wait_for_consul_services() {
  local services=("$@")
  local attempts=0
  local max_attempts=120
  while (( attempts < max_attempts )); do
    local missing=()
    for svc in "${services[@]}"; do
      if ! check_consul_service_passing "$svc"; then
        missing+=("$svc")
      fi
    done

    if (( ${#missing[@]} == 0 )); then
      log_info "Consul reports services healthy: ${services[*]}"
      return 0
    fi

    attempts=$((attempts + 1))
    log_warn "Waiting for Consul services (missing: ${missing[*]} | attempt ${attempts}/${max_attempts})"
    sleep "$POLL_INTERVAL"
  done

  log_error "Consul services did not become healthy in time: ${services[*]}"
  return 1
}

wait_for_kong_admin() {
  local attempts=0
  local max_attempts=60
  while (( attempts < max_attempts )); do
    attempts=$((attempts + 1))
    local http_code
    http_code=$(curl -sS -o /dev/null -w "%{http_code}" "${KONG_ADMIN_URL}/status" 2>/dev/null || true)
    if [[ "$http_code" == "200" ]]; then
      log_info "Kong admin API is ready"
      return 0
    fi
    log_warn "Waiting for Kong admin API (attempt ${attempts}/${max_attempts} | last status ${http_code:-n/a})"
    sleep "$POLL_INTERVAL"
  done
  log_error "Kong admin API did not become ready"
  return 1
}

cleanup_base_stack() {
  log_info "Cleaning base development stack remnants..."
  (
    cd "$REPO_ROOT"
    docker compose down --remove-orphans >/dev/null 2>&1 || true
  )
  docker ps -a --filter "network=${BASE_PROJECT_NETWORK}" -q | xargs -r docker rm -f >/dev/null 2>&1 || true
  docker images -f "dangling=true" -q | xargs -r docker rmi >/dev/null 2>&1 || true
  docker network rm "${BASE_PROJECT_NETWORK}" >/dev/null 2>&1 || true
  docker network prune -f >/dev/null 2>&1 || true
}

cleanup_project_stack() {
  local project_network="${PROJECT_NAME}_micro-network"
  log_info "Cleaning previous resources for project '${PROJECT_NAME}'..."
  compose down -v --remove-orphans >/dev/null 2>&1 || true
  docker ps -a --filter "network=${project_network}" -q | xargs -r docker rm -f >/dev/null 2>&1 || true
  docker network rm "${project_network}" >/dev/null 2>&1 || true
}

wait_for_stack() {
  local deadline=$((SECONDS + WAIT_TIMEOUT))
  while (( SECONDS < deadline )); do
    mapfile -t containers < <(compose ps -q)
    if (( ${#containers[@]} == 0 )); then
      sleep "$POLL_INTERVAL"
      continue
    fi

    local not_ready=0
    for container in "${containers[@]}"; do
      local status
      status=$(docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' "$container")
      if [[ "$status" != "healthy" && "$status" != "running" ]]; then
        not_ready=1
        break
      fi
    done

    if (( not_ready == 0 )); then
      return 0
    fi

    local snapshot
    snapshot=$(compose ps --format '{{.Name}}={{.State}}/{{.Health}}' | tr '\n' ' ')
    log_info "Waiting for containers... $snapshot"
    sleep "$POLL_INTERVAL"
  done

  compose ps
  return 1
}

create_tenant() {
  local payload
  payload=$(printf '{"tenant_id":"%s","tenant_name":"%s"}' "$TENANT_ID" "$TENANT_NAME")
  local attempts=0
  local max_attempts=40
  local tmp
  tmp=$(mktemp)

  while (( attempts < max_attempts )); do
    attempts=$((attempts + 1))
    set +e
    local http_code
    http_code=$(curl -sS -o "$tmp" -w "%{http_code}" -X POST "$API_BASE/api/tenants" \
      -H "Content-Type: application/json" -d "$payload")
    local exit_code=$?
    set -e

    if [[ $exit_code -ne 0 ]]; then
      log_warn "Attempt $attempts: tenant API unreachable (exit $exit_code)"
    else
      if [[ "$http_code" == "201" || "$http_code" == "409" ]]; then
        log_info "Tenant '$TENANT_ID' provisioned (HTTP $http_code)"
        rm -f "$tmp"
        return 0
      fi
      log_warn "Attempt $attempts: tenant creation returned HTTP $http_code: $(cat "$tmp")"
    fi
    sleep "$POLL_INTERVAL"
  done

  log_error "Unable to provision tenant '$TENANT_ID'"
  log_error "Last response: $(cat "$tmp")"
  rm -f "$tmp"
  return 1
}

fetch_firebase_token() {
  local payload
  payload=$(printf '{"email":"%s","password":"%s","returnSecureToken":true}' "$EMAIL" "$PASSWORD")
  local tmp
  tmp=$(mktemp)

  set +e
  local http_code
  http_code=$(curl -sS -o "$tmp" -w "%{http_code}" -X POST \
    "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=$API_KEY" \
    -H "Content-Type: application/json" -d "$payload")
  local exit_code=$?
  set -e

  if [[ $exit_code -ne 0 ]]; then
    log_error "Failed to contact Firebase Auth (exit $exit_code)"
    rm -f "$tmp"
    return 1
  fi

  if [[ "$http_code" != "200" ]]; then
    log_error "Firebase Auth error (HTTP $http_code): $(cat "$tmp")"
    rm -f "$tmp"
    return 1
  fi

  ID_TOKEN=$(python3 -c 'import json,sys; data=json.load(open(sys.argv[1])); print(data.get("idToken",""))' "$tmp")
  REFRESH_TOKEN=$(python3 -c 'import json,sys; data=json.load(open(sys.argv[1])); print(data.get("refreshToken",""))' "$tmp")
  rm -f "$tmp"

  if [[ -z "$ID_TOKEN" ]]; then
    log_error "Firebase Auth response did not include idToken"
    return 1
  fi

  log_info "Firebase token acquired"
  return 0
}

create_user_via_me() {
  local attempts=0
  local max_attempts=40
  local tmp
  tmp=$(mktemp)
  local url="$API_BASE/${TENANT_ID}/api/users/me"

  while (( attempts < max_attempts )); do
    attempts=$((attempts + 1))
    set +e
    local http_code
    http_code=$(curl -sS -o "$tmp" -w "%{http_code}" -X GET "$url" \
      -H "Authorization: Bearer $ID_TOKEN" \
      -H "X-Tenant-ID: $TENANT_ID" \
      -H "Accept: application/json")
    local exit_code=$?
    set -e

    if [[ $exit_code -ne 0 ]]; then
      log_warn "Attempt $attempts: /users/me unreachable (exit $exit_code)"
    else
      if [[ "$http_code" == "200" ]]; then
        USER_ID=$(python3 -c 'import json,sys; data=json.load(open(sys.argv[1])); print(data.get("id",""))' "$tmp")
        USER_ROLE=$(python3 -c 'import json,sys; data=json.load(open(sys.argv[1])); print(data.get("rol",""))' "$tmp")
        rm -f "$tmp"
        if [[ -z "$USER_ID" ]]; then
          log_error "User creation response missing id"
          return 1
        fi
        log_info "User ready (id=$USER_ID, role=$USER_ROLE)"
        return 0
      fi
      log_warn "Attempt $attempts: /users/me returned HTTP $http_code: $(cat "$tmp")"
    fi
    sleep "$POLL_INTERVAL"
  done

  log_error "Unable to create user via /users/me"
  rm -f "$tmp"
  return 1
}

create_product() {
  local payload
  payload=$(cat <<EOF
{"category":"fruta","price":5500,"description":"E2E product for tenant $TENANT_ID","name":"Mango E2E","stock":25,"unit":"kg","photo_url":"https://example.com/e2e-product.jpg"}
EOF
)
  local attempts=0
  local max_attempts=40
  local tmp
  tmp=$(mktemp)
  local url="$API_BASE/${TENANT_ID}/api/products"

  while (( attempts < max_attempts )); do
    attempts=$((attempts + 1))
    set +e
    local http_code
    http_code=$(curl -sS -o "$tmp" -w "%{http_code}" -X POST "$url" \
      -H "Authorization: Bearer $ID_TOKEN" \
      -H "X-Tenant-ID: $TENANT_ID" \
      -H "Content-Type: application/json" \
      -H "Accept: application/json" -d "$payload")
    local exit_code=$?
    set -e

    if [[ $exit_code -ne 0 ]]; then
      log_warn "Attempt $attempts: product API unreachable (exit $exit_code)"
    else
      if [[ "$http_code" == "201" ]]; then
        PRODUCT_ID=$(python3 -c 'import json,sys; data=json.load(open(sys.argv[1])); print(data.get("id",""))' "$tmp")
        rm -f "$tmp"
        if [[ -z "$PRODUCT_ID" ]]; then
          log_error "Product creation response missing id"
          return 1
        fi
        log_info "Product created (id=$PRODUCT_ID)"
        return 0
      fi
      if [[ "$http_code" == "401" || "$http_code" == "404" || "$http_code" == "409" || "$http_code" == "500" ]]; then
        log_warn "Attempt $attempts: product creation returned HTTP $http_code: $(cat "$tmp")"
      else
        log_error "Unexpected product creation status (HTTP $http_code): $(cat "$tmp")"
        rm -f "$tmp"
        return 1
      fi
    fi
    sleep "$POLL_INTERVAL"
  done

  log_error "Unable to create product"
  rm -f "$tmp"
  return 1
}

add_to_cart() {
  local payload
  payload=$(printf '{"product_id":"%s","quantity":1}' "$PRODUCT_ID")
  local attempts=0
  local max_attempts=40
  local tmp
  tmp=$(mktemp)
  local url="$API_BASE/${TENANT_ID}/api/cart/items"

  while (( attempts < max_attempts )); do
    attempts=$((attempts + 1))
    set +e
    local http_code
    http_code=$(curl -sS -o "$tmp" -w "%{http_code}" -X POST "$url" \
      -H "Authorization: Bearer $ID_TOKEN" \
      -H "X-Tenant-ID: $TENANT_ID" \
      -H "Content-Type: application/json" \
      -H "Accept: application/json" -d "$payload")
    local exit_code=$?
    set -e

    if [[ $exit_code -ne 0 ]]; then
      log_warn "Attempt $attempts: cart API unreachable (exit $exit_code)"
    else
      if [[ "$http_code" == "201" ]]; then
        CART_ITEM_ID=$(python3 -c 'import json,sys; data=json.load(open(sys.argv[1])); print(data.get("id",""))' "$tmp")
        rm -f "$tmp"
        if [[ -z "$CART_ITEM_ID" ]]; then
          log_error "Cart response missing id"
          return 1
        fi
        log_info "Product added to cart (cart_item_id=$CART_ITEM_ID)"
        return 0
      fi
      if [[ "$http_code" == "404" || "$http_code" == "409" || "$http_code" == "500" ]]; then
        log_warn "Attempt $attempts: cart add returned HTTP $http_code: $(cat "$tmp")"
      else
        log_error "Unexpected cart status (HTTP $http_code): $(cat "$tmp")"
        rm -f "$tmp"
        return 1
      fi
    fi
    sleep "$POLL_INTERVAL"
  done

  log_error "Unable to add product to cart"
  rm -f "$tmp"
  return 1
}

verify_cart_contents() {
  local tmp
  tmp=$(mktemp)
  set +e
  local http_code
  http_code=$(curl -sS -o "$tmp" -w "%{http_code}" -X GET "$API_BASE/${TENANT_ID}/api/cart" \
    -H "Authorization: Bearer $ID_TOKEN" \
    -H "X-Tenant-ID: $TENANT_ID" \
    -H "Accept: application/json")
  local exit_code=$?
  set -e

  if [[ $exit_code -ne 0 ]]; then
    log_warn "Could not fetch cart summary (exit $exit_code)"
    rm -f "$tmp"
    return
  fi

  if [[ "$http_code" != "200" ]]; then
    log_warn "Cart summary returned HTTP $http_code: $(cat "$tmp")"
    rm -f "$tmp"
    return
  fi

  local item_count
  item_count=$(python3 -c 'import json,sys; data=json.load(open(sys.argv[1])); print(len(data) if isinstance(data, list) else 0)' "$tmp")
  log_info "Cart now contains $item_count item(s)"
  rm -f "$tmp"
}

require_cmd docker
require_cmd curl
require_cmd python3

STACK_UP=0
cleanup() {
  if [[ "$STACK_UP" -eq 1 && "$KEEP_STACK" -eq 0 ]]; then
    log_info "Tearing down docker resources (project $PROJECT_NAME)"
    compose down -v --remove-orphans >/dev/null 2>&1 || true
  elif [[ "$STACK_UP" -eq 1 ]]; then
    log_warn "Leaving docker stack running (project $PROJECT_NAME)"
  fi
}
trap cleanup EXIT

log_info "Starting E2E flow with project '$PROJECT_NAME' and tenant '$TENANT_ID'"

cleanup_base_stack
cleanup_project_stack

ensure_network_capacity
ensure_container_names_available

log_info "Launching docker compose stack..."
if ! compose up --build -d; then
  log_warn "Initial docker compose up failed; retrying after forced network cleanup"
  docker network rm "${PROJECT_NAME}_micro-network" >/dev/null 2>&1 || true
  compose up --build -d
fi
STACK_UP=1

log_info "Waiting for services to become healthy (timeout ${WAIT_TIMEOUT}s)"
if ! wait_for_stack; then
  log_error "Services did not become healthy within $WAIT_TIMEOUT seconds"
  exit 1
fi
log_info "All containers report healthy"

if ! wait_for_consul_leader; then
  exit 1
fi

if ! wait_for_consul_services "${CONSUL_SERVICES[@]}"; then
  exit 1
fi

if ! wait_for_kong_admin; then
  exit 1
fi

if ! create_tenant; then
  exit 1
fi

sleep 5

if ! fetch_firebase_token; then
  exit 1
fi

if ! create_user_via_me; then
  exit 1
fi

if ! create_product; then
  exit 1
fi

if ! add_to_cart; then
  exit 1
fi

verify_cart_contents

log_info "E2E test flow completed successfully"
log_info "Tenant: $TENANT_ID"
log_info "User ID: $USER_ID (role: $USER_ROLE)"
log_info "Product ID: $PRODUCT_ID"
log_info "Cart item ID: $CART_ITEM_ID"
