# Todo API Demo

This is a tiny Go HTTP project for learning:

- basic syntax
- struct, slice, and methods
- error handling
- go modules
- net/http

## Run

```bash
go run .
```

## Routes

### Health check

```bash
curl http://localhost:8080/health
```

### List todos

```bash
curl http://localhost:8080/todos
```

### Get one todo

```bash
curl http://localhost:8080/todos/1
```

### Create todo

```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"text":"Read about go-zero next"}'
```

### Update todo

```bash
curl -X PATCH http://localhost:8080/todos/2 \
  -H "Content-Type: application/json" \
  -d '{"text":"Practice Go handlers","done":true}'
```

### Delete todo

```bash
curl -X DELETE http://localhost:8080/todos/2
```
