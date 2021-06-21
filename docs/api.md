## Add app

```
POST /app/
{
  "app": "myapp",
  "desc": "my first app for test"
}
=>
{
  "error": 0,
  "message": "",
  "data": {
    "key": "xxx"
  }
}
```

## Get app info

```
GET /app/
{
  "app": "myapp",
  "desc": "",
  "roles": [{
    "role": "role1",
    "desc": "",
    "permissions": [
      "p1",
      "p2"
    ]
  }, {
    "role": "role2",
    "desc": "",
    "permissions": [
      "p2",
      "p3"
    ]
  }],
  "permissions": [{
    "permission": "p1",
    "desc": ""
  }, {
    "permission": "p2",
    "desc": ""
  }, {
    "permission": "p3",
    "desc": ""
  }, {
    "permission": "p4",
    "desc": ""
  }]
}
```

## Add role to app

```
POST /app/role/
{
  "app": "myapp",
  "role": "role1",
  "desc": "description of myapp's role1"
}
```

## Remove role from app

```
DELETE /app/role/
{
  "app": "myapp",
  "role": "role1"
}
```

## Add permission to app

```
POST /app/permission/
{
  "app": "myapp",
  "permission": "p1",
  "desc": "description of permission p1"
}
```

## Remove permission from app

```
DELETE /app/permission/
{
  "app": "myapp",
  "permission": "p1"
}
```

## Bind permission to role

```
POST /app/role/permission/
{
  "app": "myapp",
  "role": "role1",
  "permission": "p1"
}
```

## Bind permission to role

```
POST /app/role/permission/
{
  "app": "myapp",
  "role": "role1",
  "permissions": [
    "p1",
    "p2"
  ]
}
```

## Unbind permission to role

```
DELETE /app/role/permission/
{
  "app": "myapp",
  "role": "role1",
  "permissions": [
    "p1",
    "p2"
  ]
}
```

## Get users

```
GET /user/?page=1&size=20
==>
{
  "error": 0,
  "message": "",
  "data": {
    "total": 123,
    "page": 1,
    "size": 20,
    "data": [{
      "id": 123,
      "phone": "13012345678",
      "name": "xxx",
      "email": "xxx",
      "extra": {
        // weixin, weibo, ...
        "weixin-session": "xxx"
      },
      "apps": [{
        "app": "myapp",
        "expired_at": "2031-10-12 00:00:00",
        "roles": [
          "role1",
          "role2"
        ]
      }]
    }]
  }
}
```

## add user

```
POST /user/
{
  "phone": "13012345678",
  "name": "xxx",
  "email": "xxx",
  "extra": {
    // weixin, weibo, ...
  }
}
```

## set user's app roles

```
PUT /user/app/role/
{
  "user": 123,
  "app": "myapp",
  "roles": [
    "role1",
    "role2"
  ]
}
```

## set user's app expired datetime

```
PUT /user/app/expired_at/
{
  "user": 123,
  "app": "myapp",
  "expired_at": "2031-10-12 00:00:00"
}
```

## sso login

```
POST /sso/login/
{
  "username": "xxx",
  "password": "xxx"
}
==>
{
  "error": 0,
  "message": "",
  "data": {
    "sso_token": "xxx"
  }
}
```

## sso backend login
```
POST /sso/login/
{
  "username": "xxx",
  "password": "xxx"
}
==>
{
  "error": 0,
  "message": "",
  "data": {
    "token": "xxx"
  }
}
```

## sso manager pages

```
GET /manager/
```
