import { Construct } from 'constructs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as ecs from 'aws-cdk-lib/aws-ecs';
import * as ecr from 'aws-cdk-lib/aws-ecr';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as logs from 'aws-cdk-lib/aws-logs';

interface IProps {
  readonly cluster: ecs.ICluster;
  readonly taskRole: iam.IRole;
  readonly taskLogGroup: logs.ILogGroup;
  readonly taskExecutionRole: iam.IRole;
  readonly taskSecurityGroup: ec2.ISecurityGroup;
  readonly taskEnvs: { [key: string]: ecs.Secret };
  readonly service: {
    name: string;
    repositoryName: string;
  };
  readonly desiredCount: number;
  readonly minCapacity: number;
  readonly maxCapacity: number;
}

export class CommonService extends Construct {
  public readonly fargateService: ecs.FargateService;

  constructor(scope: Construct, id: string, props: IProps) {
    super(scope, id);

    const taskDefinition = this.newTaskDefinition(props);
    this.fargateService = this.newFargateService(taskDefinition, props);
  }

  private newTaskDefinition(props: IProps): ecs.TaskDefinition {
    const ns = this.node.tryGetContext('ns') as string;

    const taskDefinition = new ecs.FargateTaskDefinition(
      this,
      `TaskDefinition`,
      {
        family: `${ns}${props.service.name}`,
        taskRole: props.taskRole,
        executionRole: props.taskExecutionRole,
        runtimePlatform: {
          operatingSystemFamily: ecs.OperatingSystemFamily.LINUX,
          cpuArchitecture: ecs.CpuArchitecture.X86_64,
        },
        cpu: 256,
        memoryLimitMiB: 1024,
      }
    );

    // service container
    const serviceRepository = ecr.Repository.fromRepositoryName(
      this,
      `ServiceRepository`,
      props.service.repositoryName
    );
    const logging = new ecs.AwsLogDriver({
      logGroup: props.taskLogGroup,
      streamPrefix: ns.toLowerCase(),
    });
    taskDefinition.addContainer(`Container`, {
      containerName: props.service.name.toLowerCase(),
      image: ecs.ContainerImage.fromEcrRepository(serviceRepository),
      logging,
      healthCheck: {
        command: [
          'CMD-SHELL',
          'curl -f http://localhost:8090/health/ || exit 1',
        ],
      },
      portMappings: [{ containerPort: 8090, protocol: ecs.Protocol.TCP }],
      secrets: props.taskEnvs,
    });
    taskDefinition.addContainer(`OTelContainer`, {
      containerName: 'aws-otel-collector',
      image: ecs.ContainerImage.fromRegistry(
        'public.ecr.aws/aws-observability/aws-otel-collector'
      ),
      command: ['--config=/etc/ecs/ecs-default-config.yaml'],
      portMappings: [
        { containerPort: 4317, protocol: ecs.Protocol.TCP },
        { containerPort: 4318, protocol: ecs.Protocol.TCP },
        { containerPort: 2000, protocol: ecs.Protocol.UDP },
      ],
    });

    return taskDefinition;
  }

  private newFargateService(
    taskDefinition: ecs.TaskDefinition,
    props: IProps
  ): ecs.FargateService {
    const ns = this.node.tryGetContext('ns') as string;

    const service = new ecs.FargateService(this, 'FargateService', {
      serviceName: `${ns}${props.service.name}`,
      platformVersion: ecs.FargatePlatformVersion.LATEST,
      cluster: props.cluster,
      vpcSubnets: { subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS },
      desiredCount: 1,
      taskDefinition,
      securityGroups: [props.taskSecurityGroup],
      cloudMapOptions: {
        name: props.service.name.toLowerCase(),
        containerPort: 8090,
      },
    });

    const scalableTarget = service.autoScaleTaskCount({
      minCapacity: 1,
      maxCapacity: 3,
    });
    scalableTarget.scaleOnCpuUtilization('CpuScaling', {
      targetUtilizationPercent: 70,
    });
    scalableTarget.scaleOnMemoryUtilization('MemoryScaling', {
      targetUtilizationPercent: 70,
    });

    return service;
  }
}
