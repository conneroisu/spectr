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
const core = __importStar(require("@actions/core"));
const core_1 = require("@octokit/core");
const plugin_paginate_rest_1 = require("@octokit/plugin-paginate-rest");
const plugin_rest_endpoint_methods_1 = require("@octokit/plugin-rest-endpoint-methods");
const semver = __importStar(require("semver"));
const update_known_checksums_1 = require("./download/checksum/update-known-checksums");
const constants_1 = require("./utils/constants");
const PaginatingOctokit = core_1.Octokit.plugin(plugin_paginate_rest_1.paginateRest, plugin_rest_endpoint_methods_1.restEndpointMethods);
async function run() {
    const checksumFilePath = process.argv.slice(2)[0];
    const github_token = process.argv.slice(2)[1];
    const octokit = new PaginatingOctokit({ auth: github_token });
    const response = await octokit.paginate(octokit.rest.repos.listReleases, {
        owner: constants_1.OWNER,
        repo: constants_1.REPO,
    });
    const downloadUrls = response.flatMap((release) => release.assets
        .filter((asset) => asset.name.endsWith(".sha256"))
        .map((asset) => asset.browser_download_url));
    await (0, update_known_checksums_1.updateChecksums)(checksumFilePath, downloadUrls);
    const latestVersion = response
        .map((release) => release.tag_name)
        .sort(semver.rcompare)[0];
    core.setOutput("latest-version", latestVersion);
}
run();
