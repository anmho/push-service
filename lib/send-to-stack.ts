import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { LambdaRestApi } from 'aws-cdk-lib/aws-apigateway';
import * as go from '@aws-cdk/aws-lambda-go-alpha';
import * as node from "aws-cdk-lib/aws-lambda-nodejs";
import { env } from './env';
import { AttributeType, Table } from 'aws-cdk-lib/aws-dynamodb';
import * as iam from "aws-cdk-lib/aws-iam";
import * as events from "aws-cdk-lib/aws-events";

import {Architecture, Runtime} from "aws-cdk-lib/aws-lambda";
import {LambdaFunction} from "aws-cdk-lib/aws-events-targets";



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

export class SendToStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const PushTable = new Table(this, 'PushTable', {
      partitionKey: {
        name: 'Id',
        type: AttributeType.STRING,
      },
      tableName: 'PushTable',

      /**
       *  The default removal policy is RETAIN, which means that cdk destroy will not attempt to delete
       * the new table, and it will remain in your account until manually deleted. By setting the policy to
       * DESTROY, cdk destroy will delete the table (even if it has data in it)
       */
    });
    const environment = this.node.tryGetContext('environment') || 'production';
    console.log(environment);
    PushTable.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);


    const sendScheduledNotificationLambda = new node.NodejsFunction(
        this,
        "sendNotification",
        {
          functionName: `send-scheduled-notification`,
          runtime: Runtime.NODEJS_20_X,
          architecture: Architecture.ARM_64,
          handler: "handler",
          entry: "async/send-scheduled-notification.ts",
          memorySize: 512,
          timeout: cdk.Duration.seconds(3),
          initialPolicy: [
            new iam.PolicyStatement({
              effect: iam.Effect.ALLOW,
              actions: [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents",
                "logs:DescribeLogStreams",
                "cloudwatch:PutMetricData",
                "iam:PassRole",
              ],
              resources: ["*"]
            })
          ]
        }
    );


    // Core infra: Eventbridge event bus
    // const eventBus = new cdk.aws_events.EventBus(this, `${id}-event-bus`, {
    //   eventBusName: `${id}-event-bus`,
    // });

    // need to create a service-linked role and policy for
    // the scheduler to be able to put events onto our bus
    const schedulerRole = new iam.Role(this, `${id}-scheduler-role`, {
      assumedBy: new iam.ServicePrincipal("scheduler.amazonaws.com"),
    });


    new iam.Policy(this, `${id}-schedule-policy`, {
      policyName: "ScheduleToPutEvents",
      roles: [schedulerRole],
      statements: [
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: [
              // "events:PutEvents",
              "lambda:InvokeFunction"
          ],
          // resources: [eventBus.eventBusArn],
          resources: [sendScheduledNotificationLambda.functionArn],
        }),
      ],
    });
    // new cdk.CfnOutput(this, `${id}-event-bus-output`, { value: eventBus.eventBusArn });
    new cdk.CfnOutput(this, `${id}-scheduler-role-output`, { value: schedulerRole.roleArn });





    // Arn:     aws.String(os.Getenv("SEND_SCHEDULED_NOTIFICATION_ARN")),
    //     RoleArn: aws.String(os.Getenv("INVOKE_SEND_SCHEDULED_NOTIFICATION_ROLE_ARN")),

    const main = new go.GoFunction(this, 'main', {
      entry: 'lambda/cmd',
      environment: {
        ...env,
        SEND_SCHEDULED_NOTIFICATION_ARN: sendScheduledNotificationLambda.functionArn,
        INVOKE_SEND_SCHEDULED_NOTIFICATION_ROLE_ARN: schedulerRole.roleArn,
        // SCHEDULE_ROLE_ARN: schedulerRole.roleArn,
        // EVENTBUS_ARN: eventBus.eventBusArn,
      },

      initialPolicy: [
        // Give lambda permission to create the group & schedule and pass IAM role to the scheduler
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: [
            "scheduler:CreateSchedule",
            "scheduler:CreateScheduleGroup",
            "iam:PassRole",
          ],
          resources: ["*"],
        }),
      ],
    });




    // Rule to match schedules for users and attach our email customer lambda.
    // new events.Rule(this, "ScheduledNotification", {
    //   description: "Send a push notification reminding user of a booked court",
    //   eventPattern: {
    //     source: ["scheduler.notifications"],
    //     detailType: ["ScheduledNotification"],
    //   },
    //   eventBus,
    // }).addTarget(new LambdaFunction(sendScheduledNotificationLambda));

    const api = new LambdaRestApi(this, 'send-to-api', {
      handler: main,
      defaultCorsPreflightOptions: {
        allowOrigins: ['*'],
        allowMethods: ['GET', 'POST', 'OPTIONS', 'DELETE', 'PUT'],
      },
      deployOptions: {
        stageName: getStageName(environment),
        throttlingBurstLimit: 123,
        throttlingRateLimit: 123,
      },
    });
  }
}
