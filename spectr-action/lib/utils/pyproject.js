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
exports.getRuffVersionFromRequirementsFile = getRuffVersionFromRequirementsFile;
const fs = __importStar(require("node:fs"));
const core = __importStar(require("@actions/core"));
const toml = __importStar(require("smol-toml"));
function getRuffVersionFromAllDependencies(allDependencies) {
    const ruffVersionDefinition = allDependencies.find((dep) => dep.startsWith("ruff"));
    if (ruffVersionDefinition) {
        const ruffVersion = ruffVersionDefinition
            .match(/^ruff([^A-Z0-9._-]+.*)$/)?.[1]
            .trim();
        if (ruffVersion?.startsWith("==")) {
            return ruffVersion.slice(2);
        }
        core.info(`Found ruff version in pyproject.toml: ${ruffVersion}`);
        return ruffVersion;
    }
    return undefined;
}
function parsePyproject(pyprojectContent) {
    const pyproject = toml.parse(pyprojectContent);
    const dependencies = pyproject?.project?.dependencies || [];
    const optionalDependencies = Object.values(pyproject?.project?.["optional-dependencies"] || {}).flat();
    const devDependencies = Object.values(pyproject?.["dependency-groups"] || {})
        .flat()
        .filter((item) => typeof item === "string");
    return (getRuffVersionFromAllDependencies(dependencies.concat(optionalDependencies, devDependencies)) || getRuffVersionFromPoetryGroups(pyproject));
}
function getRuffVersionFromPoetryGroups(pyproject) {
    // Special handling for Poetry until it supports PEP 735
    // See: <https://github.com/python-poetry/poetry/issues/9751>
    const poetry = pyproject?.tool?.poetry || {};
    const poetryGroups = Object.values(poetry.group || {});
    if (poetry.dependencies) {
        poetryGroups.unshift({ dependencies: poetry.dependencies });
    }
    return poetryGroups
        .flatMap((group) => Object.entries(group.dependencies))
        .map(([name, spec]) => {
        if (name === "ruff" && typeof spec === "string")
            return spec;
        return undefined;
    })
        .find((version) => version !== undefined);
}
function getRuffVersionFromRequirementsFile(filePath) {
    if (!fs.existsSync(filePath)) {
        core.warning(`Could not find file: ${filePath}`);
        return undefined;
    }
    const pyprojectContent = fs.readFileSync(filePath, "utf-8");
    if (filePath.endsWith(".txt")) {
        return getRuffVersionFromAllDependencies(pyprojectContent.split("\n"));
    }
    try {
        return parsePyproject(pyprojectContent);
    }
    catch (err) {
        const message = err.message;
        core.warning(`Error while parsing ${filePath}: ${message}`);
        return undefined;
    }
}
