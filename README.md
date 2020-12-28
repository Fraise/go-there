# Go-there

[![Go Report Card](https://goreportcard.com/badge/github.com/Fraise/go-there)](https://goreportcard.com/report/github.com/Fraise/go-there)

Go-there is a simple, yet configurable, URL shortener. The end goal is to create a simple and coherent API that can be
used easily by any front-end application.

Go-there aims at being decently fast, be able to use the commonly used database and cache, as well as having proper
logging and monitoring out of the box.

It is currently a work in progress, and the API should not be considered stable before version 1.0.

## Authentication

The authentication is handled by the `auth` middleware. It can be enabled per service as described in the 
[Configuration](#Configuration) section. A user can either use a username/password combination or an API key.

The username/password can be provided either in JSON form in the request body:

```http request
{
  "username": "alice",
  "password": "secretpassword"
}
```

The API key can be provided either in JSON form in the request body:

```http request
{
  "api_key": "bi44RkM4YWwueFE0d2RvTkF5akpJTzpPSC1rbkdMcm91VlA3N01pZkJ1Y0F3PT0="
}
```

or in a X-Api-Key header:

```http request
X-Api-Key: bi44RkM4YWwueFE0d2RvTkF5akpJTzpPSC1rbkdMcm91VlA3N01pZkJ1Y0F3PT0=
```

If both authentication methods are used at the same time, only the API key will be checked.

## Configuration

The configuration uses the toml format.

### [Server]

`Mode` Set the server to debug or release mode. Set to "debug" for extra gin logging or "release" for production

`ListenAddress` The ip address the application listen to, formatted as "0.0.0.0"

`ListenPort` Port used by the application

### [Endpoints]

All endpoints can be configured using the array of values :

```toml
endpoint={ Enabled=true, Auth=false, AdminOnly=false, Log=true }
```

The available configuration groups are :

`health` represents the */health* endpoint

`create_users` represents the user creation method and endpoint: `POST` on */api/users*

`manage_users` represents the user management endpoint: `GET`, `DELETE` and `PATCH` on */api/:user*

`go` represents the redirection endpoint: `GET` on */go/:path*

`manage_path` represents the path management endpoint: `POST` and `DELETE` on */api/:path*

### [Cache]

The cache supports both Redis and local cache. It is only used to cache redirection requests, and local and network
caching can be enabled at the same time. It currently only supports a single Redis instance.

`Enabled` Enable the Redis cache

`Type` Unused

`Address` Network address of the Redis instance

`Port` Port used by Redis

`User` Username used to connect to the Redis instance

`Password` Password used to connect to the Redis instance

`LocalCacheEnabled` Enable the local cache

`LocalCacheSize` Size of the cache (in number of path/target pair)

`LocalCacheTtlSec` Lifetime in seconds of the elements in the local cache

### [Database]

Currently, the supported database type is mysql. In the future, postgresql and sqlite should also be supported.

`Type` The database type: "mysql"

`Address` A string representing the address of the database. Can be a domain or IP

`Port` The port to connect to

`SslMode` Should SSL be used for the connection: *true* or *false*

`Protocol` The connection protocol to use

`Name` The name of the database

`User` The user to identify as for the connection

`Password` The password of the connection user

### [Logs]

Base logging is enabled for the base operations (initialization...) but request logging should be enabled on an endpoint
basis as described in the *[Endpoints]* section. The available options common to every endpoint are :

`File` The file where the logs will be appended. The `$stdout` and `$stderr` string will respectively output the logs in
the OS' stdout or stderr. If left empty, it will output to stdout

`AsJSON` Format the logs as a JSON string. If it is set to false, the logs will be formatted and colored for the 
console, so they will be difficult to parse in a file