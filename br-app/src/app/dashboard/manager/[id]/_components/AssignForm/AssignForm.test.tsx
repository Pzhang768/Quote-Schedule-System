import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import AssignForm from "./AssignForm";
import { assignJob } from "@/api/jobs";
import type { Quote } from "@/api/quotes";
import type { Technician } from "@/api/technicians";

jest.mock("@/api/jobs");
const assignJobMock = jest.mocked(assignJob);

const quote: Quote = {
  ID: "q-1",
  CustomerName: "John Smith",
  Address: "123 Main St",
  Description: "",
  Status: "unscheduled",
  CreatedAt: "2026-05-01T00:00:00Z",
  UpdatedAt: "2026-05-01T00:00:00Z",
};

const technician: Technician = {
  ID: "t-1",
  Name: "Jane Doe",
  Email: "jane@example.com",
};

const baseProps = {
  managerId: "m-1",
  quote,
  technician,
  date: "2026-05-03",
  time: "09:00",
  onSuccess: jest.fn(),
};

describe("AssignForm", () => {
  afterEach(() => {
    jest.resetAllMocks();
  });

  test("renders quote, technician, and time details", () => {
    render(<AssignForm {...baseProps} />);

    expect(screen.getByText("John Smith")).toBeInTheDocument();
    expect(screen.getByText("123 Main St")).toBeInTheDocument();
    expect(screen.getByText("Jane Doe")).toBeInTheDocument();
    expect(screen.getByText("2026-05-03 at 09:00")).toBeInTheDocument();
  });

  test("renders Assign Job button", () => {
    render(<AssignForm {...baseProps} />);

    expect(screen.getByRole("button", { name: "Assign Job" })).toBeInTheDocument();
  });

  test("calls assignJob with correct payload on submit", async () => {
    assignJobMock.mockResolvedValue({ id: "j-1" } as never);
    render(<AssignForm {...baseProps} />);

    await userEvent.click(screen.getByRole("button", { name: "Assign Job" }));

    expect(assignJobMock).toHaveBeenCalledWith({
      quote_id: "q-1",
      technician_id: "t-1",
      manager_id: "m-1",
      starts_at: "2026-05-03T09:00:00Z",
    });
  });

  test("calls onSuccess after successful submit", async () => {
    const onSuccess = jest.fn();
    assignJobMock.mockResolvedValue({ id: "j-1" } as never);
    render(<AssignForm {...baseProps} onSuccess={onSuccess} />);

    await userEvent.click(screen.getByRole("button", { name: "Assign Job" }));

    await waitFor(() => expect(onSuccess).toHaveBeenCalledTimes(1));
  });

  test("shows conflict error on 409", async () => {
    assignJobMock.mockRejectedValue({ response: { status: 409 } });
    render(<AssignForm {...baseProps} />);

    await userEvent.click(screen.getByRole("button", { name: "Assign Job" }));

    await waitFor(() =>
      expect(
        screen.getByText("This technician has a conflicting job at that time.")
      ).toBeInTheDocument()
    );
  });

  test("shows generic error on other failures", async () => {
    assignJobMock.mockRejectedValue({ response: { status: 500, data: { error: "Server error" } } });
    render(<AssignForm {...baseProps} />);

    await userEvent.click(screen.getByRole("button", { name: "Assign Job" }));

    await waitFor(() => expect(screen.getByText("Server error")).toBeInTheDocument());
  });

  test("shows fallback error when no error message", async () => {
    assignJobMock.mockRejectedValue({});
    render(<AssignForm {...baseProps} />);

    await userEvent.click(screen.getByRole("button", { name: "Assign Job" }));

    await waitFor(() => expect(screen.getByText("Something went wrong.")).toBeInTheDocument());
  });

  test("disables button while submitting", async () => {
    let resolve: () => void;
    assignJobMock.mockReturnValue(
      new Promise((res) => {
        resolve = () => res({} as never);
      })
    );
    render(<AssignForm {...baseProps} />);

    await userEvent.click(screen.getByRole("button", { name: "Assign Job" }));

    expect(screen.getByRole("button", { name: "Assigning..." })).toBeDisabled();
    resolve!();
  });
});
