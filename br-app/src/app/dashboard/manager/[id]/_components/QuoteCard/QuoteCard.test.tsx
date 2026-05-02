import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import QuoteCard from "./QuoteCard";
import type { Quote } from "@/api/quotes";

const quote: Quote = {
  ID: "q-1",
  CustomerName: "John Smith",
  Address: "123 Main St",
  Description: "Leaky tap",
  Status: "unscheduled",
  CreatedAt: "2026-05-01T00:00:00Z",
  UpdatedAt: "2026-05-01T00:00:00Z",
};

describe("QuoteCard", () => {
  test("renders customer name and address", () => {
    render(<QuoteCard quote={quote} selected={false} onClick={jest.fn()} />);

    expect(screen.getByText("John Smith")).toBeInTheDocument();
    expect(screen.getByText("123 Main St")).toBeInTheDocument();
  });

  test("renders description when present", () => {
    render(<QuoteCard quote={quote} selected={false} onClick={jest.fn()} />);

    expect(screen.getByText("Leaky tap")).toBeInTheDocument();
  });

  test("does not render description when absent", () => {
    const q = { ...quote, Description: "" };
    render(<QuoteCard quote={q} selected={false} onClick={jest.fn()} />);

    expect(screen.queryByText("Leaky tap")).not.toBeInTheDocument();
  });

  test("calls onClick when clicked", async () => {
    const onClick = jest.fn();
    render(<QuoteCard quote={quote} selected={false} onClick={onClick} />);

    await userEvent.click(screen.getByText("John Smith"));

    expect(onClick).toHaveBeenCalledTimes(1);
  });

  test("applies selected styles when selected is true", () => {
    const { container } = render(<QuoteCard quote={quote} selected={true} onClick={jest.fn()} />);

    expect(container.firstChild).toHaveClass("border-ink", "bg-ink/5");
  });
});
