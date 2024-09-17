import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { LambdaRestApi } from 'aws-cdk-lib/aws-apigateway';
import * as go from '@aws-cdk/aws-lambda-go-alpha';
import { env } from './env';
import { AttributeType, Table } from 'aws-cdk-lib/aws-dynamodb';
import { RemovalPolicy } from 'aws-cdk-lib';
import {Effect, Policy, PolicyStatement, Role, ServicePrincipal} from "aws-cdk-lib/aws-iam";
import * as events from "aws-cdk-lib/aws-events";
import {LambdaFunction} from "aws-cdk-lib/aws-events-targets";
import {CfnScheduleGroup} from "aws-cdk-lib/aws-scheduler";


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
    PushTable.applyRemovalPolicy(RemovalPolicy.DESTROY);

    // // Core infra: Eventbridge event bus
    // const eventBus = new cdk.aws_events.EventBus(this, `${id}-event-bus`, {
    //   eventBusName: `${id}-event-bus`,
    // });

    // // need to create a service-linked role and policy for
    // // the scheduler to be able to put events onto our bus
    // const schedulerRole = new Role(this, `${id}-scheduler-role`, {
    //   assumedBy: new ServicePrincipal("scheduler.amazonaws.com"),
    // });

    // new Policy(this, `${id}-schedule-policy`, {
    //   policyName: "ScheduleToPutEvents",
    //   roles: [schedulerRole],
    //   statements: [
    //     new PolicyStatement({
    //       effect: Effect.ALLOW,
    //       actions: ["events:PutEvents"],
    //       resources: [eventBus.eventBusArn],
    //     }),
    //   ],
    // });

    const main = new go.GoFunction(this, 'main', {
      entry: 'lambda/cmd',
      environment: {
        ...env,
        // EVENTBUS_ARN: eventBus.eventBusArn,
      },
      initialPolicy: [
        // Give lambda permission to create the group & schedule and pass IAM role to the scheduler
        new PolicyStatement({
          actions: [
            "scheduler:CreateSchedule",
            "scheduler:CreateScheduleGroup",
            "iam:PassRole",
          ],
          resources: ["*"],
        }),
      ],
    });


    // new events.Rule(this, "ReminderNotification", {
    //   description: "Send a push notification reminding user of a booked court",
    //   eventPattern: {
    //     source: ["scheduler.notifications"],
    //     detailType: ["ReminderNotification"],
    //   },
    //   ]
    //   eventBus,
    // }).addTarget(new LambdaFunction(main, ));




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
