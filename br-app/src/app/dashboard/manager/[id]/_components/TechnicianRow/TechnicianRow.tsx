"use client";
import { useEffect, useState } from "react";
import { getTechnicianJobs, Job } from "@/api/jobs";
import { Technician } from "@/api/technicians";

interface Props {
  technician: Technician;
  date: string;
  selected: boolean;
  onClick: () => void;
  onTimeSelect: (time: string) => void;
  selectedTime: string;
}

function buildAvailableSlots(date: string, jobs: Job[]): string[] {
  const slots: string[] = [];
  for (let h = 7; h < 17; h++) {
    for (const m of [0, 30]) {
      const slotStart = new Date(
        `${date}T${String(h).padStart(2, "0")}:${String(m).padStart(2, "0")}:00Z`
      );
      const slotEnd = new Date(slotStart.getTime() + 2 * 60 * 60 * 1000);
      const blocked = jobs.some((j) => {
        const jStart = new Date(j.starts_at);
        const jEnd = new Date(j.ends_at);
        return slotStart < jEnd && slotEnd > jStart;
      });
      if (!blocked) {
        slots.push(`${String(h).padStart(2, "0")}:${String(m).padStart(2, "0")}`);
      }
    }
  }
  return slots;
}

export default function TechnicianRow({
  technician,
  date,
  selected,
  onClick,
  onTimeSelect,
  selectedTime,
}: Props) {
  const [slots, setSlots] = useState<Job[]>([]);
  const [open, setOpen] = useState(false);

  useEffect(() => {
    let cancelled = false;
    getTechnicianJobs(technician.ID, date).then((data) => {
      if (!cancelled) setSlots(data);
    });
    return () => {
      cancelled = true;
    };
  }, [technician.ID, date]);

  const available = buildAvailableSlots(date, slots);

  return (
    <article className="relative">
      <div
        className={`border rounded-xl p-4 cursor-pointer transition-colors ${selected ? "border-ink bg-ink/5" : "border-divider hover:bg-hover"}`}
        onClick={() => {
          onClick();
          setOpen((o) => !o);
        }}
      >
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-heading">{technician.Name}</h3>
            <p className="text-caption text-muted">{technician.Email}</p>
          </div>
          <div className="flex items-center gap-2" onClick={(e) => e.stopPropagation()}>
            <p className="text-caption text-muted">
              {`${slots.length} job${slots.length !== 1 ? "s" : ""} booked`}
            </p>
            <select
              value={selectedTime}
              onChange={(e) => onTimeSelect(e.target.value)}
              className="text-caption border border-divider rounded-lg px-2 py-1 bg-transparent"
            >
              <option value="">Pick time</option>
              {available.map((t) => (
                <option key={t} value={t}>
                  {t}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {slots.length > 0 && selected && open && (
        <ul className="absolute left-0 right-0 top-full z-10 mt-1 border border-divider rounded-xl bg-surface shadow-md p-3 flex flex-col gap-1">
          {slots.map((s) => (
            <li key={s.id} className="flex items-center gap-2">
              <div
                className={`w-2 h-2 rounded-full ${s.status === "completed" ? "bg-ok" : "bg-accent"}`}
              />
              <time className="text-caption text-muted">
                {new Date(s.starts_at).toLocaleTimeString([], {
                  hour: "2-digit",
                  minute: "2-digit",
                  hour12: false,
                  timeZone: "UTC",
                })}
                {" – "}
                {new Date(s.ends_at).toLocaleTimeString([], {
                  hour: "2-digit",
                  minute: "2-digit",
                  hour12: false,
                  timeZone: "UTC",
                })}
              </time>
              <span className="text-caption text-muted capitalize">{s.status}</span>
            </li>
          ))}
        </ul>
      )}
    </article>
  );
}
