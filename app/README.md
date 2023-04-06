# DevOps Academy Demo App

devops academy demo app

<img src="/docs/exports/context.png" />
<img src="/docs/exports/overall-container.png" />

# Prerequisites

- Docker
- Go 1.20+
- [Wire](https://github.com/google/wire) (for DI)
- [Taskfile](https://taskfile.dev/#/installation)

# Installation

# Run infrastructure

```bash
$ docker-compose --profile backend up -d
```

# Run services

## Copy .env to project root folder

```bash
$ cp env/local.env .env
```

## Run service

```bash
$ docker-compose up dev
```

## Create trip record

```bash
$ http --json -v post localhost:8090/v1/trips/ tripId=myTrip

POST /v1/trips/ HTTP/1.1
Accept: application/json, */*;q=0.5
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 20
Content-Type: application/json
Host: localhost:8090
User-Agent: HTTPie/3.2.1

{
    "tripId": "myTrip"
}


HTTP/1.1 200 OK
Content-Length: 116
Content-Type: application/json; charset=utf-8
Date: Thu, 23 Mar 2023 12:43:06 GMT

{
    "data": {
        "createdAt": "2023-03-23T21:43:06+09:00",
        "id": "myTrip",
        "status": "Initialized",
        "updatedAt": ""
    },
    "status": true
}
```

## Query created trips

```bash
$ http get localhost:8090/v1/trips/

HTTP/1.1 200 OK
Content-Length: 177
Content-Type: application/json; charset=utf-8
Date: Sun, 28 Aug 2022 12:38:01 GMT

{
    "data": [
        {
            "createdAt": "2022-08-28T19:36:47+07:00",
            "id": 1,
            "status": "INITIALIZED",
            "updatedAt": "0001-01-01T00:00:00Z",
        },
        {
            "createdAt": "2022-08-28T19:38:52+07:00",
            "id": 2,
            "status": "INITIALIZED",
            "updatedAt": "0001-01-01T00:00:00Z",
        },
        {
            "createdAt": "2022-08-28T19:38:53+07:00",
            "id": 3,
            "status": "INITIALIZED",
            "updatedAt": "0001-01-01T00:00:00Z",
        }
    ],
    "status": true
}
```
