import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import TechnicianRow from "./TechnicianRow";
import { getTechnicianJobs } from "@/api/jobs";
import type { Job } from "@/api/jobs";
import type { Technician } from "@/api/technicians";

jest.mock("@/api/jobs");
const getTechnicianJobsMock = jest.mocked(getTechnicianJobs);

const technician: Technician = {
  ID: "t-1",
  Name: "Jane Doe",
  Email: "jane@example.com",
};

const baseProps = {
  technician,
  date: "2026-05-03",
  selected: false,
  onClick: jest.fn(),
  onTimeSelect: jest.fn(),
  selectedTime: "",
};

const makeJob = (id: string, startsAt: string, endsAt: string): Job => ({
  id,
  quote_id: "q-1",
  technician_id: "t-1",
  manager_id: "m-1",
  starts_at: startsAt,
  ends_at: endsAt,
  status: "scheduled",
  completed_at: null,
  created_at: "2026-05-01T00:00:00Z",
});

describe("TechnicianRow", () => {
  afterEach(() => {
    jest.resetAllMocks();
  });

  test("renders technician name and email", async () => {
    getTechnicianJobsMock.mockResolvedValue([]);
    render(<TechnicianRow {...baseProps} />);

    expect(screen.getByText("Jane Doe")).toBeInTheDocument();
    expect(screen.getByText("jane@example.com")).toBeInTheDocument();
  });

  test("shows 0 jobs booked when no jobs", async () => {
    getTechnicianJobsMock.mockResolvedValue([]);
    render(<TechnicianRow {...baseProps} />);

    await waitFor(() => expect(screen.getByText("0 jobs booked")).toBeInTheDocument());
  });

  test("shows job count after fetch", async () => {
    const job = makeJob("j-1", "2026-05-03T07:00:00Z", "2026-05-03T09:00:00Z");
    getTechnicianJobsMock.mockResolvedValue([job]);
    render(<TechnicianRow {...baseProps} />);

    await waitFor(() => expect(screen.getByText("1 job booked")).toBeInTheDocument());
  });

  test("renders job time range for each booked job", async () => {
    const job = makeJob("j-1", "2026-05-03T07:00:00Z", "2026-05-03T09:00:00Z");
    getTechnicianJobsMock.mockResolvedValue([job]);
    render(<TechnicianRow {...baseProps} />);

    await waitFor(() => expect(screen.getByText(/07:00/)).toBeInTheDocument());
    expect(screen.getByText(/09:00/)).toBeInTheDocument();
  });

  test("calls onClick when row clicked", async () => {
    const onClick = jest.fn();
    getTechnicianJobsMock.mockResolvedValue([]);
    render(<TechnicianRow {...baseProps} onClick={onClick} />);

    await userEvent.click(screen.getByText("Jane Doe"));

    expect(onClick).toHaveBeenCalledTimes(1);
  });

  test("calls onTimeSelect when time option selected", async () => {
    const onTimeSelect = jest.fn();
    getTechnicianJobsMock.mockResolvedValue([]);
    render(<TechnicianRow {...baseProps} onTimeSelect={onTimeSelect} />);

    await waitFor(() => expect(screen.getByRole("combobox")).toBeInTheDocument());

    await userEvent.selectOptions(screen.getByRole("combobox"), "07:00");

    expect(onTimeSelect).toHaveBeenCalledWith("07:00");
  });

  test("excludes blocked slots from available times", async () => {
    const job = makeJob("j-1", "2026-05-03T07:00:00Z", "2026-05-03T09:00:00Z");
    getTechnicianJobsMock.mockResolvedValue([job]);
    render(<TechnicianRow {...baseProps} />);

    await waitFor(() => expect(screen.getByText("1 job booked")).toBeInTheDocument());

    const options = screen.getAllByRole("option").map((o) => o.textContent);
    expect(options).not.toContain("07:00");
    expect(options).not.toContain("07:30");
  });

  test("renders bg-ok dot for completed job", async () => {
    const job: Job = {
      ...makeJob("j-1", "2026-05-03T07:00:00Z", "2026-05-03T09:00:00Z"),
      status: "completed",
    };
    getTechnicianJobsMock.mockResolvedValue([job]);
    render(<TechnicianRow {...baseProps} selected={true} />);

    await waitFor(() => expect(screen.getByText("1 job booked")).toBeInTheDocument());
    await userEvent.click(screen.getByText("Jane Doe"));

    const dot = document.querySelector(".bg-ok");
    expect(dot).toBeInTheDocument();
  });

  test("renders bg-accent dot for scheduled job", async () => {
    const job = makeJob("j-1", "2026-05-03T07:00:00Z", "2026-05-03T09:00:00Z");
    getTechnicianJobsMock.mockResolvedValue([job]);
    render(<TechnicianRow {...baseProps} selected={true} />);

    await waitFor(() => expect(screen.getByText("1 job booked")).toBeInTheDocument());
    await userEvent.click(screen.getByText("Jane Doe"));

    const dot = document.querySelector(".bg-accent");
    expect(dot).toBeInTheDocument();
  });
});
