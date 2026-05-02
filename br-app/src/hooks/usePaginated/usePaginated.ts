import { useCallback, useEffect, useState } from "react";

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

const usePaginated = <T>(
  fetcher: (page: number, pageSize: number) => Promise<PaginatedResponse<T>>,
  pageSize = 5
) => {
  const [data, setData] = useState<T[]>([]);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(1);
  const [page, setPage] = useState(1);

  const fetch = useCallback(() => {
    let cancelled = false;
    fetcher(page, pageSize).then((res) => {
      if (!cancelled) {
        setData(res.data);
        setTotal(res.total);
        setTotalPages(res.total_pages);
      }
    });
    return () => {
      cancelled = true;
    };
  }, [fetcher, page, pageSize]);

  useEffect(() => {
    return fetch();
  }, [fetch]);

  return {
    data,
    total,
    totalPages,
    page,
    refetch: fetch,
    nextPage: () => setPage((p) => Math.min(p + 1, totalPages)),
    prevPage: () => setPage((p) => Math.max(p - 1, 1)),
    setPage,
  };
};

export default usePaginated;
