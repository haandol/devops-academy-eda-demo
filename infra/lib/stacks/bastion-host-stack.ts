import { Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as iam from 'aws-cdk-lib/aws-iam';

interface IProps extends StackProps {
  vpcId: string;
  mskSecurityGroupId: string;
}

export class BastionHostStack extends Stack {
  constructor(scope: Construct, id: string, props: IProps) {
    super(scope, id, props);

    const vpc = ec2.Vpc.fromLookup(this, 'Vpc', { vpcId: props.vpcId });
    this.newBastionHost(vpc, props);
  }

  newBastionHost(vpc: ec2.IVpc, props: IProps): ec2.BastionHostLinux {
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
        ec2.InstanceClass.M5,
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

    const mskSecurityGroup = ec2.SecurityGroup.fromSecurityGroupId(
      this,
      `MskSecurityGorup`,
      props.mskSecurityGroupId
    );
    mskSecurityGroup.addIngressRule(
      bastionHost.connections.securityGroups[0],
      ec2.Port.tcpRange(9092, 9094),
      'BastionHost'
    );

    return bastionHost;
  }
}
