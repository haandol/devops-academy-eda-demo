import { Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as iam from 'aws-cdk-lib/aws-iam';

interface IProps extends StackProps {
  vpcId: string;
}

export class BastionHostStack extends Stack {
  constructor(scope: Construct, id: string, props: IProps) {
    super(scope, id, props);

    const vpc = ec2.Vpc.fromLookup(this, 'Vpc', { vpcId: props.vpcId });
    this.newBastionHost(vpc);
  }

  private newBastionHost(vpc: ec2.IVpc): ec2.BastionHostLinux {
    const bastionHost = new ec2.BastionHostLinux(this, `BastionHost`, {
      vpc,
      blockDevices: [
        {
          deviceName: '/dev/xvda',
          volume: ec2.BlockDeviceVolume.ebs(128, {
            encrypted: true,
          }),
        },
      ],
      instanceType: ec2.InstanceType.of(
        ec2.InstanceClass.T3,
        ec2.InstanceSize.LARGE
      ),
    });
    bastionHost.role.addManagedPolicy(
      iam.ManagedPolicy.fromAwsManagedPolicyName(
        `AmazonEC2ContainerRegistryFullAccess`
      )
    );
    bastionHost.role.addManagedPolicy(
      iam.ManagedPolicy.fromAwsManagedPolicyName(`AWSCodeCommitPowerUser`)
    );
    bastionHost.role.addManagedPolicy(
      iam.ManagedPolicy.fromAwsManagedPolicyName(`AmazonSSMManagedInstanceCore`)
    );
    return bastionHost;
  }
}
