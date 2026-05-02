import { render, screen, fireEvent } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import TechnicianList from "./TechnicianList";
import { getTechnicianJobs } from "@/api/jobs";
import type { Technician } from "@/api/technicians";

jest.mock("@/api/jobs");
jest.mocked(getTechnicianJobs).mockResolvedValue([]);

const makeTechnician = (id: string, name: string): Technician => ({
  ID: id,
  Name: name,
  Email: `${name.toLowerCase()}@example.com`,
});

const baseProps = {
  technicians: [],
  quoteSelected: false,
  date: "2026-05-03",
  onDateChange: jest.fn(),
  selectedTechnicianId: null,
  selectedTime: "",
  onSelect: jest.fn(),
  onTimeSelect: jest.fn(),
  page: 1,
  totalPages: 1,
  onPrev: jest.fn(),
  onNext: jest.fn(),
};

describe("TechnicianList", () => {
  afterEach(() => {
    jest.resetAllMocks();
    jest.mocked(getTechnicianJobs).mockResolvedValue([]);
  });

  test("shows prompt when no quote selected", () => {
    render(<TechnicianList {...baseProps} quoteSelected={false} />);

    expect(screen.getByText("Select a quote first.")).toBeInTheDocument();
  });

  test("does not show prompt when quote selected", () => {
    render(<TechnicianList {...baseProps} quoteSelected={true} />);

    expect(screen.queryByText("Select a quote first.")).not.toBeInTheDocument();
  });

  test("renders a row for each technician when quote selected", () => {
    const technicians = [makeTechnician("t-1", "Alice"), makeTechnician("t-2", "Bob")];
    render(<TechnicianList {...baseProps} quoteSelected={true} technicians={technicians} />);

    expect(screen.getByText("Alice")).toBeInTheDocument();
    expect(screen.getByText("Bob")).toBeInTheDocument();
  });

  test("does not render technicians when quote not selected", () => {
    const technicians = [makeTechnician("t-1", "Alice")];
    render(<TechnicianList {...baseProps} quoteSelected={false} technicians={technicians} />);

    expect(screen.queryByText("Alice")).not.toBeInTheDocument();
  });

  test("calls onDateChange when date input changes", async () => {
    const onDateChange = jest.fn();
    render(<TechnicianList {...baseProps} onDateChange={onDateChange} />);

    fireEvent.change(screen.getByDisplayValue("2026-05-03"), { target: { value: "2026-05-04" } });

    expect(onDateChange).toHaveBeenCalledWith("2026-05-04");
  });

  test("does not render Pagination when totalPages <= 1", () => {
    render(<TechnicianList {...baseProps} totalPages={1} />);

    expect(screen.queryByRole("button", { name: "Prev" })).not.toBeInTheDocument();
  });

  test("renders Pagination when totalPages > 1", () => {
    render(<TechnicianList {...baseProps} totalPages={3} page={2} />);

    expect(screen.getByRole("button", { name: "Prev" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Next" })).toBeInTheDocument();
  });
});
