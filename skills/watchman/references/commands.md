# Watchman command patterns

Read only the section needed for the current task. Confirm flags against `watchman <command> --help` because unsupported inputs must fail instead of being guessed.

For the full option surface, precedence rules, defaults, limits, aliases and environment variables, read [cli.md](cli.md).

## Node triage

Collect a point-in-time report and process attribution:

```sh
watchman snapshot --all --details --format json
watchman ps --format json
watchman doctor --format json
```

Add a loopback inference endpoint when it exists:

```sh
watchman doctor --probe http://127.0.0.1:8000 --format json
```

Use `watchman top --watch 2s` only for an interactive terminal. Prefer snapshot JSON for automation and handoff.

## Runtime fingerprint

Inspect explicit local server PIDs without storing paths, model names, environment values, or raw argv:

```sh
watchman runtime inspect --pid 4242 --pid 4243 --format json
```

An incomplete report exits `2` by default. Use `--allow-incomplete` only when the caller explicitly wants best-effort evidence; it does not make the evidence complete.

## Safetensors and capacity

Validate metadata without loading tensor payloads:

```sh
watchman artifact inspect /models/served-model/model.safetensors.index.json --format json
```

Plan memory with explicit geometry and hardware assumptions:

```sh
watchman capacity \
  --model-config /models/served-model/config.json \
  --artifact /models/served-model/model.safetensors.index.json \
  --gpu-vram 24 --gpus 1 --tp 1 --weight-bits 16 \
  --format json
```

Treat capacity output as planning evidence. It does not prove runtime placement, loaded residency, allocator behavior, or that the model will start.

## Canary

Run a small correctness and latency probe against a loopback OpenAI-compatible server:

```sh
watchman canary \
  --base-url http://127.0.0.1:8000/v1 \
  --model served-model \
  --count 10 --concurrency 2 \
  --max-ttft 2s --max-e2e 10s \
  --min-success-percent 100 \
  --format json
```

Use `--api-key-file PATH` for authenticated endpoints. Custom prompt content requires a non-secret `--workload-id`; prefer `--prompt-file` and never serialize the prompt.

## Saturation ladder

Begin with a bounded ladder and exact-stage deployment gate:

```sh
watchman benchmark saturation \
  --base-url http://127.0.0.1:8000/v1 \
  --model served-model \
  --concurrency-stages 1,2,4,8 \
  --warmup-requests-per-worker 2 \
  --requests-per-worker 20 \
  --verify-concurrency 8 \
  --max-error-percent 1 \
  --format json
```

Review the emitted stage evidence before increasing the ladder. Repeated prompts can benefit from prefix caching. The generator and network can bottleneck before the server.

## Offline benchmark comparison

Compare like-for-like saved ladders without contacting an endpoint:

```sh
watchman benchmark compare baseline.json candidate.json \
  --max-p95-ttft-regression-percent 10 \
  --max-p95-e2e-regression-percent 10 \
  --min-successful-rps-ratio 0.95 \
  --min-completion-token-goodput-ratio 0.95 \
  --max-error-percent-increase 1 \
  --fail-on-regression \
  --format ndjson
```

With enforcement, only `pass` exits `0`. A demonstrated regression or unavailable/incompatible required evidence exits `2`. Read/decode/version failures exit `1`.

## Canary rollout comparison

Compare saved canaries under explicit gates:

```sh
watchman rollout baseline-canary.json candidate-canary.json \
  --max-p95-ttft-regression-percent 10 \
  --max-p95-e2e-regression-percent 10 \
  --fail-on-regression \
  --format json
```

Do not compare different workload identities or policies as though they were equivalent.

## Monitoring service

Keep the listener on loopback and use a private token file:

```sh
watchman serve \
  --interval 5s \
  --listen 127.0.0.1:9400 \
  --api-token-file /run/secrets/watchman-api-token \
  --history /var/lib/watchman/history.ndjson
```

Do not use `--no-api-auth` outside an explicitly accepted local debugging run. Put remote access behind an authenticated, rate-limited TLS proxy.

## Profiles and configuration

Create a private starter profile, validate it without contacting hardware, then inspect only redacted normalized values:

```sh
watchman config init watchman.toml
watchman --config watchman.toml config validate
watchman --config watchman.toml --profile production config show --format json
```

Profiles are explicit and never auto-discovered. Values resolve as built-in defaults, selected profile, environment, then CLI. Read [config.md](config.md) before editing the schema or file permissions.

## History and node comparison

Summarize a private NDJSON history file:

```sh
watchman history /var/lib/watchman/history.ndjson --format json
```

Compare point-in-time hardware and telemetry reports:

```sh
watchman compare before.json after.json --fail-on-regression --format json
```

Use `watchman rollout` for canary reports and `watchman benchmark compare` for saturation reports. These formats are intentionally separate.

## Incident bundle

Collect one portable support artifact after confirming its selected evidence scope:

```sh
watchman bundle --output node-support.json
```

Treat bundles and histories as private operational artifacts. Inspect `watchman bundle --help` for probe, driver and source-selection flags before collection.

## Shell completions

Generate completion text without editing shell startup files:

```sh
watchman completions bash
watchman completions zsh
watchman completions fish
watchman completions powershell
watchman completions elvish
```

Write or source the generated output only when the user asks to configure that shell.

## Evidence files

- Use JSON for one pretty object.
- Use NDJSON for streams, history, and CI artifacts.
- Keep canary and benchmark reports private even though Watchman omits prompt and generated content.
- Compare closed-open UTC windows and retain the exact model, workload ID, plan, policy, and Watchman version.
- Preserve exit code `2` as a policy result; do not collapse it into an execution crash.
