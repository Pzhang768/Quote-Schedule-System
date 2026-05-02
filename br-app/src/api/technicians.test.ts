import api from "./index";
import { getTechnicians } from "./technicians";

jest.mock("./index");

const apiMock = jest.mocked(api);

describe("technicians api", () => {
  afterEach(() => {
    jest.resetAllMocks();
  });

  test("getTechnicians: fetches paginated technicians", async () => {
    const response = { data: [], total: 0, page: 1, page_size: 5, total_pages: 1 };
    (apiMock.get as jest.Mock) = jest.fn().mockResolvedValue({ data: response });

    const result = await getTechnicians(1, 5);

    expect(apiMock.get).toHaveBeenCalledWith("/technicians", {
      params: { page: 1, page_size: 5 },
    });
    expect(result).toEqual(response);
  });
});
