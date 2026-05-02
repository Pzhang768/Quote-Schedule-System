import api from "./index";
import { getQuotes } from "./quotes";

jest.mock("./index");

const apiMock = jest.mocked(api);

describe("quotes api", () => {
  afterEach(() => {
    jest.resetAllMocks();
  });

  test("getQuotes: fetches paginated quotes without status filter", async () => {
    const response = { data: [], total: 0, page: 1, page_size: 5, total_pages: 1 };
    (apiMock.get as jest.Mock) = jest.fn().mockResolvedValue({ data: response });

    const result = await getQuotes(1, 5);

    expect(apiMock.get).toHaveBeenCalledWith("/quotes", {
      params: { page: 1, page_size: 5 },
    });
    expect(result).toEqual(response);
  });

  test("getQuotes: fetches paginated quotes with status filter", async () => {
    const response = { data: [], total: 0, page: 1, page_size: 5, total_pages: 1 };
    (apiMock.get as jest.Mock) = jest.fn().mockResolvedValue({ data: response });

    await getQuotes(1, 5, "unscheduled");

    expect(apiMock.get).toHaveBeenCalledWith("/quotes", {
      params: { page: 1, page_size: 5, status: "unscheduled" },
    });
  });
});
