/**
 * Example usage of spectr validation types
 * This file demonstrates how to use the types in the action code
 */

import type { BulkResult, ValidationIssue, ValidationOutput } from "./spectr";
import {
  allValid,
  formatAllIssues,
  formatIssue,
  getAllErrors,
  getFailedResults,
  getTotalErrorCount,
  getTotalWarningCount,
  hasAnyErrors,
  hasError,
  hasReport,
  isValid,
} from "./spectr";

/**
 * Example: Parse and process spectr validation output
 */
function processValidationOutput(jsonOutput: string): void {
  // Parse JSON output from spectr command
  const validationOutput: ValidationOutput = JSON.parse(jsonOutput);

  // Check if all items are valid
  if (allValid(validationOutput)) {
    console.log("✓ All items passed validation");
    return;
  }

  // Check for any errors
  if (hasAnyErrors(validationOutput)) {
    console.error("✗ Validation failed with errors");
  }

  // Print summary
  const totalErrors = getTotalErrorCount(validationOutput);
  const totalWarnings = getTotalWarningCount(validationOutput);
  console.log(`Errors: ${totalErrors}, Warnings: ${totalWarnings}`);

  // Process each result
  for (const result of validationOutput) {
    console.log(`\nChecking: ${result.name} (${result.type})`);

    // Handle validation errors
    if (hasError(result)) {
      console.error(`  Error: ${result.error}`);
      continue;
    }

    // Handle validation report
    if (hasReport(result)) {
      if (isValid(result)) {
        console.log("  ✓ Valid");
      } else {
        console.error("  ✗ Invalid");
        // Print all issues
        const issues = formatAllIssues(result);
        for (const issue of issues) {
          console.error(`    ${issue}`);
        }
      }
    }
  }

  // Get all errors across all results
  const allErrors: ValidationIssue[] = getAllErrors(validationOutput);
  if (allErrors.length > 0) {
    console.error("\n=== All Errors ===");
    for (const error of allErrors) {
      console.error(formatIssue(error));
    }
  }

  // Get failed results only
  const failedResults = getFailedResults(validationOutput);
  console.log(`\nFailed items: ${failedResults.map((r) => r.name).join(", ")}`);
}

/**
 * Example: Create GitHub Actions annotations from validation output
 */
function createGitHubAnnotations(validationOutput: ValidationOutput): string[] {
  const annotations: string[] = [];

  for (const result of validationOutput) {
    if (!hasReport(result)) continue;

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
function exitWithValidationStatus(validationOutput: ValidationOutput): never {
  const hasErrors = hasAnyErrors(validationOutput);
  const exitCode = hasErrors ? 1 : 0;
  process.exit(exitCode);
}

/**
 * Example: Filter results by type
 */
function getChangeResults(validationOutput: ValidationOutput): BulkResult[] {
  return validationOutput.filter((r) => r.type === "change");
}

function getSpecResults(validationOutput: ValidationOutput): BulkResult[] {
  return validationOutput.filter((r) => r.type === "spec");
}

/**
 * Example: Create summary for GitHub Actions
 */
function createSummary(validationOutput: ValidationOutput): string {
  const totalItems = validationOutput.length;
  const validItems = validationOutput.filter(isValid).length;
  const totalErrors = getTotalErrorCount(validationOutput);
  const totalWarnings = getTotalWarningCount(validationOutput);

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
      const status = isValid(change) ? "✓" : "✗";
      summary += `- ${status} ${change.name}\n`;
    }
    summary += "\n";
  }

  if (specs.length > 0) {
    summary += `## Specs (${specs.length})\n\n`;
    for (const spec of specs) {
      const status = isValid(spec) ? "✓" : "✗";
      summary += `- ${status} ${spec.name}\n`;
    }
    summary += "\n";
  }

  return summary;
}

// Export examples for documentation
export {
  processValidationOutput,
  createGitHubAnnotations,
  exitWithValidationStatus,
  getChangeResults,
  getSpecResults,
  createSummary,
};
