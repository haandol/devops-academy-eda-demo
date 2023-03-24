import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';

interface IProps extends cdk.StackProps {
  readonly tableName: string;
}

export class DatabaseStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: IProps) {
    super(scope, id, props);

    const eventTable = new dynamodb.Table(this, 'TripTable', {
      tableName: props.tableName,
      partitionKey: {
        name: 'PK',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'SK',
        type: dynamodb.AttributeType.STRING,
      },
      encryption: dynamodb.TableEncryption.AWS_MANAGED,
      contributorInsightsEnabled: true,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });
    eventTable.addGlobalSecondaryIndex({
      indexName: 'GSI1',
      partitionKey: {
        name: 'GSI1PK',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'GSI1SK',
        type: dynamodb.AttributeType.STRING,
      },
    });
  }
}
