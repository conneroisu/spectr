import assert from "node:assert/strict";
import { describe, it } from "node:test";
import { getArch, getPlatform } from "../../../src/utils/platforms";

describe("Platform Utilities", () => {
  describe("getArch()", () => {
    it("should map arm64 to arm64", () => {
      // Mock process.arch
      const originalArch = process.arch;
      Object.defineProperty(process, "arch", {
        configurable: true,
        value: "arm64",
      });

      const result = getArch();
      assert.equal(result, "arm64");

      // Restore original
      Object.defineProperty(process, "arch", {
        configurable: true,
        value: originalArch,
      });
    });

    it("should map ia32 to i386", () => {
      const originalArch = process.arch;
      Object.defineProperty(process, "arch", {
        configurable: true,
        value: "ia32",
      });

      const result = getArch();
      assert.equal(result, "i386");

      Object.defineProperty(process, "arch", {
        configurable: true,
        value: originalArch,
      });
    });

    it("should map x64 to x86_64", () => {
      const originalArch = process.arch;
      Object.defineProperty(process, "arch", {
        configurable: true,
        value: "x64",
      });

      const result = getArch();
      assert.equal(result, "x86_64");

      Object.defineProperty(process, "arch", {
        configurable: true,
        value: originalArch,
      });
    });

    it("should return undefined for unsupported architecture", () => {
      const originalArch = process.arch;
      Object.defineProperty(process, "arch", {
        configurable: true,
        value: "unsupported",
      });

      const result = getArch();
      assert.equal(result, undefined);

      Object.defineProperty(process, "arch", {
        configurable: true,
        value: originalArch,
      });
    });
  });

  describe("getPlatform()", () => {
    it("should map darwin to Darwin", () => {
      const originalPlatform = process.platform;
      Object.defineProperty(process, "platform", {
        configurable: true,
        value: "darwin",
      });

      const result = getPlatform();
      assert.equal(result, "Darwin");

      Object.defineProperty(process, "platform", {
        configurable: true,
        value: originalPlatform,
      });
    });

    it("should map linux to Linux", () => {
      const originalPlatform = process.platform;
      Object.defineProperty(process, "platform", {
        configurable: true,
        value: "linux",
      });

      const result = getPlatform();
      assert.equal(result, "Linux");

      Object.defineProperty(process, "platform", {
        configurable: true,
        value: originalPlatform,
      });
    });

    it("should map win32 to Windows", () => {
      const originalPlatform = process.platform;
      Object.defineProperty(process, "platform", {
        configurable: true,
        value: "win32",
      });

      const result = getPlatform();
      assert.equal(result, "Windows");

      Object.defineProperty(process, "platform", {
        configurable: true,
        value: originalPlatform,
      });
    });

    it("should return undefined for unsupported platform", () => {
      const originalPlatform = process.platform;
      Object.defineProperty(process, "platform", {
        configurable: true,
        value: "unsupported",
      });

      const result = getPlatform();
      assert.equal(result, undefined);

      Object.defineProperty(process, "platform", {
        configurable: true,
        value: originalPlatform,
      });
    });
  });
});
