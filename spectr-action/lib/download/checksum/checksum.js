"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.validateChecksum = validateChecksum;
exports.isknownVersion = isknownVersion;
const crypto = __importStar(require("node:crypto"));
const fs = __importStar(require("node:fs"));
const core = __importStar(require("@actions/core"));
const known_checksums_1 = require("./known-checksums");
async function validateChecksum(checkSum, downloadPath, arch, platform, version) {
    let isValid;
    if (checkSum !== undefined && checkSum !== "") {
        isValid = await validateFileCheckSum(downloadPath, checkSum);
    }
    else {
        core.debug("Checksum not provided. Checking known checksums.");
        const key = `${arch}-${platform}-${version}`;
        if (key in known_checksums_1.KNOWN_CHECKSUMS) {
            const knownChecksum = known_checksums_1.KNOWN_CHECKSUMS[`${arch}-${platform}-${version}`];
            core.debug(`Checking checksum for ${arch}-${platform}-${version}.`);
            isValid = await validateFileCheckSum(downloadPath, knownChecksum);
        }
        else {
            core.debug(`No known checksum found for ${key}.`);
        }
    }
    if (isValid === false) {
        throw new Error(`Checksum for ${downloadPath} did not match ${checkSum}.`);
    }
    if (isValid === true) {
        core.debug(`Checksum for ${downloadPath} is valid.`);
    }
}
async function validateFileCheckSum(filePath, expected) {
    return new Promise((resolve, reject) => {
        const hash = crypto.createHash("sha256");
        const stream = fs.createReadStream(filePath);
        stream.on("error", (err) => reject(err));
        stream.on("data", (chunk) => hash.update(chunk));
        stream.on("end", () => {
            const actual = hash.digest("hex");
            resolve(actual === expected);
        });
    });
}
function isknownVersion(version) {
    const pattern = new RegExp(`^.*-.*-${version}$`);
    return Object.keys(known_checksums_1.KNOWN_CHECKSUMS).some((key) => pattern.test(key));
}
