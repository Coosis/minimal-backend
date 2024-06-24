# What is this?
A backend boilerplate, uses mongodb at localhost:27017 by default.

# How to use?
1. Clone the repository.
2. Build and run or simply do `go run .`.

# Endpoints?
1. `POST /auth/add` - Add a user
```bash
curl -X POST http://localhost:8080/auth/add \
    -d 'username=test' \
    -d 'pswdhash=testhash'
```

2. `POST /auth/del` - Delete a user
```bash
curl -X DELETE http://localhost:8080/auth/del \
    -H 'Authorization: bearer sometoken' \
    -d 'username=test' \
```

3. `POST /auth/login` - Login, returns a JWT token with a 24 hour expiry
```bash
curl -X POST http://localhost:8080/auth/login \
    -d 'username=test' \
    -d 'pswdhash=testhash'
```

4. `POST /auth/addtogroup` - Add a user to a group
```bash
curl -X POST http://localhost:8080/auth/addtogroup \
    -H 'Authorization: bearer sometoken' \
    -d 'groupname=testgroup' \
    -d 'username=test'
```

5. `POST /auth/rmfromgroup` - Remove a user from a group
```bash
curl -X DELETE http://localhost:8080/auth/rmfromgroup \
    -H 'Authorization: bearer sometoken' \
    -d 'groupname=testgroup' \
    -d 'username=test'
```

6. `POST /auth/addrighttogroup` - Add a right to a group
```bash
curl -X POST http://localhost:8080/auth/addrighttogroup \
    -H 'Authorization: bearer sometoken' \
    -d 'groupname=testgroup' \
    -d 'right=testright'
```

7. `POST /auth/rmrightfromgroup` - Remove a right from a group
```bash
curl -X DELETE http://localhost:8080/auth/rmrightfromgroup \
    -H 'Authorization: bearer sometoken' \
    -d 'groupname=testgroup' \
    -d 'right=testright'
```

8. `POST /auth/createadmin` - Create an admin user(only accessible through localhost)
```bash
curl -X POST http://localhost:8080/auth/createadmin \
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
