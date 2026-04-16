import { getPublicApiV1BaseUrl } from "../api-base-url";

describe("getPublicApiV1BaseUrl", () => {
  const original = process.env.NEXT_PUBLIC_API_URL;

  afterEach(() => {
    if (original === undefined) {
      delete process.env.NEXT_PUBLIC_API_URL;
    } else {
      process.env.NEXT_PUBLIC_API_URL = original;
    }
  });

  it("appends /api/v1 when env is unset (defaults to localhost:8080)", () => {
    delete process.env.NEXT_PUBLIC_API_URL;
    expect(getPublicApiV1BaseUrl()).toBe("http://localhost:8080/api/v1");
  });

  it("does not duplicate when env ends with /api/v1", () => {
    process.env.NEXT_PUBLIC_API_URL = "http://localhost:8080/api/v1";
    expect(getPublicApiV1BaseUrl()).toBe("http://localhost:8080/api/v1");
  });

  it("collapses repeated /api/v1 suffix", () => {
    process.env.NEXT_PUBLIC_API_URL = "http://localhost:8080/api/v1/api/v1";
    expect(getPublicApiV1BaseUrl()).toBe("http://localhost:8080/api/v1");
  });

  it("strips trailing slashes before appending", () => {
    process.env.NEXT_PUBLIC_API_URL = "http://api.example.com///";
    expect(getPublicApiV1BaseUrl()).toBe("http://api.example.com/api/v1");
  });
});
