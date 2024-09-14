import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
// import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import {
  RestApi,
  LambdaIntegration,
  ProxyResource,
  LambdaRestApi,
} from 'aws-cdk-lib/aws-apigateway';
import * as go from '@aws-cdk/aws-lambda-go-alpha';

export class AwsLambdaGoStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const main = new go.GoFunction(this, 'main', {
      entry: 'api',
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
