# Task Manager API Documentation

This API allows you to manage tasks: create, read, update, and delete them with register and login functionality.

---

## Base URL

http://localhost:8080

## Endpoints

### 1. Register User

POST /register

**Description:**
Register user

Request Body Example:

```json
{
  "username": "ruth",
  "password": "password"
}
```

Include role **admin** to create admin user

```json
{
  "username": "ruth",
  "password": "password",
  "role": "admin"
}
```

Response Body Example:

```json
{
  "id": "uu984gr..",
  "username": "ruth",
  "role": "user"
}
```

### 1. Login User

POST /login

**Description:**
Login

Request Body Example:

```json
{
  "username": "ruth",
  "Password": "password"
}
```

Response Body Example:

```json
{
  "message": "Login successfu",
  "token": "eyJhbGciOiJIUzI1NiI....",
  "id": "uu984gr..",
  "username": "ruth",
  "role": "user"
}
```

### 1. Get All Tasks

GET /tasks

```bash
  Authorization: Bearer <JWT_TOKEN>
  Content-Type: application/json
```

**Description:**  
Retrieve a list of all tasks.

**Response Example:**

```json
[
  { "id": "1", "title": "Finish Go tutorial" },
  { "id": "2", "title": "Write blog post" },
  { "id": "3", "title": "Read Gin documentation" }
]
```

Status Codes:

```bash
200 OK – Successfully retrieved tasks
```

### 2. Get Task by ID

GET /tasks/:id

```bash
  Authorization: Bearer <JWT_TOKEN>
  Content-Type: application/json
```

Description:
Retrieve a single task by its ID.

URL Parameters:

id – The ID of the task

Response Example (if found):

```json
{ "id": "1", "title": "Finish Go tutorial" }
```

Response Example (if not found):

```json
{ "message": "task not found" }
```

Status Codes:

```bash
200 OK – Task found

404 Not Found – Task does not exist
```

### 3. Create a New Task

POST /tasks

```bash
  Authorization: Bearer <JWT_TOKEN>
  Content-Type: application/json
```

Description:
Add a new task.

Request Body Example:

```json
{
  "title": "Learn Gin framework"
}
```

Response Example:

```json
{
  "id": "4",
  "title": "Learn Gin framework"
}
```

Status Codes:

```bash
201 Created – Task successfully created

400 Bad Request – Missing id or title
```

### 4. Update a Task

PUT /tasks/:id

```bash
  Authorization: Bearer <JWT_TOKEN>
  Content-Type: application/json
```

Description:
Update the title of an existing task. You only need to provide the title in the body.

URL Parameters:

id – The ID of the task to update

Request Body Example:

```json
{
  "title": "Learn Go basics"
}
```

Response Example:

```json
{
  "id": "4",
  "title": "Learn Go basics"
}
```

Response Example (if not found):

```json
{ "message": "task not found" }
```

Status Codes:

```bash
200 OK – Task successfully updated

400 Bad Request – Missing title in body

404 Not Found – Task does not exist
```

### 5. Delete a Task

DELETE /tasks/:id

```bash
  Authorization: Bearer <JWT_TOKEN>
  Content-Type: application/json
```

Description:
Delete a task by its ID.

URL Parameters:

id – The ID of the task to delete

Response Example (if deleted):

```json
{ "message": "task deleted" }
```

Response Example (if not found):

```json
{ "message": "task not found" }
```

Status Codes:

```bash
200 OK – Task successfully deleted

404 Not Found – Task does not exist
```
