import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Pagination from "./Pagination";

describe("Pagination", () => {
  test("renders nothing when totalPages <= 1", () => {
    const { container } = render(
      <Pagination page={1} totalPages={1} onPrev={jest.fn()} onNext={jest.fn()} />
    );

    expect(container).toBeEmptyDOMElement();
  });

  test("renders prev and next buttons when totalPages > 1", () => {
    render(<Pagination page={2} totalPages={3} onPrev={jest.fn()} onNext={jest.fn()} />);

    expect(screen.getByRole("button", { name: "Prev" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Next" })).toBeInTheDocument();
  });

  test("renders current page and total pages", () => {
    render(<Pagination page={2} totalPages={5} onPrev={jest.fn()} onNext={jest.fn()} />);

    expect(screen.getByText("2 / 5")).toBeInTheDocument();
  });

  test("disables Prev on first page", () => {
    render(<Pagination page={1} totalPages={3} onPrev={jest.fn()} onNext={jest.fn()} />);

    expect(screen.getByRole("button", { name: "Prev" })).toBeDisabled();
  });

  test("disables Next on last page", () => {
    render(<Pagination page={3} totalPages={3} onPrev={jest.fn()} onNext={jest.fn()} />);

    expect(screen.getByRole("button", { name: "Next" })).toBeDisabled();
  });

  test("calls onPrev when Prev clicked", async () => {
    const onPrev = jest.fn();
    render(<Pagination page={2} totalPages={3} onPrev={onPrev} onNext={jest.fn()} />);

    await userEvent.click(screen.getByRole("button", { name: "Prev" }));

    expect(onPrev).toHaveBeenCalledTimes(1);
  });

  test("calls onNext when Next clicked", async () => {
    const onNext = jest.fn();
    render(<Pagination page={1} totalPages={3} onPrev={jest.fn()} onNext={onNext} />);

    await userEvent.click(screen.getByRole("button", { name: "Next" }));

    expect(onNext).toHaveBeenCalledTimes(1);
  });
});
