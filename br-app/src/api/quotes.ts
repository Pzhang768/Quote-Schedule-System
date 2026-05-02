import api from "./index";
import { PaginatedResponse } from "@/hooks/usePaginated/usePaginated";

export interface Quote {
  ID: string;
  CustomerName: string;
  Address: string;
  Description: string;
  Status: "unscheduled" | "scheduled";
  CreatedAt: string;
  UpdatedAt: string;
}

export const getQuotes = (page: number, pageSize: number, status?: string) =>
  api
    .get<
      PaginatedResponse<Quote>
    >("/quotes", { params: { page, page_size: pageSize, ...(status ? { status } : {}) } })
    .then((res) => res.data);
