---
name: watchman
description: Operate the Watchman CLI for NVIDIA GPU and AI inference-node inspection, vLLM/OpenAI-compatible endpoint validation, privacy-safe runtime discovery, model-capacity planning, active canaries, bounded saturation benchmarks, and offline regression gates. Use when an agent needs to diagnose a GPU server, inspect inference processes or model artifacts, validate a deployment, compare benchmark evidence, collect support evidence, or configure Watchman monitoring.
---

# Watchman

Use the `watchman` executable to collect typed, privacy-safe evidence from inference nodes. Prefer machine output, preserve nonclaims, and distinguish passive inspection from active traffic.

## Start safely

1. Run `watchman version` and `watchman --help` before choosing flags.
2. Require Watchman `0.8.0` or newer for the workflows in this skill. If the version is older or a subcommand is absent, stop and request an upgrade; do not invent a compatibility flag.
3. Read [references/commands.md](references/commands.md) for the applicable workflow and exact command pattern.
4. Use `--format json` for one report and `--format ndjson` for streams or automation.
5. Interpret exit codes with the emitted report. Exit `0` completed under the selected policy, `1` is a usage/setup/I/O failure, and `2` is a completed fail-closed policy result.

## Choose the workflow

| Need | Command |
| --- | --- |
| One node inventory and health report | `watchman snapshot` |
| Live terminal control room | `watchman top` |
| GPU process, VRAM, owner and container attribution | `watchman ps` |
| Driver, telemetry and endpoint validation | `watchman doctor` |
| Local vLLM/SGLang/TGI/Triton/TensorRT-LLM evidence | `watchman runtime inspect` |
| Safetensors metadata and byte validation | `watchman artifact inspect` |
| Model and KV-cache memory planning | `watchman capacity` |
| Small correctness and latency probe | `watchman canary` |
| Explicit closed-loop concurrency ladder | `watchman benchmark saturation` |
| Offline saturation regression gate | `watchman benchmark compare` |
| Offline canary rollout gate | `watchman rollout` |
| Historical or before/after analysis | `watchman history` / `watchman compare` |
| Prometheus and report service | `watchman serve` |
| Incident handoff artifact | `watchman bundle` |

## Preserve the safety boundary

- Treat `snapshot`, `top`, `ps`, `doctor`, `runtime`, `artifact`, `capacity`, `history`, `compare`, `rollout`, and `bundle` as passive or offline unless their explicit options say otherwise.
- Treat `canary` and `benchmark saturation` as active workloads. Obtain authorization for the endpoint, token cost, concurrency ladder, and test window before running them.
- Never run a saturation benchmark against unrelated production traffic without explicit approval. Start with `1,2,4,8`; expand only after reviewing errors, temperature, memory, queueing, and endpoint health.
- Prefer `--api-key-file`, `--prompt-file`, and file-backed service tokens. Do not place credentials, customer prompts, or sensitive expectations in process arguments.
- Keep the HTTP service on loopback unless the user supplies an authenticated TLS boundary and explicitly authorizes remote exposure.
- Do not reinterpret concurrency as server batch size, GPU occupancy, or proven capacity. Do not reinterpret endpoint token usage as raw GPU decode throughput.
- Do not claim statistical significance, causality, production capacity, cost, or an SLA from a single canary or saturation ladder.

## Return useful evidence

Report:

1. Exact command and Watchman version, with secrets omitted.
2. Target scope, UTC time window, model/workload identity, and tested concurrency points.
3. Typed report status, source completeness, failed gates, and exit code.
4. Latency, successful request rate, completion-token goodput, errors, and GPU evidence only when present.
5. Missing evidence and Watchman's stated nonclaims.
6. Recommended next test, keeping observed facts separate from inference.

Never paste credentials, prompts, generated output, raw environment values, URL queries, or arbitrary server diagnostics into a report.
