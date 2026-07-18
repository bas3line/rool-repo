# Tools Host

Standalone Go service for serving Watchman installation, versioned release artifacts, and a cross-agent skill landing page from Railway.

## Public Layout

```text
https://tools.yshubham.com/watchman/install.sh
https://tools.yshubham.com/skills
https://tools.yshubham.com/watchman/releases/v0.8.0/gpu-watchman_linux_amd64.tar.gz
https://tools.yshubham.com/watchman/releases/v0.8.0/gpu-watchman_linux_amd64.tar.gz.sha256
```

The container packages this repository's canonical `install.sh` and serves it without a runtime GitHub dependency. Release archives retain the internal `gpu-watchman` member for compatibility; the installer verifies and installs it as the public `watchman` command. The `skills/watchman/` package follows the Agent Skills `SKILL.md` convention, and `/skills` returns only plain-text install commands.

Install the binary:

```sh
curl -fsSL https://tools.yshubham.com/watchman/install.sh | sh
watchman version
```

Install the Watchman skill globally into every agent supported by the Skills CLI:

```sh
npx skills add bas3line/rool-repo --skill watchman --agent '*' --global --yes
```

Review `skills/watchman/SKILL.md` before installation. The hosted CLI and complete agent skill both target Rust Watchman v0.8.0 and the skill tells agents to stop on an older command surface.

Create another hosted tool by adding:

```text
public/exporter/releases/vX.Y.Z/
```

Use `templates/install.sh.tmpl` as the starting point for a new canonical installer, then configure a proxy route for it.

## Add GPU Watchman Artifacts

Build and package each target from the GPU Watchman repository root, then copy the immutable archives and checksum files here:

```sh
cargo build --release --locked
mkdir -p /path/to/rool-repo/public/watchman/releases/v0.8.0
tar -C target/release -czf /path/to/rool-repo/public/watchman/releases/v0.8.0/gpu-watchman_darwin_arm64.tar.gz gpu-watchman
shasum -a 256 /path/to/rool-repo/public/watchman/releases/v0.8.0/gpu-watchman_darwin_arm64.tar.gz > /path/to/rool-repo/public/watchman/releases/v0.8.0/gpu-watchman_darwin_arm64.tar.gz.sha256
```

Repeat for `linux_amd64`, `linux_arm64`, and `darwin_amd64`. Each archive must contain one executable named `gpu-watchman`.

## Local Run

```sh
go run ./cmd/tools-host
curl http://localhost:8080/healthz
curl http://localhost:8080/watchman/install.sh
curl http://localhost:8080/skills
```

Set `PORT` and `TOOLS_ROOT` when required:

```sh
PORT=9000 TOOLS_ROOT=./public go run ./cmd/tools-host
```

## Railway Deployment

1. In Railway, create a project and deploy this repository.
2. Use the repository root as the service root directory.
3. Railway detects the `Dockerfile`, starts the service on its supplied `PORT`, and checks `/healthz`.
4. Add the custom domain `tools.yshubham.com` in Railway's service settings.
5. Create exactly the DNS record Railway shows for that domain in the `yshubham.com` DNS zone. Wait for Railway verification and TLS issuance.

After deployment, the command is:

```sh
curl -fsSL https://tools.yshubham.com/watchman/install.sh | sh
```

## HTTP Behavior

- `GET` and `HEAD` are allowed.
- `GET /healthz` returns `ok`.
- Directory indexes and every other method return an error.
- `install.sh` uses `Cache-Control: no-cache`.
- Files under `/releases/` use immutable one-year caching. Publish a new versioned path for every release; never overwrite an existing release artifact.
