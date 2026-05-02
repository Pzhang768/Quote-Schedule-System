import { renderHook, act, waitFor } from "@testing-library/react";
import useNotification from "./useNotification";
import { getNotifications, streamNotifications, markNotificationRead } from "@/api/notifications";
import type { Notification } from "@/api/notifications";

jest.mock("@/api/notifications");

const getNotificationsMock = jest.mocked(getNotifications);
const streamNotificationsMock = jest.mocked(streamNotifications);
const markNotificationReadMock = jest.mocked(markNotificationRead);

const makeNotification = (id: string, readAt: string | null = null): Notification => ({
  id,
  message: "",
  type: "job_assigned",
  read_at: readAt,
  created_at: "2026-05-03T07:00:00Z",
});

const createMockEventSource = () => {
  const eventSource = {
    onmessage: null as ((event: MessageEvent) => void) | null,
    close: jest.fn(),
  };
  streamNotificationsMock.mockReturnValue(eventSource as unknown as EventSource);
  return eventSource;
};

describe("useNotification", () => {
  afterEach(() => {
    jest.resetAllMocks();
  });

  test("does not fetch when recipientId is null", () => {
    renderHook(() => useNotification("manager", null));

    expect(getNotificationsMock).not.toHaveBeenCalled();
    expect(streamNotificationsMock).not.toHaveBeenCalled();
  });

  test("fetches notifications and opens stream on mount", async () => {
    const notifications = [makeNotification("n-1")];
    getNotificationsMock.mockResolvedValue(notifications);
    createMockEventSource();

    const { result } = renderHook(() => useNotification("manager", "m-1"));

    await waitFor(() => expect(result.current.notifications).toEqual(notifications));

    expect(getNotificationsMock).toHaveBeenCalledWith("manager", "m-1");
    expect(streamNotificationsMock).toHaveBeenCalledWith("manager", "m-1");
  });

  test("prepends incoming SSE notification", async () => {
    const existing = makeNotification("n-1");
    getNotificationsMock.mockResolvedValue([existing]);
    const eventSource = createMockEventSource();

    const { result } = renderHook(() => useNotification("manager", "m-1"));
    await waitFor(() => expect(result.current.notifications).toHaveLength(1));

    const incoming = makeNotification("n-2");
    act(() => {
      eventSource.onmessage?.({ data: JSON.stringify(incoming) } as MessageEvent);
    });

    expect(result.current.notifications[0]).toEqual(incoming);
    expect(result.current.notifications[1]).toEqual(existing);
  });

  test("closes event source on unmount", async () => {
    getNotificationsMock.mockResolvedValue([]);
    const eventSource = createMockEventSource();

    const { unmount } = renderHook(() => useNotification("manager", "m-1"));
    await waitFor(() => expect(streamNotificationsMock).toHaveBeenCalled());

    unmount();

    expect(eventSource.close).toHaveBeenCalled();
  });

  test("markRead calls api and sets read_at on notification", async () => {
    const notification = makeNotification("n-1");
    getNotificationsMock.mockResolvedValue([notification]);
    markNotificationReadMock.mockResolvedValue();
    createMockEventSource();

    const { result } = renderHook(() => useNotification("manager", "m-1"));
    await waitFor(() => expect(result.current.notifications).toHaveLength(1));

    await act(() => result.current.markRead("n-1"));

    expect(markNotificationReadMock).toHaveBeenCalledWith("n-1", "m-1");
    expect(result.current.notifications[0].read_at).not.toBeNull();
  });

  test("markRead does nothing when recipientId is null", async () => {
    const { result } = renderHook(() => useNotification("manager", null));

    await act(() => result.current.markRead("n-1"));

    expect(markNotificationReadMock).not.toHaveBeenCalled();
  });

  test("ignores fetch result after unmount", async () => {
    let resolve: (value: Notification[]) => void;
    getNotificationsMock.mockReturnValue(
      new Promise<Notification[]>((res) => {
        resolve = res;
      })
    );
    createMockEventSource();

    const { result, unmount } = renderHook(() => useNotification("manager", "m-1"));

    unmount();
    resolve!([makeNotification("n-1")]);

    await new Promise((r) => setTimeout(r, 0));
    expect(result.current.notifications).toEqual([]);
  });

  test("markRead leaves non-matching notifications unchanged", async () => {
    const n1 = makeNotification("n-1");
    const n2 = makeNotification("n-2");
    getNotificationsMock.mockResolvedValue([n1, n2]);
    markNotificationReadMock.mockResolvedValue();
    createMockEventSource();

    const { result } = renderHook(() => useNotification("manager", "m-1"));
    await waitFor(() => expect(result.current.notifications).toHaveLength(2));

    await act(() => result.current.markRead("n-1"));

    expect(result.current.notifications[0].read_at).not.toBeNull();
    expect(result.current.notifications[1].read_at).toBeNull();
  });

  test("resets notifications when recipientId changes", async () => {
    const notifications1 = [makeNotification("n-1")];
    const notifications2 = [makeNotification("n-2")];
    getNotificationsMock
      .mockResolvedValueOnce(notifications1)
      .mockResolvedValueOnce(notifications2);
    createMockEventSource();

    const { result, rerender } = renderHook(
      ({ id }: { id: string }) => useNotification("manager", id),
      { initialProps: { id: "m-1" } }
    );

    await waitFor(() => expect(result.current.notifications).toHaveLength(1));

    createMockEventSource();
    rerender({ id: "m-2" });

    await waitFor(() => expect(result.current.notifications[0].id).toBe("n-2"));
  });
});
