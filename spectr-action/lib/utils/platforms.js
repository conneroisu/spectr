"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getArch = getArch;
exports.getPlatform = getPlatform;
function getArch() {
    const arch = process.arch;
    const archMapping = {
        arm64: "aarch64",
        ia32: "i686",
        x64: "x86_64",
    };
    if (arch in archMapping) {
        return archMapping[arch];
    }
}
function getPlatform() {
    const platform = process.platform;
    const platformMapping = {
        darwin: "apple-darwin",
        linux: "unknown-linux-gnu",
        win32: "pc-windows-msvc",
    };
    if (platform in platformMapping) {
        return platformMapping[platform];
    }
}
