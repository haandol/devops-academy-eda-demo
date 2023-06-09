import { Stack } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { IServiceProps } from '../../interfaces/types';
import { CommonService } from '../../constructs/service';
import * as elbv2 from 'aws-cdk-lib/aws-elasticloadbalancingv2';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as ecs from 'aws-cdk-lib/aws-ecs';
import * as ssm from 'aws-cdk-lib/aws-ssm';

interface IProps extends IServiceProps {
  vpc: ec2.IVpc;
  alb: elbv2.IApplicationLoadBalancer;
  kafkaSeeds: string;
}

export class TripServiceStack extends Stack {
  constructor(scope: Construct, id: string, props: IProps) {
    super(scope, id, props);

    const taskEnvs = {
      OTEL_EXPORTER_OTLP_ENDPOINT: ecs.Secret.fromSsmParameter(
        new ssm.StringParameter(this, 'EnvOtelExporterEndpoint', {
          stringValue: '127.0.0.1:4317',
        })
      ),
      KAFKA_SEEDS: ecs.Secret.fromSsmParameter(
        new ssm.StringParameter(this, 'EnvKafkaSeeds', {
          stringValue: props.kafkaSeeds,
        })
      ),
    };

    const tripService = new CommonService(this, 'TripService', {
      ...props,
      taskEnvs,
    });

    this.registerServiceToLoadBalancer(tripService.fargateService, props);
  }

  registerServiceToLoadBalancer(
    fargateService: ecs.FargateService,
    props: IProps
  ) {
    const targetGroup = new elbv2.ApplicationTargetGroup(this, 'ListenerRule', {
      protocol: elbv2.ApplicationProtocol.HTTP,
      port: props.service.port,
      vpc: props.vpc,
      targets: [fargateService],
      healthCheck: {
        enabled: true,
        path: '/healthz',
        healthyHttpCodes: '200-299',
      },
    });

    new elbv2.ApplicationListener(this, 'Listener', {
      loadBalancer: props.alb,
      protocol: elbv2.ApplicationProtocol.HTTP,
      port: 80,
      defaultTargetGroups: [targetGroup],
    });
  }
}
