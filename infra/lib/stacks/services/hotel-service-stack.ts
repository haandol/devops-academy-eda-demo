import { Stack } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { IServiceProps } from '../../interfaces/types';
import { CommonService } from '../../constructs/service';
import * as ecs from 'aws-cdk-lib/aws-ecs';
import * as ssm from 'aws-cdk-lib/aws-ssm';

export class HotelServiceStack extends Stack {
  constructor(scope: Construct, id: string, props: IServiceProps) {
    super(scope, id, props);

    const taskEnvs = {
      OTEL_EXPORTER_OTLP_ENDPOINT: ecs.Secret.fromSsmParameter(
        new ssm.StringParameter(this, 'EnvOtelExporterEndpoint', {
          stringValue: 'aws-otel-collector:4317',
        })
      ),
    };

    new CommonService(this, 'HotelService', {
      ...props,
      taskEnvs,
    });
  }
}
