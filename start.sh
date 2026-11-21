#!/bin/bash
# Start LoopForge on Railway: map Railway's PORT, keep every data directory on
# the /data volume, create the administrator once, apply upstream's starter SQL
# once per version, then run the platform.
set -euo pipefail

DATA="${LOOPFORGE_DATA_DIR:-/data}"
export AGENTFLOW_HTTP_PORT="${PORT:-${AGENTFLOW_HTTP_PORT:-8080}}"
export AGENTFLOW_AUTH_REQUIRED=true
export AGENTFLOW_ALLOW_SIGNUP="${AGENTFLOW_ALLOW_SIGNUP:-false}"
export AGENTFLOW_SECURE_COOKIES="${AGENTFLOW_SECURE_COOKIES:-true}"
export HOME="${DATA}/home"
export AGENTFLOW_ASSISTANT_SESSION_DIR="${DATA}/assistant"
export AGENTFLOW_CONVERSATION_DIR="${DATA}/conversations"
export AGENTFLOW_CLARIFY_SESSION_DIR="${DATA}/clarify"
export AGENTFLOW_RUNTIME_CLI_DIR="${DATA}/cli"
export AGENTFLOW_RUNNER_BINARY="${AGENTFLOW_RUNNER_BINARY:-/app/bin/agentflow-runner-{os}-{arch}}"
export AGENTFLOW_GATE_MCP_BINARY="${AGENTFLOW_GATE_MCP_BINARY:-/app/bin/agentflow-gate-mcp-{os}-{arch}}"
if [[ -n "${RAILWAY_PUBLIC_DOMAIN:-}" ]]; then
  # Execution machines and their sandboxes reach the platform at its public URL.
  export AGENTFLOW_PUBLIC_URL="${AGENTFLOW_PUBLIC_URL:-https://${RAILWAY_PUBLIC_DOMAIN}}"
  export AGENTFLOW_PLATFORM_MCP_URL="${AGENTFLOW_PLATFORM_MCP_URL:-https://${RAILWAY_PUBLIC_DOMAIN}}"
fi
mkdir -p "$HOME" "$AGENTFLOW_ASSISTANT_SESSION_DIR" "$AGENTFLOW_CONVERSATION_DIR" \
  "$AGENTFLOW_CLARIFY_SESSION_DIR" "$AGENTFLOW_RUNTIME_CLI_DIR" "${DATA}/workspaces"

if [[ -z "${AGENTFLOW_SESSION_SECRET:-}" ]]; then
  echo "AGENTFLOW_SESSION_SECRET is required" >&2
  exit 1
fi

lfdb wait

if [[ "$(lfdb accounts)" == "0" ]]; then
  # Without an account LoopForge lets the first browser visitor become the
  # administrator, so a public deploy creates it here instead.
  ADMIN_USER="${LOOPFORGE_ADMIN_USER:-admin}"
  if [[ -z "${LOOPFORGE_ADMIN_PASSWORD:-}" ]]; then
    echo "LOOPFORGE_ADMIN_PASSWORD is required for the first start (it creates the administrator)" >&2
    exit 1
  fi
  printf '%s\n' "$LOOPFORGE_ADMIN_PASSWORD" | /app/bin/loopforge-account -user "$ADMIN_USER" -admin
  echo "Created administrator account '${ADMIN_USER}' (password: the LOOPFORGE_ADMIN_PASSWORD variable)."
fi

marker="${DATA}/.loopforge-railway-sql-${LOOPFORGE_VERSION:-unknown}"
if [[ ! -f "$marker" ]]; then
  lfdb extensions /opt/loopforge-railway/sql
  touch "$marker"
fi

exec /app/bin/agentflow -workspaces "${DATA}/workspaces" "$@"
