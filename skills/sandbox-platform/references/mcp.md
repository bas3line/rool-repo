# MCP tool and context map

## Setup

```sh
export SANDBOX_URL=https://sandbox.example.com
export SANDBOX_TOKEN='read-from-your-secret-store'
codex mcp add sandbox -- sandbox-mcp
```

`sandbox-mcp` is a local stdio bridge. It connects to the public controller API and does not need worker, database, Docker, or NATS access.

## Tools

| Tool | Use |
|---|---|
| `sandbox_health` | Verify controller reachability and version |
| `sandbox_create` | Create with resources, policy signals, labels, and placement |
| `sandbox_exec` | Execute argv; wait by default or return an operation |
| `sandbox_list` | List visible sandboxes, optionally by tenant |
| `sandbox_inspect` | Read one sandbox and its selected isolation |
| `sandbox_delete` | Remove runtime resources; wait by default |
| `sandbox_operation` | Read one asynchronous operation snapshot |
| `sandbox_wait` | Poll an operation with a bounded timeout |
| `sandbox_agent_list` | Discover built-in coding-agent profiles |
| `sandbox_agent_run` | Create from a coding-agent profile |

For `sandbox_create`, supply `tenant` and `image`. Optional fields cover startup `command`, non-secret `env`, CPU, memory, disk, PIDs, TTL, network, isolation, sensitivity, risk signals, labels, required worker labels, preferred region, and anti-affinity keys.

For `sandbox_exec`, supply `sandbox_id` and `argv`. Optional fields are `cwd`, non-secret `env`, `timeout_seconds`, and `wait`.

## Resources

- Read `sandbox://capabilities` before claiming a deployment feature exists.
- Read `sandbox://agents` to discover profile defaults without a controller call.
- Read `sandbox://workflow` for the compact lifecycle runbook.

## Prompts

- `sandbox-task` accepts `tenant`, `image`, `task`, and optional `network`.
- `sandbox-agent-session` accepts `agent`, `tenant`, and `task`.

Prompts guide workflow only. They do not grant authorization.

## Result handling

Inspect `isError` before consuming `structuredContent`. A non-zero remote exit is a tool error. Check the nested operation `state`, `error`, `output.exit_code`, `stdout`, `stderr`, and `truncated` fields. A protocol-success response does not mean the remote command succeeded.
