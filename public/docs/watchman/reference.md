# Watchman documentation

> Agent entry point for Watchman v0.8.2. Verify the installed binary before operating, use the narrowest workflow that answers the task, and preserve typed incomplete or `not_evaluable` evidence.

## Canonical endpoints

- Human documentation: <https://tools.yshubham.com/docs/watchman/>
- Raw Markdown: <https://tools.yshubham.com/docs/watchman/reference.md>
- Installer: <https://tools.yshubham.com/watchman/install.sh>
- Source: <https://github.com/bas3line/gpu-watchman>
- Agent skill: <https://github.com/bas3line/rool-repo/tree/main/skills/watchman>

## Install

Install the checksum-verified release as `watchman`:

```sh
curl -fsSL https://tools.yshubham.com/watchman/install.sh | sh
watchman version
```

Supported release targets are macOS and Linux on x86-64 and ARM64.

The installer requires HTTPS, verifies the published SHA-256 digest and exact archive member set, and automatically verifies GitHub artifact provenance when a compatible `gh` CLI is installed. Set `WATCHMAN_VERIFY_ATTESTATION=required` to fail unless provenance verification succeeds, or `disabled` to rely on the mandatory checksum only.

Install the complete Watchman skill globally for every agent supported by the Skills CLI:

```sh
npx skills add bas3line/rool-repo \
  --skill watchman --agent '*' --global --yes
```

Build from source when a hosted binary is unsuitable:

```sh
cargo build --locked --release
install -m 0755 target/release/gpu-watchman /usr/local/bin/watchman
```

## Agent operating protocol

1. Run `watchman version` and `watchman --help`.
2. Require version `0.8.2` or newer for this reference. If it is older, stop and upgrade.
3. Run `watchman <command> --help` before composing flags. Do not guess options.
4. Choose the narrowest workflow that answers the request.
5. Use `--format json` for a single automation artifact and `--format ndjson` for streams, history, or CI.
6. Record the exact redacted command, Watchman version, UTC window, target scope, workload identity, source completeness, typed status, failed or unavailable gates, exit code, and nonclaims.

## Choose a workflow

| Need | Command | Endpoint load |
| --- | --- | --- |
| Point-in-time GPU and endpoint report | `watchman snapshot` | Passive metrics probes only when selected |
| Live terminal control room | `watchman top` / `watchman watch` | Passive metrics probes only when selected |
| Prometheus, health and report service | `watchman serve` / `watchman exporter` | Passive metrics probes only when selected |
| GPU process ownership | `watchman ps` | None |
| Driver, telemetry and attribution checks | `watchman doctor` | Passive metrics probes only when selected |
| Local runtime fingerprint for explicit PIDs | `watchman runtime inspect` | None |
| Safetensors metadata and shard validation | `watchman artifact inspect` | None |
| Weights, KV cache and topology planning | `watchman capacity` | None |
| Small correctness and latency SLO probe | `watchman canary` | **Active inference traffic** |
| Bounded closed-loop concurrency ladder | `watchman benchmark saturation` | **Active inference traffic** |
| Offline saturation regression gate | `watchman benchmark compare` | None |
| Offline canary rollout gate | `watchman rollout` | None |
| Historical peaks and recurring findings | `watchman history` | None |
| Before/after node comparison | `watchman compare` | None |
| Portable incident evidence | `watchman bundle` | Passive metrics probes only when selected |
| Profiles | `watchman config init|validate|show` | None |

## Node triage

Collect a point-in-time report, process ownership and diagnostic evidence:

```sh
watchman snapshot --all --details --format json > snapshot.json
watchman ps --format json > processes.json
watchman doctor --format json > doctor.json
```

Add a loopback inference runtime when one exists:

```sh
watchman doctor --probe http://127.0.0.1:8000 --format json
```

Use `watchman top --watch 2s` for an interactive terminal. Prefer snapshot JSON for automation and incident handoff.

## Runtime, artifacts and capacity

Inspect bounded, privacy-safe evidence for explicit local server PIDs:

```sh
watchman runtime inspect --pid 4242 --pid 4243 --format json
```

An incomplete runtime report exits `2` by default. `--allow-incomplete` accepts best-effort evidence but does not make the evidence complete.

Validate safetensors metadata without loading tensor payloads:

```sh
watchman artifact inspect \
  /models/served-model/model.safetensors.index.json \
  --format json
```

Plan memory using explicit model, artifact and GPU assumptions:

```sh
watchman capacity \
  --model-config /models/served-model/config.json \
  --artifact /models/served-model/model.safetensors.index.json \
  --gpu-vram 24 --gpus 1 --tp 1 --weight-bits 16 \
  --format json
```

Capacity output is planning evidence. It does not prove runtime placement, allocator behavior, loaded residency, or that a model will start.

## Active canary

`canary` sends real chat-completion requests to an OpenAI-compatible endpoint. Obtain authorization for the endpoint, workload and request count first.

```sh
watchman canary \
  --base-url http://127.0.0.1:8000/v1 \
  --model served-model \
  --count 10 --concurrency 2 \
  --max-ttft 2s --max-e2e 10s \
  --min-success-percent 100 \
  --format json > canary.json
```

