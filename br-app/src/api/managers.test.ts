import api from "./index";
import { getManagers } from "./managers";

jest.mock("./index");

const apiMock = jest.mocked(api);

describe("managers api", () => {
  afterEach(() => {
    jest.resetAllMocks();
  });

  test("getManagers: fetches paginated managers", async () => {
    const response = { data: [], total: 0, page: 1, page_size: 5, total_pages: 1 };
    (apiMock.get as jest.Mock) = jest.fn().mockResolvedValue({ data: response });

    const result = await getManagers(1, 5);

    expect(apiMock.get).toHaveBeenCalledWith("/managers", {
      params: { page: 1, page_size: 5 },
    });
    expect(result).toEqual(response);
  });
});
