import assert from "node:assert/strict";
import { afterEach, describe, it } from "node:test";
import { getArch, getPlatform } from "../../../src/utils/platforms";
import { mockArch, mockPlatform } from "../../helpers/test-utils";

describe("Platform Utilities (Improved)", () => {
  describe("getArch()", () => {
    let restore: (() => void) | undefined;

    afterEach(() => {
      restore?.();
    });

    it("should map arm64 to arm64", () => {
      restore = mockArch("arm64");
      const result = getArch();
      assert.equal(result, "arm64");
    });

    it("should map ia32 to i386", () => {
      restore = mockArch("ia32");
      const result = getArch();
      assert.equal(result, "i386");
    });

    it("should map x64 to x86_64", () => {
      restore = mockArch("x64");
      const result = getArch();
      assert.equal(result, "x86_64");
    });

    it("should return undefined for unsupported architecture", () => {
      restore = mockArch("unsupported");
      const result = getArch();
      assert.equal(result, undefined);
    });
  });

  describe("getPlatform()", () => {
    let restore: (() => void) | undefined;

    afterEach(() => {
      restore?.();
    });

    it("should map darwin to Darwin", () => {
      restore = mockPlatform("darwin");
      const result = getPlatform();
      assert.equal(result, "Darwin");
    });

    it("should map linux to Linux", () => {
      restore = mockPlatform("linux");
      const result = getPlatform();
      assert.equal(result, "Linux");
    });

    it("should map win32 to Windows", () => {
      restore = mockPlatform("win32");
      const result = getPlatform();
      assert.equal(result, "Windows");
    });

    it("should return undefined for unsupported platform", () => {
      restore = mockPlatform("unsupported");
      const result = getPlatform();
      assert.equal(result, undefined);
    });
  });
});
