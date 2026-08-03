# Organizations

Base path: `/api/v1/organization`

All routes require JWT via `Authorization: Bearer <token>` or the `token` cookie.

```text
organization
├── POST   Create Organization
├── GET    Get Organization Details
├── GET    List My Organizations
├── PUT    Change Owner
└── DEL    Delete Organization
```

---

## Create Organization

### POST /api/v1/organization/create

Auth: required (JWT)

```bash
curl -X POST http://localhost:8080/api/v1/organization/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"Acme Corp"}'
```

---

## Get Organization Details

### GET /api/v1/organization/get/{id}

Auth: required (JWT)

```bash
curl http://localhost:8080/api/v1/organization/get/11111111-1111-1111-1111-111111111111 \
  -H "Authorization: Bearer <token>"
```

---

## List My Organizations

### GET /api/v1/organization/list

Auth: required (JWT)

Returns organizations for the authenticated user.

```bash
curl http://localhost:8080/api/v1/organization/list \
  -H "Authorization: Bearer <token>"
```

---

## Change Owner

### PUT /api/v1/organization/change-owner/{id}

Auth: required (JWT) — caller must be the current owner

```bash
curl -X PUT http://localhost:8080/api/v1/organization/change-owner/11111111-1111-1111-1111-111111111111 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"new_owner_id":"22222222-2222-2222-2222-222222222222"}'
```

---

## Delete Organization

### DELETE /api/v1/organization/delete/{id}

Auth: required (JWT) — caller must be the owner; org must have no members

```bash
curl -X DELETE http://localhost:8080/api/v1/organization/delete/11111111-1111-1111-1111-111111111111 \
  -H "Authorization: Bearer <token>"
```
