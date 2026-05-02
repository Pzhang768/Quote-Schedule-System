import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import SideBar from "./SideBar";
import useNotification from "@/hooks/useNotification/useNotification";

const mockPush = jest.fn();

jest.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
  usePathname: jest.fn(),
}));

jest.mock("@/hooks/useNotification/useNotification");
const useNotificationMock = jest.mocked(useNotification);

import { usePathname } from "next/navigation";
const usePathnameMock = jest.mocked(usePathname);

describe("SideBar", () => {
  afterEach(() => {
    jest.resetAllMocks();
    mockPush.mockReset();
  });

  test("renders Manager View and Technician View buttons", () => {
    usePathnameMock.mockReturnValue("/");
    render(<SideBar />);

    expect(screen.getByRole("button", { name: /Manager View/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /Technician View/i })).toBeInTheDocument();
  });

  test("does not render notifications when no recipient in path", () => {
    usePathnameMock.mockReturnValue("/dashboard/manager");
    render(<SideBar />);

    expect(useNotificationMock).not.toHaveBeenCalled();
  });

  test("renders notifications when on manager detail page", () => {
    usePathnameMock.mockReturnValue("/dashboard/manager/m-1");
    useNotificationMock.mockReturnValue({ notifications: [], markRead: jest.fn() });
    render(<SideBar />);

    expect(useNotificationMock).toHaveBeenCalledWith("manager", "m-1");
    expect(screen.getByText("Notifications")).toBeInTheDocument();
  });

  test("renders notification messages", () => {
    usePathnameMock.mockReturnValue("/dashboard/manager/m-1");
    useNotificationMock.mockReturnValue({
      notifications: [
        {
          id: "n-1",
          message: "Job assigned",
          type: "job_assigned",
          read_at: null,
          created_at: "2026-05-03T07:00:00Z",
        },
      ],
      markRead: jest.fn(),
    });
    render(<SideBar />);

    expect(screen.getByText("Job assigned")).toBeInTheDocument();
  });

  test("calls markRead when notification clicked", async () => {
    const markRead = jest.fn();
    usePathnameMock.mockReturnValue("/dashboard/manager/m-1");
    useNotificationMock.mockReturnValue({
      notifications: [
        {
          id: "n-1",
          message: "Job assigned",
          type: "job_assigned",
          read_at: null,
          created_at: "2026-05-03T07:00:00Z",
        },
      ],
      markRead,
    });
    render(<SideBar />);

    await userEvent.click(screen.getByText("Job assigned"));

    expect(markRead).toHaveBeenCalledWith("n-1");
  });

  test("navigates to manager dashboard when Manager View clicked", async () => {
    usePathnameMock.mockReturnValue("/");
    render(<SideBar />);

    await userEvent.click(screen.getByRole("button", { name: /Manager View/i }));

    expect(mockPush).toHaveBeenCalledWith("/dashboard/manager");
  });

  test("navigates to technician dashboard when Technician View clicked", async () => {
    usePathnameMock.mockReturnValue("/");
    render(<SideBar />);

    await userEvent.click(screen.getByRole("button", { name: /Technician View/i }));

    expect(mockPush).toHaveBeenCalledWith("/dashboard/technician");
  });

  test("applies active styles to Manager View button when role is manager", () => {
    usePathnameMock.mockReturnValue("/dashboard/manager/m-1");
    useNotificationMock.mockReturnValue({ notifications: [], markRead: jest.fn() });
    render(<SideBar />);

    const btn = screen.getByRole("button", { name: /Manager View/i });
    expect(btn).toHaveClass("bg-ink", "text-white");
  });

  test("applies active styles to Technician View button when role is technician", () => {
    usePathnameMock.mockReturnValue("/dashboard/technician/t-1");
    useNotificationMock.mockReturnValue({ notifications: [], markRead: jest.fn() });
    render(<SideBar />);

    const btn = screen.getByRole("button", { name: /Technician View/i });
    expect(btn).toHaveClass("bg-ink", "text-white");
  });
});
