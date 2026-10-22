# Realtime

Base path: `/api/v1`

```text
realtime
└── GET    WebSocket
```

The handshake is a GET. Auth is the same JWT as the REST API: `Authorization: Bearer <token>` or the `token` cookie.

---

## GET /api/v1/ws/

Auth: required (JWT)

Upgrades the HTTP connection to a WebSocket. The server pushes text frames to every connected client and keeps the socket alive with ping frames. Inbound client frames are read and discarded.

```bash
curl --include \
  -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Version: 13" \
  -H "Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==" \
  -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/ws/
```
