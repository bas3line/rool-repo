---
name: objects-storage
description: Upload, retrieve, inspect, and early-delete small temporary public build artifacts through objects.yshubham.com. Use when an agent must safely share a non-sensitive log bundle, benchmark, release artifact, report, or other public file up to 10 MiB with a seven-day lifetime.
---

# Objects Storage

## Overview

Use Objects for small, temporary artifacts that are safe for anyone with the returned URL to access. It is not a backup service, package registry, secret store, or access-controlled file system.

Read [the API reference](references/api.md) before relying on a status code or field not covered here.

## Keep uploads safe

- Only upload public, non-sensitive files. Never upload credentials, tokens, private source, personal data, database dumps, or material that needs access control.
- Confirm the local file is at most **10 MiB** before uploading.
- Do not attempt custom paths, directory listing, multipart form data, or a direct bucket endpoint. The public API intentionally does not provide them.
- Objects expire after **seven days**. A later `404` can mean either expiry or an unknown object.
- The returned `deleteToken` is a private capability. Never include it in a URL, repository, issue, chat summary, or build log.

## Upload an artifact

1. Choose an accurate content type. Use `application/octet-stream` only when the file type is genuinely unknown.
2. Send the raw file bytes to the endpoint. Do not wrap them in JSON or form data.
3. Report the returned URL, object ID, byte count, and expiry time. Keep the delete token private unless the user explicitly asks how to retain it.

```sh
curl --fail-with-body \
  --data-binary @artifact.tar.zst \
  -H 'content-type: application/zstd' \
  https://objects.yshubham.com/v1/objects
```

The successful `201` response includes `id`, `url`, `bytes`, `expiresAt`, and `deleteToken`.

## Retrieve or inspect an object

Use the exact `url` returned from upload:

```sh
curl --fail --remote-name 'https://objects.yshubham.com/v1/objects/OBJECT_ID'
curl --fail https://objects.yshubham.com/health
curl --fail https://objects.yshubham.com/api
```

Do not guess object IDs or scan the service. There is no listing endpoint by design.

## Delete early

Delete only an object you created or for which the user supplied the delete token. Pass the token in an authorization header, not on the command line if that command could be logged.

```sh
curl --fail -X DELETE \
  -H "authorization: Bearer $OBJECT_DELETE_TOKEN" \
  'https://objects.yshubham.com/v1/objects/OBJECT_ID'
```

If a token must survive the current task, store it only in a user-approved private secret store. State that it was stored without printing it.

## Handle errors deliberately

- `400`: the request is malformed or the body is empty. Correct the request; do not retry unchanged.
- `401` / `403`: the delete token is missing or invalid. Do not guess a replacement.
- `404`: the object is unknown or has expired. Re-upload only if the user still needs a public copy.
- `413`: the artifact exceeds 10 MiB. Compress, split, or use an appropriate user-approved storage system.
- `429`: stop and retry later with backoff. Do not run a retry loop.
- `503` or another `5xx`: storage is unavailable. Stop, report the failure, and avoid repeated uploads.

## Report the result

For a successful upload, give the user the public URL, expiry time, and byte count. Mention that it is public and temporary. Do not print `deleteToken` by default.
