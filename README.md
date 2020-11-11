# Go-there

[![Go Report Card](https://goreportcard.com/badge/github.com/Fraise/go-there)](https://goreportcard.com/report/github.com/Fraise/go-there)

Go-there is a simple, yet configurable, URL shortener.

## Authentication

The authentication is handled by the `auth` middleware. It can be enabled per service as described in the 
[Configuration](#Configuration) section. A user can either use a username/password combination or an API key.

The username/password can be provided either in URL form input:

```http request
http://example.com/users?username=alice&password=secretpassword
```

or in JSON form in the request body:

```http request
{
  "username": "alice",
  "password": "secretpassword"
}
```

The API key can be provided either in URL form input:

```http request
http://example.com/users?api_key=MEVEejdDZWR3TmFUaVNyaTRQdlJFdQ==.KaK8b3OwgOk6VW-MaXOwPA==
```

in JSON form in the request body:

```http request
{
  "api_key": "MEVEejdDZWR3TmFUaVNyaTRQdlJFdQ==.KaK8b3OwgOk6VW-MaXOwPA==",
}
```

or in a X-Api-Key header:

```http request
X-Api-Key: MEVEejdDZWR3TmFUaVNyaTRQdlJFdQ==.KaK8b3OwgOk6VW-MaXOwPA==
```

If both authentication methods are used at the same time, only the API key will be checked.

## Configuration

*TODO*