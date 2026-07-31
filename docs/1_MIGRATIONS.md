# Migration Architecture

Schema changes in this project are managed with [golang-migrate](https://github.com/golang-migrate/migrate) (`v4`), PostgreSQL via `pgx/v5`, and Make targets that wrap a small Go CLI.

---

## Directory layout

```text
server/
├── Makefile                          # create / up / down / force targets
├── .env                              # DATABASE_URL (loaded by config)
├── config/env.go                     # reads DATABASE_URL for migrate + app
├── sqlc.yaml                         # uses migrations folder as schema source
├── db/
│   ├── schema.sql                    # auto-dumped schema (do not hand-edit)
│   ├── queries/                      # hand-written SQL for sqlc
│   └── sqlc/                         # generated Go from schema + queries
└── cmd/migrate/
    ├── main.go                       # migrate CLI entrypoint
    └── migrations/                   # versioned .up.sql / .down.sql pairs
        ├── 20260728132001_create_user_table.up.sql
        ├── 20260728132001_create_user_table.down.sql
        ├── 20260731003123_added_organization_schema.up.sql
        └── 20260731003123_added_organization_schema.down.sql
```

---

## Components

| Piece | Role |
| --- | --- |
| `github.com/golang-migrate/migrate/v4` | Core migrator: versions, up/down/force, dirty flag |
| `database/pgx/v5` driver | Applies SQL against Postgres through `database/sql` + pgx |
| `source/file` | Loads SQL from `file://cmd/migrate/migrations` |
| `cmd/migrate/main.go` | Thin CLI: `up`, `down`, `force` |
| `Makefile` | Developer ergonomics around the CLI, `migrate create`, and schema dump |
| `config.Envs.DatabaseUrl` | Connection string from `.env` (`DATABASE_URL`) |
| Postgres `schema_migrations` | Bookkeeping table (version + dirty) managed by golang-migrate |
| `db/schema.sql` | Auto-dumped full schema after every up/down — committed source of truth |
| `sqlc.yaml` → `schema: cmd/migrate/migrations` | Generated query code tracks the same SQL as migrations |

---

## How the runner works

```text
make migrate-up / migrate-down / migrate-force
        │
        ▼
go run cmd/migrate/main.go <command> [args]
        │
        ├─ sql.Open("pgx", DATABASE_URL)
        ├─ pgxmigrate.WithInstance(...)
        └─ migrate.NewWithDatabaseInstance(
               "file://cmd/migrate/migrations",
               "pgx",
               driver,
           )
                │
                ▼
         argsCommands(m)
           up    → m.Steps(1)      # apply next 1 migration
           down  → m.Steps(-1)     # rollback last 1 migration
           force → m.Force(version) # fix version bookkeeping only
                │
                ▼  (up / down only)
         make schema-dump
           pg_dump --schema-only → db/schema.sql
```

Important behaviors:

- **`up` / `down` are one step at a time** (`Steps(±1)`), not “all pending / all the way to zero”.
- Source path is relative to the process working directory — run Make targets from `server/`.
- `ErrNoChange` is ignored (already at latest / already at zero).
- **`force` does not run SQL.** It only updates `schema_migrations` (and clears dirty). The real schema must already match the version you force to.

---

## Makefile commands

Run from `server/`.

### Create a new migration pair

Requires the [`migrate` CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) installed locally.

```bash
make migration add_something_descriptive
```

Creates:

```text
cmd/migrate/migrations/<timestamp>_add_something_descriptive.up.sql
cmd/migrate/migrations/<timestamp>_add_something_descriptive.down.sql
```

### Apply / rollback / repair

```bash
make migrate-up                      # apply next migration, then dump schema
make migrate-down                    # roll back latest migration, then dump schema
make migrate-force 20260731003123    # set version (clears dirty); no SQL, no dump
make schema-dump                     # refresh db/schema.sql without migrating
```

`migrate-up` and `migrate-down` always finish with `schema-dump`. Commit the updated `db/schema.sql` with the migration.

---

## File naming convention

```text
{version}_{snake_case_name}.{up|down}.sql
```

| Part | Meaning |
| --- | --- |
| `version` | Timestamp / integer sort key (e.g. `20260728132001`) |
| `name` | Human label (`create_user_table`) |
| `.up.sql` | Forward change |
| `.down.sql` | Exact reverse of that up |

Migrations are applied in ascending version order; downs run in reverse.

---

## Schema source of truth

Do not maintain a hand-written “current schema” section in docs — it goes stale. After every successful `migrate-up` / `migrate-down`, the Makefile runs:

```bash
pg_dump --schema-only --no-owner --no-privileges $DATABASE_URL > db/schema.sql
```

Read `db/schema.sql` for the full current schema. Do not edit that file by hand; regenerate with `make schema-dump`.

### Dependency rules (especially for downs)

Foreign keys dictate drop order. Dependent objects must go first:

```sql
-- correct down for org schema
DROP TABLE IF EXISTS memberships;
DROP TABLE IF EXISTS organizations;
DROP TYPE IF EXISTS memberships_role;
```

Dropping `organizations` before `memberships`, or `users` while org tables still exist, fails with `SQLSTATE 2BP01`.

---

## Version tracking & dirty state

golang-migrate stores state in Postgres:

| Column | Meaning |
| --- | --- |
| `version` | Last successfully applied migration version |
| `dirty` | `true` if the last attempt failed mid-migration |

When dirty:

```text
Dirty database version <N>. Fix and force version.
```

Recovery pattern:

1. Inspect schema — what tables/types actually exist?
2. Fix the failing `.up.sql` / `.down.sql` if the SQL was wrong.
3. Manually align the DB with a known good version if needed.
4. `make migrate-force <version_that_matches_reality>`
5. Re-run `migrate-up` or `migrate-down`.

Never force a version that does not match the real schema.

---

## Relationship to sqlc

```text
cmd/migrate/migrations/*.up.sql  ──►  sqlc schema
db/queries/*.sql                 ──►  sqlc queries
                │
                ▼
         make sqlc  /  sqlc generate
                │
                ▼
            db/sqlc/   (generated Go)
```

Migrations are the source of truth for tables/types. After adding or changing migrations, regenerate sqlc so store code stays in sync.

---

## Config & prerequisites

| Requirement | Detail |
| --- | --- |
| `.env` with `DATABASE_URL` | Loaded by `config` via `godotenv`; `schema-dump` reads the same key from `.env` |
| Postgres reachable | Same DB the app uses |
| `pg_dump` ≥ server major | Needed for `make schema-dump` (Makefile prefers Homebrew `libpq`) |
| Working directory `server/` | So `file://cmd/migrate/migrations` resolves |
| Optional: `migrate` CLI | Only needed for `make migration <name>` |

---

## Mental model

```text
┌─────────────────┐     Steps(1)      ┌──────────────────┐
│  migration N-1  │ ───────────────►  │   migration N    │
│  (applied)      │  ◄───────────────  │   (applied)      │
└─────────────────┘     Steps(-1)     └──────────────────┘
         ▲                                      │
         │         schema_migrations            │
         └──────── version + dirty ◄────────────┘
```

- **Up SQL** moves the schema forward one version.
- **Down SQL** must fully reverse that version (order-sensitive with FKs).
- **Force** rewrites the bookmark when bookkeeping and reality diverge after a failed run.
