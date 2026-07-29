# Backend Challenge: Project Management API

## Objective

Build a REST API using **Go** (or any backend language) with **PostgreSQL** or **MySQL**.

The goal is to practice:

- CRUD operations
- Database relationships
- Foreign keys
- Transactions
- Pagination
- Filtering
- SQL joins
- Proper database design

---

# Database Schema

## 1. Users

A user can own organizations, belong to organizations, create projects, create tasks, and comment on tasks.

| Column | Type |
|---------|------|
| id | UUID / INT |
| name | VARCHAR |
| email | VARCHAR (Unique) |
| created_at | TIMESTAMP |

---

## 2. Organizations

An organization is owned by a user.

| Column | Type |
|---------|------|
| id | UUID / INT |
| name | VARCHAR |
| owner_id | FK → users.id |
| created_at | TIMESTAMP |

Relationship

```
User
 └── Organization
```

---

## 3. Memberships (Many-to-Many)

A user may belong to multiple organizations.

An organization has multiple members.

| Column | Type |
|---------|------|
| organization_id | FK |
| user_id | FK |
| role | admin/member |
| joined_at | TIMESTAMP |

Composite Primary Key

```
organization_id
user_id
```

Relationship

```
Users
   ▲
   │
Membership
   │
   ▼
Organizations
```

---

## 4. Projects

Each project belongs to one organization.

| Column | Type |
|---------|------|
| id | UUID |
| organization_id | FK |
| name | VARCHAR |
| description | TEXT |
| status | active/completed |
| created_by | FK → users.id |
| created_at | TIMESTAMP |

Relationship

```
Organization
    └── Projects
```

---

## 5. Tasks

Each project contains many tasks.

| Column | Type |
|---------|------|
| id | UUID |
| project_id | FK |
| title | VARCHAR |
| description | TEXT |
| status | todo/in_progress/done |
| priority | low/medium/high |
| due_date | DATE |
| assigned_to | FK → users.id |
| created_by | FK → users.id |
| created_at | TIMESTAMP |

Relationship

```
Project
    └── Tasks
```

---

## 6. Comments

Each task can have many comments.

| Column | Type |
|---------|------|
| id | UUID |
| task_id | FK |
| user_id | FK |
| message | TEXT |
| created_at | TIMESTAMP |

Relationship

```
Task
    └── Comments
```

---

## 7. Labels

Labels belong to an organization.

| Column | Type |
|---------|------|
| id | UUID |
| organization_id | FK |
| name | VARCHAR |
| color | VARCHAR |

---

## 8. Task Labels (Many-to-Many)

A task may have multiple labels.

A label may belong to multiple tasks.

| Column | Type |
|---------|------|
| task_id | FK |
| label_id | FK |

Composite Primary Key

```
task_id
label_id
```

---

# Entity Relationship Diagram

```
User
│
├── owns Organizations
│
├── Memberships
│      │
│      ▼
│  Organization
│       │
│       ├── Projects
│       │      │
│       │      └── Tasks
│       │              │
│       │              ├── Comments
│       │              └── TaskLabels
│       │                      │
│       │                      ▼
│       │                    Labels
```

---

# Required API Endpoints

## Users

- Create User
- Get User
- Update User
- Delete User
- List Users

---

## Organizations

- Create Organization
- Get Organization
- Update Organization
- Delete Organization
- List Organizations

---

## Memberships

- Add Member
- Remove Member
- List Members

---

## Projects

- Create Project
- Get Project
- Update Project
- Delete Project
- List Projects

---

## Tasks

- Create Task
- Get Task
- Update Task
- Delete Task
- List Tasks

---

## Comments

- Create Comment
- Get Comments
- Delete Comment

---

## Labels

- Create Label
- Update Label
- Delete Label
- Assign Label to Task
- Remove Label from Task

---

# Filtering

Implement support for:

- Task status
- Priority
- Assigned user
- Due date
- Organization
- Project

Example

```
GET /tasks?status=todo
GET /tasks?priority=high
GET /tasks?assigned_to=12
GET /tasks?project_id=5
```

---

# Pagination

Support

```
?page=1
&limit=20
```

---

# Sorting

Support

```
?sort=created_at

?sort=priority

?sort=due_date
```

Ascending and descending order.

---

# Search

Support keyword searching.

Example

```
GET /tasks?search=database

GET /projects?search=mobile
```

---

# Transactions

Implement transactions where appropriate.

Example:

Creating a project

- Create project
- Create default labels
- Create default project settings

If one step fails, rollback everything.

---

# Validation Rules

- Email must be unique.
- Organization names should be unique per owner.
- Project must belong to an existing organization.
- Task must belong to an existing project.
- Assigned user must be a member of the organization.
- Labels cannot be attached across organizations.
- Prevent duplicate memberships.
- Prevent duplicate task labels.

---

# Soft Delete (Bonus)

Instead of deleting records permanently:

```
deleted_at TIMESTAMP NULL
```

Filter deleted records automatically.

---

# Authentication (Bonus)

Use JWT authentication.

Protected endpoints should require a valid token.

---

# Authorization (Bonus)

Implement role-based permissions.

Example:

Admin can:

- Invite members
- Delete projects
- Delete labels

Member can:

- Create tasks
- Update assigned tasks
- Comment

---

# Advanced SQL Challenges

Try implementing the following queries:

- Number of tasks per project
- Number of completed tasks
- Number of tasks assigned to each user
- Most active member
- Projects with overdue tasks
- Top 5 busiest users
- Organization dashboard statistics

---

# Bonus Features

- File attachments
- Activity logs
- Task history
- Notifications
- Project archive
- Favorite projects
- Recently viewed tasks
- Audit logs

---

# Expected Learning Outcomes

By completing this project you should become comfortable with:

- CRUD APIs
- RESTful design
- SQL joins
- One-to-Many relationships
- Many-to-Many relationships
- Foreign keys
- Composite keys
- Transactions
- Pagination
- Filtering
- Sorting
- Search
- Database indexing
- Authorization
- Database normalization
- Clean backend architecture
```

# Learn Make files and all comands should be written down in Makefile
