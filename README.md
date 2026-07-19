# YShubham Tools Registry

The source behind [tools.yshubham.com](https://tools.yshubham.com): a tiny Go file host for checksum-verified binaries, one-command workstation setup, MCP configuration, and cross-agent skills.

Sandbox documentation is published for both humans and agents:

```text
https://tools.yshubham.com/docs/sandbox/
https://tools.yshubham.com/docs/sandbox/index.md
https://tools.yshubham.com/docs/sandbox/llms.txt
```

## Sandbox in one command

Install `sandbox`, `sandboxd`, `sandbox-mcp`, the `sandbox-platform` skill across supported agents, and register the local MCP server with detected Codex, Claude Code, and Gemini CLIs:

```sh
curl -fsSL https://tools.yshubham.com/sandbox/setup.sh | sh
```

The full setup requires `npx` for cross-agent skill installation. Binary-only installation has no Node.js dependency:

```sh
curl -fsSL https://tools.yshubham.com/sandbox/install.sh | sh
sandbox --help
```

The current v0.1.4 registry release includes macOS ARM64 and Linux x86-64 archives. The Sandbox release workflow is ready to add macOS x86-64 and Linux ARM64 without changing the installer. Release binaries bundle `libpq` and vendored OpenSSL, avoiding machine-specific PostgreSQL client paths.

Share a local frontend or API without running a Sandbox controller:

```sh
sandbox http 4321
```

The hosted first-party relay returns a temporary `https://local-….tunnel.yshubham.com` URL with HTTP, WebSocket, and Vite HMR support. Ctrl-C revokes the exact-host route.

Both commands are safe to rerun. The installer verifies the requested archive and embedded binary versions, stages all three binaries in the destination filesystem, and atomically replaces each existing command instead of leaving an older installation untouched.

Then connect the client to a self-hosted controller:

```sh
export SANDBOX_URL=https://sandbox.example.com
export SANDBOX_TOKEN='read-from-your-secret-store'
sandbox doctor
```

The setup script updates binaries in place while leaving existing `sandbox` MCP configuration entries unchanged. It also publishes verified templates for OpenCode, Cursor, VS Code/Copilot, Windsurf, Cline, Roo Code, Goose, Claude Desktop, Gemini Code Assist, and generic MCP hosts. Skip optional parts with:

```sh
curl -fsSL https://tools.yshubham.com/sandbox/setup.sh | sh -s -- --no-skill
curl -fsSL https://tools.yshubham.com/sandbox/setup.sh | sh -s -- --no-mcp
```

## Agent skills

Install the complete Sandbox workflow globally for every agent supported by the Skills CLI:

```sh
npx skills add bas3line/rool-repo --skill sandbox-platform --agent '*' --global --yes
```

Install Watchman the same way:

```sh
npx skills add bas3line/rool-repo --skill watchman --agent '*' --global --yes
```

Inspect before installing:

```sh
npx skills add bas3line/rool-repo --list
curl -fsSL https://tools.yshubham.com/skills
```

## Sandbox MCP

`sandbox-mcp` is installed with the Sandbox binary bundle. It is a local stdio bridge with 12 lifecycle tools, three resources, and two prompts.

```sh
codex mcp add sandbox -- sandbox-mcp
claude mcp add --scope user --transport stdio sandbox -- sandbox-mcp
gemini mcp add sandbox sandbox-mcp --scope user
curl -fsSL https://tools.yshubham.com/sandbox/clients/index.md
```

Keep `SANDBOX_TOKEN` in the environment or client secret store. Do not bake it into this repository or a shared MCP config.

## Registry layout

```text
public/
  index.html                    # human landing page
  skills.txt                    # curl-friendly skill commands
  sandbox/
    install.sh                  # checksum-verifying binary installer
    setup.sh                    # binaries + skill + detected MCP CLIs
    mcp.json                    # generic client template
    clients/                    # verified configs for major MCP hosts
    latest                      # current immutable release pointer
    releases/vX.Y.Z/            # archives and .sha256 files
  docs/sandbox/                 # themed portal + raw Markdown + llms.txt
  watchman/releases/vX.Y.Z/     # archives and .sha256 files
skills/
  sandbox-platform/             # complete Sandbox agent skill
  watchman/                     # complete Watchman agent skill
internal/server/                # hardened static HTTP behavior
```

The Go service serves regular files only. Mutable installers, setup files, config templates, and version pointers use `Cache-Control: no-cache`. Existing release assets use one-year immutable caching; missing release paths use `no-store` to prevent CDN cache poisoning during rollout.

## Publish a Sandbox release

The Sandbox release archive name is:

```text
sandbox_vX.Y.Z_<linux|darwin>_<amd64|arm64>.tar.gz
```

Each archive must contain exactly three regular files in this order:

```text
sandbox
sandboxd
sandbox-mcp
```

Copy each archive and its `.sha256` file into `public/sandbox/releases/vX.Y.Z/`. Publish all supported platform pairs before changing `public/sandbox/latest`; never overwrite an existing versioned asset.

Example from a Sandbox source checkout:

```sh
cargo build --profile dist --locked --package sandbox-cli --package sandboxd --package sandbox-mcp
version=v0.1.4
os=darwin
arch=arm64
archive="sandbox_${version}_${os}_${arch}.tar.gz"
tar -C target/dist -czf "$archive" sandbox sandboxd sandbox-mcp
shasum -a 256 "$archive" > "$archive.sha256"
mkdir -p "/path/to/rool-repo/public/sandbox/releases/$version"
cp "$archive" "$archive.sha256" "/path/to/rool-repo/public/sandbox/releases/$version/"
```

## Watchman

Install the current Watchman binary:

```sh
curl -fsSL https://tools.yshubham.com/watchman/install.sh | sh
watchman version
```

The root `install.sh` remains the canonical Watchman installer for backward compatibility. It requires HTTPS, verifies the release checksum and exact archive member set, and automatically verifies GitHub artifact provenance when a compatible `gh` CLI is available. Set `WATCHMAN_VERIFY_ATTESTATION=required` to make provenance verification mandatory. Linux packages are static musl executables without a build-runner GLIBC requirement. Watchman release archives retain the internal `gpu-watchman` member and install it as the public `watchman` command.

## Local run and verification

```sh
go test ./...
go run ./cmd/tools-host
curl -i http://127.0.0.1:8080/
curl -i http://127.0.0.1:8080/sandbox/setup.sh
curl -i http://127.0.0.1:8080/sandbox/mcp.json
curl -i http://127.0.0.1:8080/skills
```

Set `PORT` and `TOOLS_ROOT` when required:

```sh
PORT=9000 TOOLS_ROOT=./public go run ./cmd/tools-host
```

## Railway deployment

1. Deploy this repository from its root; Railway builds the `Dockerfile`.
2. Keep `/healthz` as the service health check.
3. Attach `tools.yshubham.com` and install the exact DNS record Railway provides.
4. Verify TLS and every mutable endpoint before publishing install commands.
5. Verify every release archive and checksum through the public domain after deployment.

The runtime is a non-root Alpine container containing one static Go binary plus the packaged public registry.
