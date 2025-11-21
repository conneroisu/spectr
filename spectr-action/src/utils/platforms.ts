export type Platform = "Linux" | "Darwin" | "Windows";
export type Architecture = "i386" | "x86_64" | "arm64";

export function getArch(): Architecture | undefined {
  const arch = process.arch;
  const archMapping: { [key: string]: Architecture } = {
    arm64: "arm64",
    ia32: "i386",
    x64: "x86_64",
  };

  if (arch in archMapping) {
    return archMapping[arch];
  }
}

export function getPlatform(): Platform | undefined {
  const platform = process.platform;
  const platformMapping: { [key: string]: Platform } = {
    darwin: "Darwin",
    linux: "Linux",
    win32: "Windows",
  };

  if (platform in platformMapping) {
    return platformMapping[platform];
  }
}
