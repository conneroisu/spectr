"use strict";
/**
 * TypeScript type definitions for spectr validation output
 *
 * These types match the JSON output structure from:
 * `spectr validate --all --strict --json`
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.formatAllIssues = exports.formatIssue = exports.getResultsWithErrors = exports.getFailedResults = exports.allValid = exports.hasAnyErrors = exports.getTotalInfoCount = exports.getTotalWarningCount = exports.getTotalErrorCount = exports.getAllInfo = exports.getAllWarnings = exports.getAllErrors = exports.getTotalIssueCount = exports.hasError = exports.hasReport = exports.isValid = void 0;
/**
 * Type guard to check if validation result is valid
 */
const isValid = (result) => result.valid;
exports.isValid = isValid;
/**
 * Type guard to check if result has a validation report
 */
const hasReport = (result) => result.report !== undefined;
exports.hasReport = hasReport;
/**
 * Type guard to check if result has an error
 */
const hasError = (result) => result.error !== undefined;
exports.hasError = hasError;
/**
 * Count total issues across all validation results
 */
const getTotalIssueCount = (output) => output.reduce((sum, r) => sum + (r.report?.issues.length || 0), 0);
exports.getTotalIssueCount = getTotalIssueCount;
/**
 * Get all errors from all validation results
 */
const getAllErrors = (output) => output
    .filter(exports.hasReport)
    .flatMap((r) => r.report.issues.filter((i) => i.level === "ERROR"));
exports.getAllErrors = getAllErrors;
/**
 * Get all warnings from all validation results
 */
const getAllWarnings = (output) => output
    .filter(exports.hasReport)
    .flatMap((r) => r.report.issues.filter((i) => i.level === "WARNING"));
exports.getAllWarnings = getAllWarnings;
/**
 * Get all info messages from all validation results
 */
const getAllInfo = (output) => output
    .filter(exports.hasReport)
    .flatMap((r) => r.report.issues.filter((i) => i.level === "INFO"));
exports.getAllInfo = getAllInfo;
/**
 * Count total errors across all validation results
 */
const getTotalErrorCount = (output) => output.reduce((sum, r) => sum + (r.report?.summary.errors || 0), 0);
exports.getTotalErrorCount = getTotalErrorCount;
/**
 * Count total warnings across all validation results
 */
const getTotalWarningCount = (output) => output.reduce((sum, r) => sum + (r.report?.summary.warnings || 0), 0);
exports.getTotalWarningCount = getTotalWarningCount;
/**
 * Count total info messages across all validation results
 */
const getTotalInfoCount = (output) => output.reduce((sum, r) => sum + (r.report?.summary.info || 0), 0);
exports.getTotalInfoCount = getTotalInfoCount;
/**
 * Check if any validation result has errors
 */
const hasAnyErrors = (output) => output.some((r) => r.report && r.report.summary.errors > 0);
exports.hasAnyErrors = hasAnyErrors;
/**
 * Check if all validation results are valid
 */
const allValid = (output) => output.every(exports.isValid);
exports.allValid = allValid;
/**
 * Get all failed validation results
 */
const getFailedResults = (output) => output.filter((r) => !r.valid);
exports.getFailedResults = getFailedResults;
/**
 * Get all validation results with errors
 */
const getResultsWithErrors = (output) => output.filter((r) => r.report && r.report.summary.errors > 0);
exports.getResultsWithErrors = getResultsWithErrors;
/**
 * Format validation issue for display
 */
const formatIssue = (issue) => {
    const location = issue.line ? `${issue.path}:${issue.line}` : issue.path;
    return `[${issue.level}] ${location}: ${issue.message}`;
};
exports.formatIssue = formatIssue;
/**
 * Format all issues from a validation result
 */
const formatAllIssues = (result) => {
    if (!(0, exports.hasReport)(result)) {
        return result.error ? [`Error: ${result.error}`] : [];
    }
    return result.report.issues.map(exports.formatIssue);
};
exports.formatAllIssues = formatAllIssues;
