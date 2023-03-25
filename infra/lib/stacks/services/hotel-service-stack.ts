import { Stack } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { IServiceProps } from '../../interfaces/types';
import { CommonService } from '../../constructs/service';
import * as ecs from 'aws-cdk-lib/aws-ecs';
import * as ssm from 'aws-cdk-lib/aws-ssm';

interface IProps extends IServiceProps {
  kafkaSeeds: string;
}

export class HotelServiceStack extends Stack {
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

    new CommonService(this, 'HotelService', {
      ...props,
      taskEnvs,
    });
  }
}
