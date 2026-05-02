import api from "./index";
import { getNotifications, streamNotifications, markNotificationRead } from "./notifications";

jest.mock("./index");

const apiMock = jest.mocked(api);

const mockEventSource = {
  onmessage: null,
  close: jest.fn(),
};

beforeAll(() => {
  global.EventSource = jest
    .fn()
    .mockImplementation(() => mockEventSource) as unknown as typeof EventSource;
});

describe("notifications api", () => {
  afterEach(() => {
    jest.resetAllMocks();
  });

  test("getNotifications: fetches notifications for recipient", async () => {
    const notifications = [{ id: "n-1", message: "Job assigned" }];
    (apiMock.get as jest.Mock) = jest.fn().mockResolvedValue({ data: { data: notifications } });

    const result = await getNotifications("manager", "m-1");

    expect(apiMock.get).toHaveBeenCalledWith("/notifications", {
      params: { recipient_type: "manager", recipient_id: "m-1" },
    });
    expect(result).toEqual(notifications);
  });

  test("streamNotifications: creates EventSource with correct url", () => {
    process.env.NEXT_PUBLIC_API_URL = "http://localhost:8081";

    const eventSource = streamNotifications("technician", "t-1");

    expect(global.EventSource).toHaveBeenCalledWith(
      "http://localhost:8081/api/v1/notifications/stream?recipient_type=technician&recipient_id=t-1"
    );
  });

  test("markNotificationRead: patches notification as read", async () => {
    (apiMock.patch as jest.Mock) = jest.fn().mockResolvedValue({});

    await markNotificationRead("n-1", "m-1");

    expect(apiMock.patch).toHaveBeenCalledWith("/notifications/n-1/read", null, {
      params: { recipient_id: "m-1" },
    });
  });
});
