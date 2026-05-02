"use client";
import useNotification from "@/hooks/useNotification/useNotification";
import { usePathname, useRouter } from "next/navigation";

interface NotificationsProps {
  recipientType: "manager" | "technician";
  recipientId: string;
}

function Notifications({ recipientType, recipientId }: NotificationsProps) {
  const { notifications, markRead } = useNotification(recipientType, recipientId);
  return (
    <div className="mt-6">
      <h2 className="text-heading mb-2">Notifications</h2>
      {notifications.map((n) => (
        <div
          key={n.id}
          className={`p-2 mb-1 rounded cursor-pointer text-body ${n.read_at ? "opacity-50" : ""}`}
          onClick={() => markRead(n.id)}
        >
          {n.message}
          <div className="text-caption text-muted">
            {new Date(n.created_at).toLocaleTimeString([], {
              hour: "2-digit",
              minute: "2-digit",
              hour12: false,
            })}
          </div>
        </div>
      ))}
    </div>
  );
}

const SideBar = () => {
  const router = useRouter();
  const pathname = usePathname();

  const segments = pathname.split("/");
  const role = segments[2] as "manager" | "technician" | undefined;
  const recipientId = segments[3] ?? null;
  const recipientType = role === "manager" || role === "technician" ? role : null;

  return (
    <div className="w-60 flex flex-col p-4 border-r border-ink/20">
      <button
        className={`w-full text-left p-4 border border-divider rounded-xl mb-2 mt-8 ${role === "manager" ? "bg-ink text-white border-ink!" : "hover:bg-accent-brass/20"}`}
        onClick={() => router.push("/dashboard/manager")}
      >
        <div className="text-body">Manager View</div>
        <div className={`text-caption ${role === "manager" ? "text-white/60" : "text-muted"}`}>
          Allocate work
        </div>
      </button>

      <button
        className={`w-full text-left p-4 border border-divider rounded-xl ${role === "technician" ? "bg-ink text-white border-ink!" : "hover:bg-accent-brass/20"}`}
        onClick={() => router.push("/dashboard/technician")}
      >
        <div className="text-body">Technician View</div>
        <div className={`text-caption ${role === "technician" ? "text-white/60" : "text-muted"}`}>
          Today&apos;s jobs
        </div>
      </button>

      {recipientType && recipientId && (
        <Notifications
          key={`${recipientType}-${recipientId}`}
          recipientType={recipientType}
          recipientId={recipientId}
        />
      )}
    </div>
  );
};

export default SideBar;
