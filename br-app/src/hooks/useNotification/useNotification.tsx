import { useEffect, useState } from "react";
import {
  getNotifications,
  streamNotifications,
  markNotificationRead,
  Notification,
} from "@/api/notifications";

const useNotification = (recipientType: "technician" | "manager", recipientId: string | null) => {
  const [notifications, setNotifications] = useState<Notification[]>([]);

  useEffect(() => {
    if (!recipientId) return;

    getNotifications(recipientType, recipientId).then(setNotifications);

    const eventSource = streamNotifications(recipientType, recipientId);

    eventSource.onmessage = (event) => {
      const notification: Notification = JSON.parse(event.data);
      setNotifications((prev) => [notification, ...prev]);
    };

    return () => eventSource.close();
  }, [recipientType, recipientId]);

  const markRead = async (id: string) => {
    if (!recipientId) return;
    await markNotificationRead(id, recipientId);
    setNotifications((prev) =>
      prev.map((notification) =>
        notification.id === id
          ? { ...notification, read_at: new Date().toISOString() }
          : notification
      )
    );
  };

  return { notifications, markRead };
};

export default useNotification;
