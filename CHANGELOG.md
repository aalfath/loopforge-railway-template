# Changelog

Each entry lists breaking changes first. Read it before you apply an update in Railway.

## 0.4.0-1 (2026-11-21)

No breaking changes. Adds the Deploy on Railway button (template `loopforge-1`) to the README. Nothing else changes; no update is needed.

## 0.4.0 (2026-11-21)

First release of this template.

- LoopForge [v0.4.0](https://github.com/Tencent/LoopForge/releases/tag/v0.4.0), image `ghcr.io/tencent/loopforge:0.4.0@sha256:87b6a3e3cba201d03f8f8d22929881f576021fc78c7c42e82f78329abbf4bda4`.
- MySQL `mysql:8.4.13@sha256:a81d043930e51745522bd5c430961cf8286f0a23f0a9a2a6e0108f7ea71b0902` on its own volume.
- The administrator account `admin` is created on the first start, with the generated `LOOPFORGE_ADMIN_PASSWORD`. Self sign-up is off (`AGENTFLOW_ALLOW_SIGNUP=false`).
- Upstream's starter catalog, welcome project and observability columns (`migrations/seed/*.sql`, `001`, `031`) are applied once per LoopForge version.
