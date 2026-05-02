import { renderHook, act, waitFor } from "@testing-library/react";
import usePaginated, { PaginatedResponse } from "./usePaginated";

const makeFetcher = <T>(pages: PaginatedResponse<T>[]) => {
  const fetcher = jest.fn();
  pages.forEach((page, i) => fetcher.mockResolvedValueOnce(pages[i]));
  fetcher.mockResolvedValue(pages[pages.length - 1]);
  return fetcher;
};

const makePage = <T>(data: T[], page: number, totalPages: number): PaginatedResponse<T> => ({
  data,
  total: totalPages * 5,
  page,
  page_size: 5,
  total_pages: totalPages,
});

describe("usePaginated", () => {
  afterEach(() => {
    jest.resetAllMocks();
  });

  test("fetches first page on mount", async () => {
    const fetcher = makeFetcher([makePage(["a", "b"], 1, 2)]);
    const { result } = renderHook(() => usePaginated(fetcher));

    await waitFor(() => expect(result.current.data).toEqual(["a", "b"]));

    expect(fetcher).toHaveBeenCalledWith(1, 5);
    expect(result.current.page).toBe(1);
    expect(result.current.totalPages).toBe(2);
  });

  test("nextPage increments page and refetches", async () => {
    const fetcher = makeFetcher([makePage(["a"], 1, 2), makePage(["b"], 2, 2)]);
    const { result } = renderHook(() => usePaginated(fetcher));

    await waitFor(() => expect(result.current.data).toEqual(["a"]));

    act(() => result.current.nextPage());

    await waitFor(() => expect(result.current.data).toEqual(["b"]));
    expect(result.current.page).toBe(2);
  });

  test("prevPage decrements page and refetches", async () => {
    const fetcher = makeFetcher([
      makePage(["a"], 1, 2),
      makePage(["b"], 2, 2),
      makePage(["a"], 1, 2),
    ]);
    const { result } = renderHook(() => usePaginated(fetcher));

    await waitFor(() => expect(result.current.data).toEqual(["a"]));

    act(() => result.current.nextPage());
    await waitFor(() => expect(result.current.page).toBe(2));

    act(() => result.current.prevPage());
    await waitFor(() => expect(result.current.page).toBe(1));
  });

  test("nextPage does not exceed totalPages", async () => {
    const fetcher = makeFetcher([makePage(["a"], 1, 1)]);
    const { result } = renderHook(() => usePaginated(fetcher));

    await waitFor(() => expect(result.current.data).toEqual(["a"]));

    act(() => result.current.nextPage());

    expect(result.current.page).toBe(1);
  });

  test("prevPage does not go below 1", async () => {
    const fetcher = makeFetcher([makePage(["a"], 1, 2)]);
    const { result } = renderHook(() => usePaginated(fetcher));

    await waitFor(() => expect(result.current.data).toEqual(["a"]));

    act(() => result.current.prevPage());

    expect(result.current.page).toBe(1);
  });

  test("refetch re-calls fetcher", async () => {
    const fetcher = makeFetcher([makePage(["a"], 1, 1), makePage(["b"], 1, 1)]);
    const { result } = renderHook(() => usePaginated(fetcher));

    await waitFor(() => expect(result.current.data).toEqual(["a"]));

    act(() => result.current.refetch());

    await waitFor(() => expect(result.current.data).toEqual(["b"]));
  });

  test("ignores fetch result after unmount", async () => {
    let resolve: (value: PaginatedResponse<string>) => void;
    const fetcher = jest.fn(
      () =>
        new Promise<PaginatedResponse<string>>((res) => {
          resolve = res;
        })
    );

    const { result, unmount } = renderHook(() => usePaginated(fetcher));

    unmount();
    resolve!(makePage(["a"], 1, 1));

    await new Promise((r) => setTimeout(r, 0));
    expect(result.current.data).toEqual([]);
  });
});
