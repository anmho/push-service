import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { LambdaRestApi } from 'aws-cdk-lib/aws-apigateway';
import * as go from '@aws-cdk/aws-lambda-go-alpha';
import { env } from './env';
import { AttributeType, Table } from 'aws-cdk-lib/aws-dynamodb';
import { RemovalPolicy } from 'aws-cdk-lib';

function getStageName(environment: string): string {
  switch (environment) {
    case 'production':
      return 'prod';
    case 'staging':
      return 'staging';
    default:
      return 'prod';
  }
}

export class AwsLambdaGoStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const booksTable = new Table(this, 'books', {
      partitionKey: {
        name: 'id',
        type: AttributeType.STRING,
      },
      tableName: 'books',

      /**
       *  The default removal policy is RETAIN, which means that cdk destroy will not attempt to delete
       * the new table, and it will remain in your account until manually deleted. By setting the policy to
       * DESTROY, cdk destroy will delete the table (even if it has data in it)
       */
    });
    const environment = this.node.tryGetContext('environment') || 'production';
    console.log(environment);
    booksTable.applyRemovalPolicy(RemovalPolicy.DESTROY);

    const main = new go.GoFunction(this, 'main', {
      entry: 'lambda/cmd',
      environment: { ...env },
    });
    const api = new LambdaRestApi(this, 'myGateway', {
      handler: main,
      defaultCorsPreflightOptions: {
        allowOrigins: ['*'],
        allowMethods: ['GET', 'POST', 'OPTIONS', 'DELETE', 'PUT'],
      },
      deployOptions: {
        stageName: getStageName(environment),
      },
    });
  }
}
