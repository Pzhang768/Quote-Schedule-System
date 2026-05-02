import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import SideBar from "./SideBar";
import useNotification from "@/hooks/useNotification/useNotification";

jest.mock("next/navigation", () => ({
  useRouter: () => ({ push: jest.fn() }),
  usePathname: jest.fn(),
}));

jest.mock("@/hooks/useNotification/useNotification");
const useNotificationMock = jest.mocked(useNotification);

import { usePathname } from "next/navigation";
const usePathnameMock = jest.mocked(usePathname);

describe("SideBar", () => {
  afterEach(() => {
    jest.resetAllMocks();
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
});
