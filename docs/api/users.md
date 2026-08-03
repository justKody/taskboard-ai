# Auth / Users

Base path: `/api/v1`

```text
auth
└── users
    ├── PUT    Update User
    ├── PUT    Change Password
    ├── GET    Get Me
    ├── DEL    Delete User
    ├── GET    Get Users List
    ├── POST   Signup
    ├── POST   Login
    └── POST   Logout
```

Protected `/users/*` routes require JWT via `Authorization: Bearer <token>` or the `token` cookie set on login.

---

## POST /api/v1/signup

Auth: none

```bash
curl -X POST http://localhost:8080/api/v1/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com","password":"password123"}'
```

---

## POST /api/v1/login

Auth: none

Sets a `token` cookie on success.

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"email":"jane@example.com","password":"password123"}'
```

---

## POST /api/v1/logout

Auth: none

Clears the `token` cookie.

```bash
curl -X POST http://localhost:8080/api/v1/logout \
  -b cookies.txt
```

---

## GET /api/v1/users/me

Auth: required (JWT)

```bash
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer <token>"
```

---

## PUT /api/v1/users/update

Auth: required (JWT)

```bash
curl -X PUT http://localhost:8080/api/v1/users/update \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"Jane Smith"}'
```

---

## PUT /api/v1/users/change/password

Auth: required (JWT)

```bash
curl -X PUT http://localhost:8080/api/v1/users/change/password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"old_password":"password123","new_password":"newpassword123"}'
```

---

## DELETE /api/v1/users/delete

Auth: required (JWT)

```bash
curl -X DELETE http://localhost:8080/api/v1/users/delete \
  -H "Authorization: Bearer <token>"
```

---

## GET /api/v1/users/

Auth: required (JWT)

```bash
curl http://localhost:8080/api/v1/users/ \
  -H "Authorization: Bearer <token>"
```
