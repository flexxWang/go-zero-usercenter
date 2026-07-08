# gozero-user-demo

This demo is a more production-shaped `go-zero` service than the beginner sample in the repo root.

## What it shows

- `goctl` generated `handler / logic / svc / types`
- MySQL-backed `users` table
- Redis-backed login session verification
- JWT protected routes
- unified success/error response envelopes
- request validation through `Validate()` on generated request types
- cache + repository split instead of all logic living in handlers

## Endpoints

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/users/me`
- `GET /api/v1/users/:userId`
- `PATCH /api/v1/users/me`

## Quick start

1. Start infra:

```bash
cd demos/gozero-user-demo/deploy
docker compose up -d
```

2. Run the service:

```bash
cd /Users/chengjiaxiang/Desktop/zwo/workspace/go
go run . -f etc/usercenter-api.yaml
```

## Example requests

Register:

```bash
curl -X POST http://localhost:8888/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"tom@example.com","password":"123456","nickname":"Tom"}'
```

Login:

```bash
curl -X POST http://localhost:8888/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"tom@example.com","password":"123456"}'
```

Get current user:

```bash
curl http://localhost:8888/api/v1/users/me \
  -H "Authorization: Bearer <accessToken>"
```

Update profile:

```bash
curl -X PATCH http://localhost:8888/api/v1/users/me \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <accessToken>" \
  -d '{"nickname":"Tommy","bio":"frontend -> go"}'
```

## Reading order

If you want to understand the flow quickly, read files in this order:

1. `usercenter.api`
2. `internal/handler/routes.go`
3. `internal/types/types.go` and `internal/types/validation.go`
4. `internal/logic/auth/loginlogic.go`
5. `internal/logic/user/updateprofilelogic.go`
6. `internal/model/usermodel.go`
7. `internal/svc/servicecontext.go`
8. `internal/responsex/response.go`
