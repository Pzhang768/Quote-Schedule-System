import api from "./index";
import { PaginatedResponse } from "@/hooks/usePaginated/usePaginated";

export interface Technician {
  ID: string;
  Name: string;
  Email: string;
}

export const getTechnicians = (page: number, pageSize: number) =>
  api
    .get<PaginatedResponse<Technician>>("/technicians", { params: { page, page_size: pageSize } })
    .then((res) => res.data);
