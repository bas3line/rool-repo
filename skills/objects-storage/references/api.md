# Objects API reference

Base URL: `https://objects.yshubham.com`

Objects is a public, temporary artifact hand-off service. Every object URL is publicly retrievable, has a maximum payload of 10 MiB, and expires seven days after creation. It has no account layer, object listing, or custom object keys.

## Endpoints

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/health` | Service and storage health |
| `GET` | `/api` | Machine-readable service contract |
| `POST` | `/v1/objects` | Upload raw bytes as one public object |
| `GET` | `/v1/objects/:id` | Download one public object |
| `DELETE` | `/v1/objects/:id` | Delete an object early with its delete token |

## Upload

Send a non-empty raw request body. Set `content-type` to the artifact MIME type. The API derives the object ID and public URL; callers cannot choose either.

```sh
curl --fail-with-body \
  --data-binary @report.json \
  -H 'content-type: application/json' \
  https://objects.yshubham.com/v1/objects
```

Successful response (`201`):

```json
{
  "id": "object-id",
  "url": "https://objects.yshubham.com/v1/objects/object-id",
  "bytes": 2048,
  "expiresAt": "2026-07-28T12:00:00.000Z",
  "deleteToken": "private-delete-capability"
}
```

`deleteToken` is intentionally shown here only as an API field. Treat real values as secrets: never publish, log, or include them in URLs.

## Download

Request the URL returned by the upload response. Do not probe IDs or assume a `404` distinguishes an expired object from an unknown one.

```sh
curl --fail --remote-name 'https://objects.yshubham.com/v1/objects/object-id'
```

## Early deletion

Use the upload response’s delete token as a bearer token. Successful deletion returns `204 No Content`.

```sh
curl --fail -X DELETE \
  -H "authorization: Bearer $OBJECT_DELETE_TOKEN" \
  'https://objects.yshubham.com/v1/objects/object-id'
```

## Errors and limits

| Status | Meaning | Safe response |
| --- | --- | --- |
| `400` | Empty or malformed upload | Correct input before retrying |
| `401` | Delete token missing | Request the token from its owner |
| `403` | Delete token invalid | Stop; do not guess tokens |
| `404` | Object unknown or expired | Re-upload only if still needed |
| `413` | Payload exceeds 10 MiB | Reduce size or use another approved store |
| `429` | Request limit reached | Back off; do not loop retries |
| `503` | Storage unavailable | Stop and report the service issue |
