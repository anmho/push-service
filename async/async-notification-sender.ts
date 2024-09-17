
import { EventBridgeEvent } from "aws-lambda";

interface ScheduledNotification {
    recipients: string[];
    title: string;
    body: string;
}
export const handler = async (
    event: EventBridgeEvent<"ReminderNotification", ScheduledNotification>
) => {
    console.log("result", JSON.stringify(event))
};
