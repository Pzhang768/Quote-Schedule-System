import api from "./index";

export interface Notification {
  id: string;
  message: string;
  type: "job_assigned" | "job_updated" | "job_completed";
  read_at: string | null;
  created_at: string;
}

export const getNotifications = (recipientType: string, recipientId: string) =>
  api
    .get<{ data: Notification[] }>("/notifications", {
      params: { recipient_type: recipientType, recipient_id: recipientId },
    })
    .then((res) => res.data.data);

export const streamNotifications = (recipientType: string, recipientId: string) =>
  new EventSource(
    `${process.env.NEXT_PUBLIC_API_URL}/api/v1/notifications/stream?recipient_type=${recipientType}&recipient_id=${recipientId}`
  );

export const markNotificationRead = (id: string, recipientId: string): Promise<void> =>
  api
    .patch(`/notifications/${id}/read`, null, {
      params: { recipient_id: recipientId },
    })
    .then(() => {});
