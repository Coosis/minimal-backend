# What is this?
A backend boilerplate, uses mongodb at localhost:27017 by default.

# How to use?
1. Clone the repository.
2. Build and run or simply do `go run .`.

# Endpoints?
1. `POST /auth/add` - Add a user
```bash
curl -X POST http://localhost:8080/auth/add -d 'username=test' -d 'pswdhash=testhash'
```
2. `POST /auth/login` - Login, returns a JWT token with a 24 hour expiry
```bash
curl -X POST http://localhost:8080/auth/login -d 'username=test' -d 'pswdhash=testhash'
```
3. `POST /auth/addtogroup` - Add a user to a group
```bash
curl -X POST http://localhost:8080/auth/addtogroup -H 'Authorization: bearer admintoken' -d 'groupname=testgroup' -d 'username=test'
```

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
