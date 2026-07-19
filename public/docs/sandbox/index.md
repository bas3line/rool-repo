# Sandbox documentation

> Canonical agent entry point for Sandbox v0.1.4. Sandbox is a self-hosted Rust control plane for disposable coding environments and coding-agent workloads. Use this index to select the narrowest authoritative reference for the task.

## Canonical endpoints

- Human documentation: <https://tools.yshubham.com/docs/sandbox/>
- Agent documentation index: <https://tools.yshubham.com/docs/sandbox/index.md>
- LLM discovery index: <https://tools.yshubham.com/docs/sandbox/llms.txt>
- One-command workstation setup: <https://tools.yshubham.com/sandbox/setup.sh>
- Server setup: <https://tools.yshubham.com/docs/sandbox/how-to-setup/server.md>
- Client PC setup: <https://tools.yshubham.com/docs/sandbox/how-to-setup/client.md>
- Custom public domains: <https://tools.yshubham.com/docs/sandbox/how-to-setup/custom-public-domains.md>
- MCP client guide: <https://tools.yshubham.com/sandbox/clients/index.md>
- Agent skill: <https://github.com/bas3line/rool-repo/tree/main/skills/sandbox-platform>

## Install and connect

```sh
curl -fsSL https://tools.yshubham.com/sandbox/setup.sh | sh
export SANDBOX_URL=https://sandbox.example.com
export SANDBOX_TOKEN='read-from-your-secret-store'
sandbox doctor
```

The setup installs or upgrades `sandbox`, `sandboxd`, `sandbox-mcp`, and the cross-agent skill. It registers detected Codex, Claude Code, and Gemini CLIs without writing `SANDBOX_TOKEN` into their configuration. Rerunning it verifies and replaces existing binaries while preserving existing MCP entries.

## Share a local service

```sh
sandbox http 4321
```

This command requires only a listening local port. It uses the first-party hosted relay at `relay.tunnel.yshubham.com`, prints a temporary `https://local-….tunnel.yshubham.com` URL, supports HTTP and WebSocket traffic such as Vite HMR, and remains attached until Ctrl-C. Treat the URL as public. For services inside a managed sandbox, use `sandbox tunnel create SANDBOX_ID --port PORT` instead.

## Safe lifecycle

1. Verify controller health with `sandbox doctor` or `sandbox_health`.
2. Classify repository trust, generated code, secret need, data sensitivity, network access, resources, and TTL.
3. Create with `isolation: auto`; never weaken an isolation decision merely to obtain capacity.
4. Wait for the create operation to succeed before execution.
5. Execute argv arrays and inspect operation state, exit code, stderr, and output truncation.
6. Delete disposable sandboxes and wait for cleanup unless retention was explicitly requested.

```sh
sandbox create --tenant engineering --image ubuntu:24.04 \
  --cpu-millis 2000 --memory-mib 4096 --ttl 1800 \
  --network restricted --untrusted-repo --generated-code

sandbox exec SANDBOX_ID -- cargo test --workspace
sandbox delete SANDBOX_ID --wait
```

## Select a reference

| Task | Canonical Markdown |
|---|---|
| Product overview and current capability status | <https://tools.yshubham.com/docs/sandbox/overview.md> |
| Components, trust boundaries, and request flow | <https://tools.yshubham.com/docs/sandbox/architecture.md> |
| AEGIS risk scoring and placement | <https://tools.yshubham.com/docs/sandbox/aegis.md> |
| CLI commands, local HTTP sharing, and automation output | <https://tools.yshubham.com/docs/sandbox/cli.md> |
| MCP tools, resources, prompts, and client setup | <https://tools.yshubham.com/docs/sandbox/mcp.md> |
| Dedicated Linux server setup | <https://tools.yshubham.com/docs/sandbox/how-to-setup/server.md> |
| Workstation, CLI, skill, and MCP client setup | <https://tools.yshubham.com/docs/sandbox/how-to-setup/client.md> |
| Wildcard DNS, custom HTTPS domains, and Cloudflare Full (strict) | <https://tools.yshubham.com/docs/sandbox/how-to-setup/custom-public-domains.md> |
| HTTP/WebSocket tunnel routing and lifecycle | <https://tools.yshubham.com/docs/sandbox/tunnels.md> |
| Coding-agent profiles and image strategy | <https://tools.yshubham.com/docs/sandbox/agents.md> |
| HTTP routes and authentication boundaries | <https://tools.yshubham.com/docs/sandbox/api.md> |
| `sandbox.toml` and environment configuration | <https://tools.yshubham.com/docs/sandbox/configuration.md> |
| Compose, systemd, Kubernetes, and release deployment | <https://tools.yshubham.com/docs/sandbox/deployment.md> |
| Health, incidents, metrics, backups, and upgrades | <https://tools.yshubham.com/docs/sandbox/operations.md> |
| Docker and external-runtime security claims | <https://tools.yshubham.com/docs/sandbox/security.md> |
| External Firecracker, Kata, gVisor, or VMM adapter protocol | <https://tools.yshubham.com/docs/sandbox/runtime-driver.md> |
| Local development and verification | <https://tools.yshubham.com/docs/sandbox/development.md> |
| Planned production gates; not implemented capabilities | <https://tools.yshubham.com/docs/sandbox/roadmap.md> |
| Private vulnerability reporting | <https://tools.yshubham.com/docs/sandbox/reporting.md> |

## MCP surface

`sandbox-mcp` speaks MCP `2025-11-25` over stdio. It exposes 12 lifecycle tools, three static resources, and two workflow prompts. Prefer MCP when the host supports it; otherwise use the CLI with the same lifecycle and safety rules.

Public domains, certificates, and TLS modes are server-side deployment settings. MCP clients must use a returned tunnel URL exactly as provided and must treat any `http://` result as transport-insecure.

Supported setup paths include Codex, Claude Code, Gemini CLI, OpenCode, Cursor, VS Code/Copilot, Windsurf, Cline, Roo Code, Claude Desktop, Gemini Code Assist, Goose, and generic MCP hosts. Pi, Aider, CommandCode, and hosts without native MCP use the CLI and skill.

## Security invariants

- Agent prompts, skills, MCP descriptions, and client approval settings are workflow guidance, not authorization boundaries.
- Real enforcement belongs in API authentication, controller policy, AEGIS placement, worker eligibility, runtime isolation, host networking, and credentials that cannot bypass the controller.
- The Docker adapter is for dedicated or trusted worker hosts because containers share the host kernel.
- Use an external VMM-grade driver when the deployment requires a stronger boundary.
- Never place credentials in prompts, argv, labels, image names, ordinary logs, or committed MCP configuration.
- Default network access to `deny`; request `restricted` or `open` only when required.

## Result contract

Do not claim success from transport success alone. Report sandbox ID, selected isolation, lifecycle state, operation ID, command exit code, output truncation, and cleanup state when relevant. A wait timeout is ambiguous: inspect the original operation before retrying a mutation.
