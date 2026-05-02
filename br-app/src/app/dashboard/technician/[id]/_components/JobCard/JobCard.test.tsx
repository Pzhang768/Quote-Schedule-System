import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import JobCard from "./JobCard";
import { completeJob } from "@/api/jobs";
import type { Job } from "@/api/jobs";

jest.mock("@/api/jobs");
const completeJobMock = jest.mocked(completeJob);

const makeJob = (status: Job["status"], extra: Partial<Job> = {}): Job => ({
  id: "j-1",
  quote_id: "q-1",
  technician_id: "t-1",
  manager_id: "m-1",
  starts_at: "2026-05-03T07:00:00Z",
  ends_at: "2026-05-03T09:00:00Z",
  status,
  completed_at: status === "completed" ? "2026-05-03T09:00:00Z" : null,
  created_at: "2026-05-01T00:00:00Z",
  ...extra,
});

describe("JobCard", () => {
  afterEach(() => {
    jest.resetAllMocks();
  });

  test("renders time range and status", () => {
    render(<JobCard job={makeJob("scheduled")} technicianId="t-1" onCompleted={jest.fn()} />);

    expect(screen.getByText(/07:00/)).toBeInTheDocument();
    expect(screen.getByText(/09:00/)).toBeInTheDocument();
    expect(screen.getByText("scheduled")).toBeInTheDocument();
  });

  test("renders customer name and address when present", () => {
    render(
      <JobCard
        job={makeJob("scheduled", { customer_name: "John Smith", address: "123 Main St" })}
        technicianId="t-1"
        onCompleted={jest.fn()}
      />
    );

    expect(screen.getByText("John Smith")).toBeInTheDocument();
    expect(screen.getByText("123 Main St")).toBeInTheDocument();
  });

  test("does not render customer name when absent", () => {
    render(<JobCard job={makeJob("scheduled")} technicianId="t-1" onCompleted={jest.fn()} />);

    expect(screen.queryByText("John Smith")).not.toBeInTheDocument();
  });

  test("renders Complete button for scheduled job", () => {
    render(<JobCard job={makeJob("scheduled")} technicianId="t-1" onCompleted={jest.fn()} />);

    expect(screen.getByRole("button", { name: "Complete" })).toBeInTheDocument();
  });

  test("does not render Complete button for completed job", () => {
    render(<JobCard job={makeJob("completed")} technicianId="t-1" onCompleted={jest.fn()} />);

    expect(screen.queryByRole("button", { name: "Complete" })).not.toBeInTheDocument();
  });

  test("renders Done label for completed job", () => {
    render(<JobCard job={makeJob("completed")} technicianId="t-1" onCompleted={jest.fn()} />);

    expect(screen.getByText("Done")).toBeInTheDocument();
  });

  test("calls completeJob with correct args on click", async () => {
    completeJobMock.mockResolvedValue(makeJob("completed"));
    render(<JobCard job={makeJob("scheduled")} technicianId="t-1" onCompleted={jest.fn()} />);

    await userEvent.click(screen.getByRole("button", { name: "Complete" }));

    expect(completeJobMock).toHaveBeenCalledWith("j-1", "t-1");
  });

  test("calls onCompleted after successful complete", async () => {
    const onCompleted = jest.fn();
    completeJobMock.mockResolvedValue(makeJob("completed"));
    render(<JobCard job={makeJob("scheduled")} technicianId="t-1" onCompleted={onCompleted} />);

    await userEvent.click(screen.getByRole("button", { name: "Complete" }));

    await waitFor(() => expect(onCompleted).toHaveBeenCalledTimes(1));
  });

  test("shows error message on failure", async () => {
    completeJobMock.mockRejectedValue(new Error("network error"));
    render(<JobCard job={makeJob("scheduled")} technicianId="t-1" onCompleted={jest.fn()} />);

    await userEvent.click(screen.getByRole("button", { name: "Complete" }));

    await waitFor(() => expect(screen.getByText("Failed to complete job.")).toBeInTheDocument());
  });

  test("disables button while completing", async () => {
    let resolve: () => void;
    completeJobMock.mockReturnValue(
      new Promise((res) => {
        resolve = () => res(makeJob("completed"));
      })
    );
    render(<JobCard job={makeJob("scheduled")} technicianId="t-1" onCompleted={jest.fn()} />);

    await userEvent.click(screen.getByRole("button", { name: "Complete" }));

    expect(screen.getByRole("button", { name: "Completing..." })).toBeDisabled();
    resolve!();
  });
});