Use `--api-key-file PATH` for authenticated endpoints. Prefer `--prompt-file` for synthetic custom prompts and assign a non-secret `--workload-id`. Never place credentials or sensitive prompts in command arguments, reports, logs, or summaries.

## Saturation ladder

`benchmark saturation` sends real requests. Begin with a bounded ladder and review error, thermal, VRAM, queue and source evidence before expanding it.

```sh
watchman benchmark saturation \
  --base-url http://127.0.0.1:8000/v1 \
  --model served-model \
  --concurrency-stages 1,2,4,8 \
  --warmup-requests-per-worker 2 \
  --requests-per-worker 20 \
  --verify-concurrency 8 \
  --max-error-percent 1 \
  --format json > benchmark.json
```

This is a single-generator, closed-loop test. It can hide coordinated omission; the client or network can bottleneck first; repeated prompts may benefit from prefix caching; concurrency is not server batch size or GPU occupancy. The highest accepted tested point is not production capacity or an SLA certification.

## Offline deployment gates

Compare like-for-like saved saturation ladders without contacting either endpoint:

```sh
watchman benchmark compare baseline.json candidate.json \
  --max-p95-ttft-regression-percent 10 \
  --max-p95-e2e-regression-percent 10 \
  --min-successful-rps-ratio 0.95 \
  --min-completion-token-goodput-ratio 0.95 \
  --max-error-percent-increase 1 \
  --fail-on-regression \
  --format json
```

Compare compatible saved canaries:

```sh
watchman rollout baseline-canary.json candidate-canary.json \
  --max-p95-ttft-regression-percent 10 \
  --max-p95-e2e-regression-percent 10 \
  --fail-on-regression \
  --format json
```

Do not compare different workload identities, plans or policies as though they were equivalent. Preserve incompatible and `not_evaluable` outcomes.

## Monitoring service

Keep the listener on loopback and use a private token file:

```sh
watchman serve \
  --interval 5s \
  --listen 127.0.0.1:9400 \
  --api-token-file /run/secrets/watchman-api-token \
  --history /var/lib/watchman/history.ndjson
```

The service exposes `/livez`, `/metrics`, `/healthz`, and `/api/v1/report`. Remote exposure requires an authenticated, rate-limited TLS boundary. Do not use `--no-api-auth` outside an explicitly accepted loopback debugging run.

## Profiles and history

Profiles are explicit and never auto-discovered:

```sh
watchman config init watchman.toml
watchman --config watchman.toml config validate
watchman --config watchman.toml --profile production config show --format json
```

Summarize private NDJSON evidence or compare two node snapshots:

```sh
watchman history /var/lib/watchman/history.ndjson --format json
watchman compare before.json after.json --fail-on-regression --format json
```

Use `watchman rollout` for canary reports and `watchman benchmark compare` for saturation reports. The formats are intentionally separate.

## Evidence and exit codes

| Code | Meaning |
| --- | --- |
| `0` | The workflow completed and every selected policy passed. |
| `1` | Usage, validation, setup, local I/O, collection or encoding prevented a valid workflow. |
| `2` | Completed fail-closed evidence: an SLO or health policy failed, evidence was incomplete, capacity did not fit, or an enforced comparison regressed or was not evaluable. Inspect the emitted report. |

Hardware reports use `schema_version: 3`. Runtime inspection uses `runtime_fingerprint_version: 1`; canaries use `canary_version: 2`; saturation reports use `saturation_benchmark_version: 1`; rollout results use `rollout_version: 1`.

Never turn missing samples, unavailable telemetry, incomplete usage, an incompatible identity, a zero baseline or `not_evaluable` evidence into a pass.

## Safety boundary

- Passive and offline workflows do not authorize active inference traffic. `canary` and `benchmark saturation` do.
- Use synthetic, non-sensitive prompts. Treat saved reports, histories and bundles as private operational artifacts.
- Prefer file-backed credentials. Never expose API keys, service tokens, prompt content, generated output, URL queries or raw environment values.
- Non-loopback cleartext HTTP requires explicit acceptance. Prefer HTTPS.
- Watchman does not change clocks, power limits, compute modes, MIG state or workloads.
- A Watchman result is evidence under its stated plan and policy, not proof of causality, maximum production capacity, cost, statistical significance or SLA certification.

## Full references

- [Complete CLI reference](https://github.com/bas3line/rool-repo/blob/main/skills/watchman/references/cli.md)
- [Command recipes](https://github.com/bas3line/rool-repo/blob/main/skills/watchman/references/commands.md)
- [Operational profiles](https://github.com/bas3line/rool-repo/blob/main/skills/watchman/references/config.md)
- [Report schemas](https://github.com/bas3line/rool-repo/blob/main/skills/watchman/references/report-schema.md)
- [History](https://github.com/bas3line/rool-repo/blob/main/skills/watchman/references/history.md)
- [Node comparison](https://github.com/bas3line/rool-repo/blob/main/skills/watchman/references/comparison.md)
- [Canary rollout gates](https://github.com/bas3line/rool-repo/blob/main/skills/watchman/references/rollout.md)
- [Benchmark comparison gates](https://github.com/bas3line/rool-repo/blob/main/skills/watchman/references/benchmark-comparison.md)

Last verified against Watchman `0.8.2`.
