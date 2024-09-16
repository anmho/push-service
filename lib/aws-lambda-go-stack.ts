import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { LambdaRestApi } from 'aws-cdk-lib/aws-apigateway';
import * as go from '@aws-cdk/aws-lambda-go-alpha';

export class AwsLambdaGoStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const main = new go.GoFunction(this, 'main', {
      entry: 'lambda/cmd',
    });
    const gateway = new LambdaRestApi(this, 'myGateway', {
      handler: main,
      defaultCorsPreflightOptions: {
        allowOrigins: ['*'],
        allowMethods: ['GET', 'POST', 'OPTIONS', 'DELETE', 'PUT'],
      },
    });
  }
}
