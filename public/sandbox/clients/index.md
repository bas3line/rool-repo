# Connect sandbox-mcp to your agent

First install the binaries and cross-agent skill:

```sh
curl -fsSL https://tools.yshubham.com/sandbox/setup.sh | sh
export SANDBOX_URL=https://sandbox.example.com
export SANDBOX_TOKEN='read-from-your-secret-store'
```

The setup script automatically registers `sandbox-mcp` with detected Codex,
Claude Code, and Gemini CLIs. It never writes `SANDBOX_TOKEN` to their config.

## Direct CLI registration

```sh
# Codex
codex mcp add sandbox -- sandbox-mcp

# Claude Code
claude mcp add --scope user --transport stdio sandbox -- sandbox-mcp

# Gemini CLI
gemini mcp add sandbox sandbox-mcp --scope user

# OpenCode has an interactive registration wizard
opencode mcp add

# Goose for one session
goose session --with-extension "sandbox-mcp"

# VS Code user profile
code --add-mcp '{"name":"sandbox","type":"stdio","command":"sandbox-mcp"}'
```

## Drop-in templates

| Host | Config location | Template |
|---|---|---|
| Cursor | `~/.cursor/mcp.json` or `.cursor/mcp.json` | [`mcp-servers.json`](mcp-servers.json) |
| Claude Desktop | platform Claude Desktop config | [`mcp-servers.json`](mcp-servers.json) |
| Windsurf | `~/.codeium/windsurf/mcp_config.json` | [`mcp-servers.json`](mcp-servers.json) |
| Cline | `.cline/mcp.json` or the MCP settings UI | [`mcp-servers.json`](mcp-servers.json) |
| Roo Code | `.roo/mcp.json` or global MCP settings | [`mcp-servers.json`](mcp-servers.json) |
| Gemini Code Assist | `~/.gemini/settings.json` | [`mcp-servers.json`](mcp-servers.json) |
| OpenCode | `opencode.json` | [`opencode.json`](opencode.json) |
| VS Code / Copilot agent | `.vscode/mcp.json` | [`vscode.json`](vscode.json) |
| Goose | `~/.config/goose/config.yaml` | [`goose.yaml`](goose.yaml) |

Merge the named `sandbox` entry when a config file already exists. Do not
overwrite unrelated MCP servers. Use an absolute path to `sandbox-mcp` when a
desktop client does not inherit your shell `PATH`.

Keep the token in the environment or the client's secret-input mechanism.
Never commit it to a project MCP configuration. Pi, Aider, CommandCode, and
clients without native MCP support can use the installed `sandbox` CLI and the
same `sandbox-platform` skill.
