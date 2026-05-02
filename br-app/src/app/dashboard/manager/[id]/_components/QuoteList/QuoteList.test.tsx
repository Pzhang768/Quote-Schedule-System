import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import QuoteList from "./QuoteList";
import type { Quote } from "@/api/quotes";

const makeQuote = (id: string, name: string): Quote => ({
  ID: id,
  CustomerName: name,
  Address: "1 St",
  Description: "",
  Status: "unscheduled",
  CreatedAt: "2026-05-01T00:00:00Z",
  UpdatedAt: "2026-05-01T00:00:00Z",
});

describe("QuoteList", () => {
  const baseProps = {
    selectedQuote: null,
    onSelect: jest.fn(),
    page: 1,
    totalPages: 1,
    onPrev: jest.fn(),
    onNext: jest.fn(),
  };

  test("renders empty state when no quotes", () => {
    render(<QuoteList {...baseProps} quotes={[]} />);

    expect(screen.getByText("No unscheduled quotes.")).toBeInTheDocument();
  });

  test("renders a card for each quote", () => {
    const quotes = [makeQuote("q-1", "Alice"), makeQuote("q-2", "Bob")];
    render(<QuoteList {...baseProps} quotes={quotes} />);

    expect(screen.getByText("Alice")).toBeInTheDocument();
    expect(screen.getByText("Bob")).toBeInTheDocument();
  });

  test("calls onSelect with the clicked quote", async () => {
    const onSelect = jest.fn();
    const quotes = [makeQuote("q-1", "Alice")];
    render(<QuoteList {...baseProps} quotes={quotes} onSelect={onSelect} />);

    await userEvent.click(screen.getByText("Alice"));

    expect(onSelect).toHaveBeenCalledWith(quotes[0]);
  });

  test("does not render Pagination when totalPages <= 1", () => {
    render(<QuoteList {...baseProps} quotes={[]} totalPages={1} />);

    expect(screen.queryByRole("button", { name: "Prev" })).not.toBeInTheDocument();
  });

  test("renders Pagination when totalPages > 1", () => {
    render(<QuoteList {...baseProps} quotes={[]} totalPages={3} page={2} />);

    expect(screen.getByRole("button", { name: "Prev" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Next" })).toBeInTheDocument();
  });
});
