---
name: watchman
description: Operate the complete Watchman CLI for NVIDIA GPU and inference-node observability, process attribution, vLLM/OpenAI-compatible telemetry, runtime fingerprints, safetensors inspection, model-capacity planning, canary SLOs, saturation benchmarks, benchmark and rollout regression gates, history, support bundles, profiles, Prometheus service operation, and shell completions. Use for GPU incident response, inference deployment validation, performance testing, capacity analysis, evidence comparison, monitoring setup, or any task involving the watchman command.
---

# Watchman

Use `watchman` as an evidence-first inference operations CLI. Match the exact installed command surface, choose the narrowest workflow that answers the question, and retain Watchman's typed limitations.

## Verify before operating

1. Run `watchman version` and `watchman --help`.
2. Require version `0.8.0` or newer for this skill. If older, report the installed version and request an upgrade.
3. Read the applicable section of [references/cli.md](references/cli.md) before composing a command. It is the complete option reference.
4. Run `watchman <command> --help` to confirm flags against the installed binary. Never guess or silently substitute an option.
5. Use `--format json` for a single automation artifact and `--format ndjson` for streams, history, or CI.

## Complete command map

| Task | Command |
| --- | --- |
| Point-in-time GPU, health, process and endpoint report | `watchman snapshot` |
| Live terminal control room | `watchman top` / `watchman watch` |
| Continuous Prometheus, health and report service | `watchman serve` / `watchman exporter` |
| GPU process, VRAM, owner, container and pod attribution | `watchman ps` |
| Driver, source, attribution and endpoint checks | `watchman doctor` |
| Local engine/framework/launch evidence for explicit PIDs | `watchman runtime inspect` |
| Safetensors metadata, byte, dtype and shard validation | `watchman artifact inspect` |
| Weights, topology, KV cache, headroom and concurrency planning | `watchman capacity` |
| Small active correctness, success and latency SLO probe | `watchman canary` |
| Bounded active closed-loop concurrency ladder | `watchman benchmark saturation` |
| Offline exact-ladder performance regression gate | `watchman benchmark compare` |
| Offline canary baseline/candidate rollout gate | `watchman rollout` |
| NDJSON availability, peak and recurring-finding summary | `watchman history` |
| Before/after hardware and telemetry comparison | `watchman compare` |
| Portable incident handoff evidence | `watchman bundle` |
| Create, validate and redact operational profiles | `watchman config init|validate|show` |
| Generate Bash, Zsh, Fish, PowerShell or Elvish completion | `watchman completions` |
| Print version | `watchman version` / `watchman --version` |

## Load only the needed reference

- Read [references/cli.md](references/cli.md) for every command, flag, environment variable, profile rule, output mode, and exit code.
- Read [references/config.md](references/config.md) when creating or reviewing `config_version = 1` profiles and precedence.
- Read [references/report-schema.md](references/report-schema.md) when parsing JSON/NDJSON or deciding what a field proves.
- Read [references/commands.md](references/commands.md) for common operator recipes.
- Read [references/history.md](references/history.md) for time-window, rotation and historical analysis.
- Read [references/comparison.md](references/comparison.md) for node before/after comparisons.
- Read [references/rollout.md](references/rollout.md) for canary rollout compatibility and formulas.
- Read [references/benchmark-comparison.md](references/benchmark-comparison.md) for saturation comparison formulas and CI behavior.

## Apply the safety boundary

- Passive and offline workflows do not authorize endpoint load. `canary` and `benchmark saturation` send synthetic requests and may consume serving capacity or billable tokens.
- Obtain explicit endpoint, workload, concurrency and time-window authorization before active testing. Start with a small ladder and expand only after reviewing errors, thermals, VRAM, queueing and source completeness.
- Prefer `--api-key-file`, `--prompt-file`, and file-backed service tokens. Never put credentials or sensitive prompts in command arguments, reports, logs, or summaries.
- Keep `serve` on loopback unless the user explicitly authorizes remote exposure with authentication and an external TLS boundary.
- Treat exit `0` as completed under the selected policy, exit `1` as usage/setup/I/O failure, and exit `2` as completed fail-closed evidence. Inspect the emitted report before explaining a nonzero result.
- Preserve `not_evaluable` and incomplete evidence. Never turn missing samples, incompatible identity, missing token usage, unavailable telemetry, or a zero baseline into a pass.
- Do not claim production capacity, causal attribution, statistical significance, GPU decode throughput, batch size, occupancy, cost, or SLA certification unless separate evidence proves it.

## Report results

Return the exact redacted command, Watchman version, UTC window, target scope, model/workload identity, tested stages, source completeness, typed status, failed or unavailable gates, exit code, and stated nonclaims. Separate observed facts from inference and recommend the smallest next test that resolves missing evidence.

Never expose credentials, prompt content, generated output, URL queries, raw environment values, arbitrary server diagnostics, filesystem identities omitted by the schema, or private artifact contents.
