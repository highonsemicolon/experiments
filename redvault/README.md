## RedVault

#### API for notes with private, public, and group access

Create a public note

```bash
curl -X POST http://localhost:3000/notes \
-H "Content-Type: application/json" \
-H "user-id: user:111" \
-d '{"title":"Public Note","content":"Hello world","type":"public"}'
```


Create a group
```bash
curl -X POST http://localhost:3000/groups \
-H "Content-Type: application/json" \
-d '{"members":["user:111","user:222","user:333"]}'
```

Create a group note
```bash
curl -X POST http://localhost:3000/notes \
-H "Content-Type: application/json" \
-H "user-id: user:111" \
-d '{"title":"Team Plan","content":"Details","type":"group","groupId":"1"}'
```

List notes (that a user has access to)
```bash
curl -X GET http://localhost:3000/notes \
-H "user-id: user:222"
```

Get note by id
```bash
curl -X GET http://localhost:3000/notes/2 \
-H "user-id: user:222"
```