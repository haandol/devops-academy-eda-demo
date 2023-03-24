# Demo Infra

# Prerequisites

- git
- awscli
- Nodejs 16.x
- AWS Account and locally configured AWS credential

# Installation

## Install dependencies

```bash
$ cd infra
$ npm i -g aws-cdk
$ npm i
```

we are using [Taskfile](https://taskfile.dev/) for running script
install taskfile cli

```bash
$ npm i -g @go-task/cli
$ task --list-all
```

## Configuration

open [**config/dev.toml**](/infra/config/dev.toml) and fill the blow fields

and copy `config/dev.toml` file to project root as `.toml`

```bash
$ cd infra
$ cp config/dev.toml .toml
```

## Setup ECR repositories for services

### Create repositories

```bash
$ task create-repo
```

### Push initial images

```bash
$ task push-echo
```

### Create topics in MSK

- trip, 10 partitions, 3 replication factor, 2 min.insync.replicas
- car, 10 partitions, 3 replication factor, 2 min.insync.replicas
- hotel, 10 partitions, 3 replication factor, 2 min.insync.replicas
- flight, 10 partitions, 3 replication factor, 2 min.insync.replicas

```bash
$ task create-topic
```

## Deploy for dev

if you never run bootstrap on the account, bootstrap it.

```bash
$ cdk bootstrap
```

deploy infrastructure

```bash
$ cdk deploy "*" --require-approval never
```

## Troubleshooting

### Unable to assume the service linked role. Please verify that the ECS service linked role exists.

> there is no ECS service linked-role because you never been trying to create a ECS cluster before. trying to ECS cluster will create the service linked-role.

A: just re-deploy using `cdk deploy` after rollback the stack.
