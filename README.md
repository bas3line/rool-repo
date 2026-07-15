# Tools Host

Standalone Go service for serving installation scripts and versioned release artifacts from Railway.

## Public Layout

```text
https://tools.yshubham.com/watchman/install.sh
https://tools.yshubham.com/watchman/releases/v0.2.0/gpu-watchman_linux_amd64.tar.gz
https://tools.yshubham.com/watchman/releases/v0.2.0/gpu-watchman_linux_amd64.tar.gz.sha256
```

The Watchman installer is proxied from `bas3line/gpu-watchman` at request time. The `public/` directory holds only release artifacts. Create another tool by adding:

```text
public/exporter/releases/vX.Y.Z/
```

Use `templates/install.sh.tmpl` as the starting point for a new canonical installer, then configure a proxy route for it.

## Add GPU Watchman Artifacts

Build and package each target from the repository root:

```sh
mkdir -p web/public/watchman/releases/v0.2.0
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o gpu-watchman ./code/cmd/gpu-watchman
tar -czf web/public/watchman/releases/v0.2.0/gpu-watchman_linux_amd64.tar.gz gpu-watchman
shasum -a 256 web/public/watchman/releases/v0.2.0/gpu-watchman_linux_amd64.tar.gz > web/public/watchman/releases/v0.2.0/gpu-watchman_linux_amd64.tar.gz.sha256
rm gpu-watchman
```

Repeat for `linux_arm64`, `darwin_amd64`, and `darwin_arm64`. Each archive must contain one executable named `gpu-watchman`.

## Local Run

```sh
go run ./cmd/tools-host
curl http://localhost:8080/healthz
curl http://localhost:8080/watchman/install.sh
```

Set `PORT` and `TOOLS_ROOT` when required:

```sh
PORT=9000 TOOLS_ROOT=./public go run ./cmd/tools-host
```

## Railway Deployment

1. In Railway, create a project and deploy this repository.
2. Set the service root directory to `web`.
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
