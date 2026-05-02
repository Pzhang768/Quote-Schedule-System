import api from "./index";

export interface Job {
  id: string;
  quote_id: string;
  technician_id: string;
  manager_id: string;
  starts_at: string;
  ends_at: string;
  status: "scheduled" | "completed" | "cancelled";
  completed_at: string | null;
  created_at: string;
  customer_name?: string;
  address?: string;
}

export interface AssignJobPayload {
  quote_id: string;
  technician_id: string;
  manager_id: string;
  starts_at: string;
}

export const assignJob = (payload: AssignJobPayload) =>
  api.post<{ data: Job }>("/jobs", payload).then((res) => res.data.data);

export const completeJob = (jobId: string, technicianId: string) =>
  api
    .patch<{ data: Job }>(`/jobs/${jobId}/complete`, { technician_id: technicianId })
    .then((res) => res.data.data);

export const getTechnicianJobs = (technicianId: string, date: string) =>
  api
    .get<{ data: Job[] }>(`/technicians/${technicianId}/jobs`, { params: { date } })
    .then((res) => res.data.data);
