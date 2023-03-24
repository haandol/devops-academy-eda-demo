#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { VpcStack } from '../lib/stacks/vpc-stack';
import { BastionHostStack } from '../lib/stacks/bastion-host-stack';
import { DatabaseStack } from '../lib/stacks/database-stack';
import { EcsClusterStack } from '../lib/stacks/ecs-cluster-stack';
import { TripServiceStack } from '../lib/stacks/services/trip-service-stack';
import { Config } from '../config/loader';

const app = new cdk.App({
  context: {
    ns: Config.app.ns,
    stage: Config.app.stage,
  },
});

new DatabaseStack(app, `${Config.app.ns}DatabaseStack`, {
  tableName: Config.ddb.tableName,
});

const vpcStack = new VpcStack(app, `${Config.app.ns}VpcStack`, {
  vpcId: Config.vpc.id,
  env: {
    account: Config.aws.account,
    region: Config.aws.region,
  },
});

new BastionHostStack(app, `${Config.app.ns}BastionHostStack`, {
  vpcId: vpcStack.vpc.vpcId,
  env: {
    account: Config.aws.account,
    region: Config.aws.region,
  },
});

const ecsClusterStack = new EcsClusterStack(
  app,
  `${Config.app.ns}EcsClusterStack`,
  {
    vpc: vpcStack.vpc,
    env: {
      account: Config.aws.account,
      region: Config.aws.region,
    },
  }
);
ecsClusterStack.addDependency(vpcStack);

const tripServiceStack = new TripServiceStack(
  app,
  `${Config.app.ns}TripServiceStack`,
  {
    vpc: vpcStack.vpc,
    alb: ecsClusterStack.alb,
    cluster: ecsClusterStack.cluster,
    taskRole: ecsClusterStack.taskRole,
    taskLogGroup: ecsClusterStack.taskLogGroup,
    taskExecutionRole: ecsClusterStack.taskExecutionRole,
    taskSecurityGroup: ecsClusterStack.taskSecurityGroup,
    service: {
      name: Config.service.trip.name,
      repositoryName: Config.service.trip.repositoryName,
      port: Config.service.common.port,
      tag: Config.service.common.tag
    },
    env: {
      account: Config.aws.account,
      region: Config.aws.region,
    },
  }
);
tripServiceStack.addDependency(ecsClusterStack);

const tags = cdk.Tags.of(app);
tags.add('namespace', Config.app.ns);
tags.add('stage', Config.app.stage);

app.synth();
