# Go Workspace

This workspace is organized around demos.
Each main demo now has its own `go.mod`, and the root uses `go.work` to keep them in one workspace.

- `demos/basic-go-demo`
  - a plain Go `net/http` todo demo for syntax and standard library learning
- `demos/gozero-user-demo`
  - a more production-shaped `go-zero` demo with MySQL, Redis, JWT, middleware, and unified responses

## Run

Basic Go demo:

```bash
cd demos/basic-go-demo
go run .
```

go-zero demo:

```bash
cd demos/gozero-user-demo
go run . -f etc/usercenter-api.yaml
```
