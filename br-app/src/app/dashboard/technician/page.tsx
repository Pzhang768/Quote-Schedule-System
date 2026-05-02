"use client";
import { useCallback } from "react";
import { useRouter } from "next/navigation";
import usePaginated from "@/hooks/usePaginated/usePaginated";
import { getTechnicians, Technician } from "@/api/technicians";

export default function TechnicianPicker() {
  const fetcher = useCallback((page: number, ps: number) => getTechnicians(page, ps), []);
  const { data: technicians } = usePaginated<Technician>(fetcher);
  const router = useRouter();

  if (technicians.length === 0)
    return (
      <div className="p-8 text-body text-muted">
        No technicians found. Please run the migration.
      </div>
    );

  return (
    <div className="p-8">
      <h1 className="text-title mb-4">Select a technician</h1>
      <ul className="flex flex-col gap-2">
        {technicians.map((t: Technician) => (
          <li key={t.ID}>
            <button
              className="w-full text-left p-4 border rounded-xl hover:bg-hover"
              onClick={() => router.push(`/dashboard/technician/${t.ID}`)}
            >
              <div className="text-heading">{t.Name}</div>
              <div className="text-caption text-muted">{t.Email}</div>
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
