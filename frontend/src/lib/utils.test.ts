import { describe, it, expect } from "vitest";
import { cn } from "./utils";

describe("cn", () => {
  it("merges class names", () => {
    expect(cn("a", "b")).toBe("a b");
  });

  it("resolves tailwind conflicts (later wins)", () => {
    expect(cn("px-2", "px-4")).toBe("px-4");
  });
});
