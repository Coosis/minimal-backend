# What is this?
A backend boilerplate, uses mongodb at localhost:27017 by default.

# How to use?
1. Clone the repository.
2. Build and run or simply do `go run .`.

# Endpoints?
1. `POST /auth/add` - Add a user
```bash
curl -X POST http://localhost:8080/auth/add \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d 'username=test' \
    -d 'pswdhash=testhash'
```
Response:
```json
{
    "message": "User {username} added"
}
```

2. `POST /auth/del` - Delete a user
```bash
curl -X POST http://localhost:8080/auth/del \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -H 'Authorization: bearer sometoken' \
    -d 'username=test' \
```
Response:
```json
{
    "message": "User {username} deleted"
}
```

3. `POST /auth/login` - Login, returns a JWT token with a 24 hour expiry
```bash
curl -X POST http://localhost:8080/auth/login \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d 'username=test' \
    -d 'pswdhash=testhash'
```
Response:
```json
{
    "username": {username},
    "token": {token}
}
```

4. `POST /auth/addtogroup` - Add a user to a group
```bash
curl -X POST http://localhost:8080/auth/addtogroup \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -H 'Authorization: bearer sometoken' \
    -d 'groupname=testgroup' \
    -d 'username=test'
```
Response:
```json
{
    "message": "User {username} added to group {groupname}"
}
```

5. `POST /auth/rmfromgroup` - Remove a user from a group
```bash
curl -X POST http://localhost:8080/auth/rmfromgroup \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -H 'Authorization: bearer sometoken' \
    -d 'groupname=testgroup' \
    -d 'username=test'
```
Response:
```json
{
    "message": "User {username} removed from group {groupname}"
}
```

6. `POST /auth/addrighttogroup` - Add a right to a group
```bash
curl -X POST http://localhost:8080/auth/addrighttogroup \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -H 'Authorization: bearer sometoken' \
    -d 'groupname=testgroup' \
    -d 'right=testright'
```
Response:
```json
{
    "message": "Right {right} added to group {groupname}"
}
```

7. `POST /auth/rmrightfromgroup` - Remove a right from a group
```bash
curl -X POST http://localhost:8080/auth/rmrightfromgroup \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -H 'Authorization: bearer sometoken' \
    -d 'groupname=testgroup' \
    -d 'right=testright'
```
Response:
```json
{
    "message": "Right {right} removed from group {groupname}"
}
```

8. `POST /auth/createadmin` - Create an admin user(only accessible through localhost)
```bash
curl -X POST http://localhost:8080/auth/createadmin \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d 'username=admin' \
    -d 'pswdhash=adminhash'
```

# Group Rights:
* `admin` - Can do everything
* `delete_user` - Can delete users
* `edit_group` - Can edit groups and its users and permissions

# Database structure
```markdown
db
|-users
| |-User
|   |-_id
|   |-username
|   |-pswdhash
|   |-usergroup
|-groups
| |-UserGroup
|   |-_id
|   |-groupname
|   |-users
|   |-permissions
```

# Logging
`logger` module is used for logging, log entries are of the spec `date time file:line message`, all logs are written to `logs.log` file, as well as stdout.
