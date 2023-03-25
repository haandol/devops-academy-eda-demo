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

## Configuration

open [**config/dev.toml**](/infra/config/dev.toml) and fill the blow fields

```toml
[aws]
account="" # e.g. 123456789012

[vpc]
id="" # e.g. vpc-xxx

[msk]
seeds="" # e.g. xxx:9094,yyy:9094
securityGroupId="" # e.g. sg-xxx
```

and copy `config/dev.toml` file to project root as `.toml`

```bash
$ cd infra
$ cp config/dev.toml .toml
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
