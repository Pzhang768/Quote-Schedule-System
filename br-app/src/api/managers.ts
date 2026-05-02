import api from "./index";
import { PaginatedResponse } from "@/hooks/usePaginated/usePaginated";

export interface Manager {
  ID: string;
  Name: string;
  Email: string;
}

export const getManagers = (page: number, pageSize: number) =>
  api
    .get<PaginatedResponse<Manager>>("/managers", { params: { page, page_size: pageSize } })
    .then((res) => res.data);
