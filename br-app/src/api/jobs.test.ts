import api from "./index";
import { assignJob, completeJob, getTechnicianJobs } from "./jobs";

jest.mock("./index");

const apiMock = jest.mocked(api);

const createMockApi = () => ({
  post: jest.fn(),
  patch: jest.fn(),
  get: jest.fn(),
});

describe("jobs api", () => {
  let mock: ReturnType<typeof createMockApi>;

  beforeEach(() => {
    mock = createMockApi();
    (apiMock.post as jest.Mock) = mock.post;
    (apiMock.patch as jest.Mock) = mock.patch;
    (apiMock.get as jest.Mock) = mock.get;
  });

  afterEach(() => {
    jest.resetAllMocks();
  });

  test("assignJob: posts payload and returns job", async () => {
    const job = { id: "job-1", status: "scheduled" };
    mock.post.mockResolvedValue({ data: { data: job } });

    const payload = {
      quote_id: "q-1",
      technician_id: "t-1",
      manager_id: "m-1",
      starts_at: "2026-05-03T07:00:00Z",
    };
    const result = await assignJob(payload);

    expect(mock.post).toHaveBeenCalledWith("/jobs", payload);
    expect(result).toEqual(job);
  });

  test("completeJob: patches job with technician id and returns updated job", async () => {
    const job = { id: "job-1", status: "completed" };
    mock.patch.mockResolvedValue({ data: { data: job } });

    const result = await completeJob("job-1", "t-1");

    expect(mock.patch).toHaveBeenCalledWith("/jobs/job-1/complete", {
      technician_id: "t-1",
    });
    expect(result).toEqual(job);
  });

  test("getTechnicianJobs: fetches jobs for technician on given date", async () => {
    const jobs = [{ id: "job-1", status: "scheduled" }];
    mock.get.mockResolvedValue({ data: { data: jobs } });

    const result = await getTechnicianJobs("t-1", "2026-05-03");

    expect(mock.get).toHaveBeenCalledWith("/technicians/t-1/jobs", {
      params: { date: "2026-05-03" },
    });
    expect(result).toEqual(jobs);
  });
});
