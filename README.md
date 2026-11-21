# LoopForge on Railway

[![Deploy on Railway](https://railway.com/button.svg)](https://railway.com/new/template/loopforge-1?utm_medium=integration&utm_source=button&utm_campaign=loopforge)

A Railway template for [LoopForge](https://github.com/Tencent/LoopForge) by Tencent: a self-hosted, browser-based workspace that runs coding agents (CodeBuddy, Codex, Claude Code and others) through repeatable SOP workflows with scripts and human approval gates.

**Version:** LoopForge **0.4.0** (`ghcr.io/tencent/loopforge:0.4.0@sha256:87b6a3e3cba201d03f8f8d22929881f576021fc78c7c42e82f78329abbf4bda4`) with MySQL **8.4.13**.

## What it deploys

| Service | Source | Volume | Notes |
|---|---|---|---|
| LoopForge | this repo (`Dockerfile`) | `/data` | Public, health check `/api/health` |
| MySQL | `mysql:8.4.13` (pinned by digest) | `/var/lib/mysql` | Private network only |

The `Dockerfile` starts from upstream's official image and adds only:

- `start.sh`: maps Railway's `PORT`, keeps all LoopForge data directories on `/data`, creates the administrator on the first start, and applies upstream's starter SQL once per version.
- `lfdb`: a small Go helper that applies upstream's unchanged `migrations/seed/*.sql`, `001_run_obs_generated_columns.sql` and `031_generic_backend_openssh.sql` (the official image has no `mysql` client, which upstream's installer uses).

## Variables

| Service | Variable | Default | Purpose |
|---|---|---|---|
| LoopForge | `LOOPFORGE_ADMIN_PASSWORD` | generated | Password of the administrator created on first start. Find it in the service's Variables tab |
| LoopForge | `LOOPFORGE_ADMIN_USER` | `admin` | Login of that administrator |
| LoopForge | `AGENTFLOW_SESSION_SECRET` | generated | Signs browser sessions. Keep it stable |
| LoopForge | `AGENTFLOW_MYSQL_DSN` | reference to MySQL | Database connection |
| LoopForge | `PORT` | `8080` | HTTP port |
| MySQL | `MYSQL_ROOT_PASSWORD`, `MYSQL_PASSWORD` | generated | Database passwords |
| MySQL | `MYSQL_DATABASE`, `MYSQL_USER` | `loopforge` | Database and its user |

The deployer is asked for nothing. Changing `LOOPFORGE_ADMIN_PASSWORD` after the first start does not change the password: change it in LoopForge (Platform administration → member accounts).

Optional settings you can add to the LoopForge service:

- `AGENTFLOW_ALLOW_SIGNUP=true` lets visitors create their own (non-admin) accounts. Off by default, because LoopForge lets every account see every project.
- `AGENTFLOW_GITHUB_CLIENT_ID` for GitHub sign-in and repository access. See upstream's [deployment guide](https://github.com/Tencent/LoopForge/blob/v0.4.0/docs/deployment.md).

## AI runs and execution machines

Railway gives a service no Docker socket. Upstream documents this mode: without the socket, the web UI, API, SOP designer, approvals and script-only SOPs (like the welcome SOP) run on the platform, and AI SOPs and conversation workspaces run on a connected execution machine. Connect one under **Execution machines** in LoopForge: it gives you a short-lived binding code and an install command for a Linux or macOS host with Docker or the agent CLIs. Model keys go in each user's personal settings.

## How updates reach you

This template builds from this repository's `main` branch. When we bump the LoopForge version, we add a `CHANGELOG.md` entry and merge to `main`; Railway then shows an "update available" notice on your project, and you choose when to apply it. Your data stays on the volumes. Read the changelog first: some LoopForge releases add SQL migrations for existing installations, and those are listed there.

## License

The wrapper files in this repository are MIT. LoopForge is MIT licensed by Tencent (see its [LICENSE](https://github.com/Tencent/LoopForge/blob/v0.4.0/LICENSE)); the files in `sql/` are unchanged copies from LoopForge v0.4.0.
