"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getArch = getArch;
exports.getPlatform = getPlatform;
function getArch() {
    const arch = process.arch;
    const archMapping = {
        arm64: "arm64",
        ia32: "i386",
        x64: "x86_64",
    };
    if (arch in archMapping) {
        return archMapping[arch];
    }
}
function getPlatform() {
    const platform = process.platform;
    const platformMapping = {
        darwin: "Darwin",
        linux: "Linux",
        win32: "Windows",
    };
    if (platform in platformMapping) {
        return platformMapping[platform];
    }
}
