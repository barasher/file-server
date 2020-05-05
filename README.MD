# file-server

**`file-server`** is an ultra-simple webservice that exposes files over *http*.

It exposes three services :

- a service that returns stored binaries (~ GET)
- a **set** service that stores binaries (~ SET)
- a **metrics** service that exposes... metrics !

The unique key for each binary is its filename.

**`file-server`** can store (and expose) files :
- on a regular filesystem
- on a S3 bucket
 
## Getting file-server

**`file-server`** is released :
- As a compiled binary on the [Github project release section](https://github.com/barasher/file-server/releases). If you have a [golang environment](https://golang.org/doc/install), you can compile the projet : `go build -o file-server main.go`
- As a docker imaged on [Docker Hub](https://hub.docker.com/r/barasher/file-server/tags). If you have a [docker environment], you can build the image : `docker build .` 

## How to configure

In the **`file-server`** docker image :
- the configuration file is located in `/etc/file-server/file-server.json`
- by default, files are stored on the filesystem in the folder `/data/file-server/`
- by default, the server is listening on port `8080`

The configuration file is a JSON file :
```json
{
  "LoggingLevel" : "debug",
  "Type": "configurationType",
  "Port": 8080,
  [Storage speficic configuration]
}
```

- `LoggingLevel` defines the logging level : `debug`, `info` (default value), `warn`, `error`, ...
- `Type` defines which type of storage that must be used. This parameter is **mandatory** and depending on the value, other parameters become also required. Acceptable values are :
  - `s3` to store data in a S3 bucket.
  - `local` to store data in the local filesystem
- `Port` defines the port where the server will be listening. The default value is `8080`

### Local storage

If the chosen storage type is `local`, then the following block will be required :
```json
{
  [...],
  "Local": {
    "Folder": "/var/data/file-server"
  }
}
```

- `Folder` defines where files will be stored. This parameter is **mandatory**.

### S3 storage

If the chosen storage type is `s3`, then the following block will be required :
```json
{
  [...],
  "S3": {
    "AccessId": "myId",
    "AccessSecret": "mySecret",
    "Bucket": "myBucket",
    "URL": "http://1.2.3.4:8080/"
  }
}
```

- `AccessId` defines the S3 access id. This parameter is **mandatory**.
- `AccessSecret` defines the secret for the id. This parameter is **mandatory**.
- `Bucket` defines the bucket name. This parameter is **mandatory**.
- `URL` defines the S3 entrypoint. This parameter is **mandatory**.

## Exposed services

### Set a binary

Request :
- Path : `/key/{keyValue}`
- Method : `POST`
- Content-type : `multipart/form-data`
- Body : multipart, part key : `file` 

Response :
- Main HTTP status codes :
  - 204 : OK
  - 500 : Server error

### Get a binary

Request :
- Path : `/key/{keyValue}`
- Method : `GET`

Response :
- Main HTTP status codes :
  - 200 : OK
  - 404 : Unknown key
  - 500 : Server error
- Body : file content

### Metrics

Path : `/metrics`

Metrics are based on (https://prometheus.io/)[Prometheus] standards :
- `get` operations are monitored as histogram : `file_server_get_request_duration_seconds_bucket`
- `set` operations are monitored as histogram :  `file_server_set_request_duration_seconds_bucket`

## Changelog