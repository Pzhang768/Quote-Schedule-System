"use client";
import { useCallback, useState } from "react";
import { useParams } from "next/navigation";
import { getTechnicianJobs, Job } from "@/api/jobs";
import usePaginated from "@/hooks/usePaginated/usePaginated";
import JobCard from "./_components/JobCard/JobCard";

const today = () => new Date().toLocaleDateString("en-CA");

export default function TechnicianDashboard() {
  const { id: technicianId } = useParams<{ id: string }>();
  const [date, setDate] = useState(today());

  const fetcher = useCallback(
    (page: number, ps: number) =>
      getTechnicianJobs(technicianId, date).then((jobs) => ({
        data: jobs,
        total: jobs.length,
        page,
        page_size: ps,
        total_pages: 1,
      })),
    [technicianId, date]
  );

  const { data: jobs, refetch } = usePaginated<Job>(fetcher);

  return (
    <div className="flex flex-col h-full w-full p-4 gap-4">
      <div className="border border-divider rounded-xl px-6 py-4 flex items-center justify-between">
        <h1 className="text-title">My Jobs</h1>
        <input
          type="date"
          value={date}
          onChange={(e) => setDate(e.target.value)}
          className="text-caption border border-divider rounded-lg px-2 py-1 bg-transparent"
        />
      </div>

      <ul className="flex flex-col gap-2">
        {jobs.length === 0 && <p className="text-body text-muted p-2">No jobs on this day.</p>}
        {jobs.map((job) => (
          <li key={job.id}>
            <JobCard job={job} technicianId={technicianId} onCompleted={refetch} />
          </li>
        ))}
      </ul>
    </div>
  );
}
