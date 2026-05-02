"use client";
import { useState } from "react";
import { completeJob, Job } from "@/api/jobs";

interface Props {
  job: Job;
  technicianId: string;
  onCompleted: () => void;
}

export default function JobCard({ job, technicianId, onCompleted }: Props) {
  const [completing, setCompleting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleComplete = async () => {
    setCompleting(true);
    setError(null);
    try {
      await completeJob(job.id, technicianId);
      onCompleted();
    } catch {
      setError("Failed to complete job.");
    } finally {
      setCompleting(false);
    }
  };

  const start = new Date(job.starts_at).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
    timeZone: "UTC",
  });
  const end = new Date(job.ends_at).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
    timeZone: "UTC",
  });

  return (
    <div className="border border-divider rounded-xl p-4 flex items-center justify-between gap-4 bg-accent-brass/10">
      <div className="flex flex-col gap-1">
        <div className="text-heading">
          {start} – {end}
        </div>
        {job.customer_name && <div className="text-body">{job.customer_name}</div>}
        {job.address && <div className="text-caption text-muted">{job.address}</div>}
        <div className="text-caption text-muted capitalize">{job.status}</div>
        {error && <div className="text-caption text-accent">{error}</div>}
      </div>

      {job.status === "scheduled" && (
        <button
          disabled={completing}
          onClick={handleComplete}
          className="text-body bg-ink text-white rounded-xl px-4 py-2 disabled:opacity-50 shrink-0 cursor-pointer"
        >
          {completing ? "Completing..." : "Complete"}
        </button>
      )}

      {job.status === "completed" && <div className="text-caption text-ok font-medium">Done</div>}
    </div>
  );
}
