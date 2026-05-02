"use client";
import { useCallback } from "react";
import { useRouter } from "next/navigation";
import usePaginated from "@/hooks/usePaginated/usePaginated";
import { getManagers } from "@/api/managers";
import { Manager } from "@/api/managers";

export default function ManagerPicker() {
  const fetcher = useCallback((page: number, ps: number) => getManagers(page, ps), []);
  const { data: managers } = usePaginated<Manager>(fetcher);
  const router = useRouter();

  if (managers.length === 0)
    return (
      <div className="p-8 text-body text-muted">No managers found. Please run the migration.</div>
    );

  return (
    <div className="p-8">
      <h1 className="text-title mb-4">Select a manager</h1>
      <ul className="flex flex-col gap-2">
        {managers.map((m: Manager) => (
          <li key={m.ID}>
            <button
              className="w-full text-left p-4 border rounded-xl hover:bg-hover"
              onClick={() => router.push(`/dashboard/manager/${m.ID}`)}
            >
              <div className="text-heading">{m.Name}</div>
              <div className="text-caption text-muted">{m.Email}</div>
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
