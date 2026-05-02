import { useState } from "react";
import { assignJob } from "@/api/jobs";
import { Quote } from "@/api/quotes";
import { Technician } from "@/api/technicians";

interface Props {
  managerId: string;
  quote: Quote;
  technician: Technician;
  date: string;
  time: string;
  onSuccess: () => void;
}

export default function AssignForm({ managerId, quote, technician, date, time, onSuccess }: Props) {
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: { preventDefault: () => void }) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      await assignJob({
        quote_id: quote.ID,
        technician_id: technician.ID,
        manager_id: managerId,
        starts_at: `${date}T${time}:00Z`,
      });
      onSuccess();
    } catch (err: unknown) {
      const status = (err as { response?: { status?: number; data?: { error?: string } } })
        .response;
      if (status?.status === 409) {
        setError("This technician has a conflicting job at that time.");
      } else {
        setError(status?.data?.error ?? "Something went wrong.");
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="border border-divider rounded-xl p-4 flex flex-col gap-3"
    >
      <div className="text-caption-upper text-muted">Assign Job</div>

      <div className="flex flex-col gap-1">
        <div className="text-caption text-muted">Quote</div>
        <div className="text-heading">{quote.CustomerName}</div>
        <div className="text-body text-muted">{quote.Address}</div>
      </div>

      <div className="flex flex-col gap-1">
        <div className="text-caption text-muted">Technician</div>
        <div className="text-heading">{technician.Name}</div>
      </div>

      <div className="flex flex-col gap-1">
        <div className="text-caption text-muted">Start time</div>
        <div className="text-body">
          {date} at {time}
        </div>
      </div>

      {error && <div className="text-caption text-accent">{error}</div>}

      <button
        type="submit"
        disabled={submitting}
        className="text-body bg-ink text-white rounded-xl px-4 py-2 disabled:opacity-50"
      >
        {submitting ? "Assigning..." : "Assign Job"}
      </button>
    </form>
  );
}
