
import { EventBridgeEvent } from "aws-lambda";
import {z} from "zod";


const PushNotificationRequestSchema = z.object({
    recipient_push_tokens: z.array(z.string()),
    title: z.string(),
    body: z.string(),
    data: z.any(),
})


export type PushNotificationRequest = z.infer<typeof PushNotificationRequestSchema>;

export const handler = async (
    // event: EventBridgeEvent<"ScheduledNotification", PushNotificationRequest>
    detail: any
) => {
    // console.error("EVENT", event)
    // if (!event.detail) {
    //     throw new Error(JSON.stringify(event))
    // }
    // console.log(request)
    console.log("event/detail", detail)
    const request = PushNotificationRequestSchema.parse(detail)
    console.log("REQUEST", request)

    console.log('handling event', detail)


    const res = await fetch("https://ftlblm5tka.execute-api.us-west-2.amazonaws.com/prod/send-push", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(
            request
            // {
            //     "recipient_push_tokens": [ "ExponentPushToken[1IoX9FHzHNRx6doi7h3nJm]" ],
            //     "title": "scheduled notification",
            //     "body": "scheduled notification body 2"
            // }
        )
    })
    console.log(res)
    const data = await res.json()
    console.log(data)
};