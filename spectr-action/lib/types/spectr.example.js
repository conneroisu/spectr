"use strict";
/**
 * Example usage of spectr validation types
 * This file demonstrates how to use the types in the action code
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.processValidationOutput = processValidationOutput;
exports.createGitHubAnnotations = createGitHubAnnotations;
exports.exitWithValidationStatus = exitWithValidationStatus;
exports.getChangeResults = getChangeResults;
exports.getSpecResults = getSpecResults;
exports.createSummary = createSummary;
const spectr_1 = require("./spectr");
/**
 * Example: Parse and process spectr validation output
 */
function processValidationOutput(jsonOutput) {
    // Parse JSON output from spectr command
    const validationOutput = JSON.parse(jsonOutput);
    // Check if all items are valid
    if ((0, spectr_1.allValid)(validationOutput)) {
        console.log("✓ All items passed validation");
        return;
    }
    // Check for any errors
    if ((0, spectr_1.hasAnyErrors)(validationOutput)) {
        console.error("✗ Validation failed with errors");
    }
    // Print summary
    const totalErrors = (0, spectr_1.getTotalErrorCount)(validationOutput);
    const totalWarnings = (0, spectr_1.getTotalWarningCount)(validationOutput);
    console.log(`Errors: ${totalErrors}, Warnings: ${totalWarnings}`);
    // Process each result
    for (const result of validationOutput) {
        console.log(`\nChecking: ${result.name} (${result.type})`);
        // Handle validation errors
        if ((0, spectr_1.hasError)(result)) {
            console.error(`  Error: ${result.error}`);
            continue;
        }
        // Handle validation report
        if ((0, spectr_1.hasReport)(result)) {
            if ((0, spectr_1.isValid)(result)) {
                console.log("  ✓ Valid");
            }
            else {
                console.error("  ✗ Invalid");
                // Print all issues
                const issues = (0, spectr_1.formatAllIssues)(result);
                for (const issue of issues) {
                    console.error(`    ${issue}`);
                }
            }
        }
    }
    // Get all errors across all results
    const allErrors = (0, spectr_1.getAllErrors)(validationOutput);
    if (allErrors.length > 0) {
        console.error("\n=== All Errors ===");
        for (const error of allErrors) {
            console.error((0, spectr_1.formatIssue)(error));
        }
    }
    // Get failed results only
    const failedResults = (0, spectr_1.getFailedResults)(validationOutput);
    console.log(`\nFailed items: ${failedResults.map((r) => r.name).join(", ")}`);
}
/**
 * Example: Create GitHub Actions annotations from validation output
 */
function createGitHubAnnotations(validationOutput) {
    const annotations = [];
    for (const result of validationOutput) {
        if (!(0, spectr_1.hasReport)(result))
            continue;
        for (const issue of result.report.issues) {
            const level = issue.level.toLowerCase();
            const file = issue.path;
            const line = issue.line || 1;
            const message = issue.message;
            // GitHub Actions annotation format
            // ::error file={name},line={line}::{message}
            annotations.push(`::${level} file=${file},line=${line}::${message}`);
        }
    }
    return annotations;
}
/**
 * Example: Exit with appropriate code based on validation
 */
function exitWithValidationStatus(validationOutput) {
    const hasErrors = (0, spectr_1.hasAnyErrors)(validationOutput);
    const exitCode = hasErrors ? 1 : 0;
    process.exit(exitCode);
}
/**
 * Example: Filter results by type
 */
function getChangeResults(validationOutput) {
    return validationOutput.filter((r) => r.type === "change");
}
function getSpecResults(validationOutput) {
    return validationOutput.filter((r) => r.type === "spec");
}
/**
 * Example: Create summary for GitHub Actions
 */
function createSummary(validationOutput) {
    const totalItems = validationOutput.length;
    const validItems = validationOutput.filter(spectr_1.isValid).length;
    const totalErrors = (0, spectr_1.getTotalErrorCount)(validationOutput);
    const totalWarnings = (0, spectr_1.getTotalWarningCount)(validationOutput);
    const changes = getChangeResults(validationOutput);
    const specs = getSpecResults(validationOutput);
    let summary = "# Spectr Validation Results\n\n";
    summary += `- **Total Items**: ${totalItems}\n`;
    summary += `- **Valid**: ${validItems}\n`;
    summary += `- **Invalid**: ${totalItems - validItems}\n`;
    summary += `- **Errors**: ${totalErrors}\n`;
    summary += `- **Warnings**: ${totalWarnings}\n\n`;
    if (changes.length > 0) {
        summary += `## Changes (${changes.length})\n\n`;
        for (const change of changes) {
            const status = (0, spectr_1.isValid)(change) ? "✓" : "✗";
            summary += `- ${status} ${change.name}\n`;
        }
        summary += "\n";
    }
    if (specs.length > 0) {
        summary += `## Specs (${specs.length})\n\n`;
        for (const spec of specs) {
            const status = (0, spectr_1.isValid)(spec) ? "✓" : "✗";
            summary += `- ${status} ${spec.name}\n`;
        }
        summary += "\n";
    }
    return summary;
}
