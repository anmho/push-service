import {PushNotificationRequest} from "../async/send-scheduled-notification";
import now = jest.now;
import {format} from "date-fns";

type ScheduleNotificationRequest = {
    recipient_push_tokens: string[];
    title: string;
    body: string;
    send_time: string; // ISO 8601 format date string
};



describe("scheduled notification", () => {
    const tests: {
        desc: string
        request: ScheduleNotificationRequest
        expectedStatus: number
    }[] = [
        {
            desc: "happy path. valid scheduled notification",
            request: {
                recipient_push_tokens: ["ExponentPushToken[1IoX9FHzHNRx6doi7h3nJm]"],
                title: "scheduled notification",
                body: "scheduled notification body from test",
                send_time: (() => {
                    const now = new Date()
                    now.setMinutes(now.getMinutes() + 1)
                    return now.toISOString().split(".")[0]
                })()
            },
            expectedStatus: 200
        }
    ]

    for (const tc of tests) {
        it(tc.desc, async () => {
            const url = "https://ftlblm5tka.execute-api.us-west-2.amazonaws.com/prod/schedule-push"
            console.log(tc.request)
            const res = await fetch(url, {
                method: "POST",
                body: JSON.stringify(tc.request)
            })

            expect(res.status).toBe(200)
            const data = await res.json()
            console.log(data)
        })
    }
})