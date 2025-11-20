/**
 * TypeScript type definitions for spectr validation output
 *
 * These types match the JSON output structure from:
 * `spectr validate --all --strict --json`
 */

/**
 * Validation issue severity levels
 */
export type ValidationLevel = "ERROR" | "WARNING" | "INFO";

/**
 * A single validation issue with location and message
 */
export interface ValidationIssue {
  /** Severity level of the issue */
  level: ValidationLevel;
  /** File path where the issue was found */
  path: string;
  /** Optional line number where the issue occurs */
  line?: number;
  /** Human-readable description of the issue */
  message: string;
}

/**
 * Summary of validation issues by severity
 */
export interface ValidationSummary {
  /** Number of error-level issues */
  errors: number;
  /** Number of warning-level issues */
  warnings: number;
  /** Number of info-level issues */
  info: number;
}

/**
 * Validation report for a single change or spec
 */
export interface ValidationReport {
  /** Whether the item passed validation */
  valid: boolean;
  /** List of validation issues found */
  issues: ValidationIssue[];
  /** Summary counts by severity level */
  summary: ValidationSummary;
}

/**
 * Type of item being validated
 */
export type ValidationType = "change" | "spec";

/**
 * Single result from bulk validation
 */
export interface BulkResult {
  /** Item identifier (change ID or spec ID) */
  name: string;
  /** Type of item being validated */
  type: ValidationType;
  /** Overall validation status */
  valid: boolean;
  /** Validation report if validation ran successfully */
  report?: ValidationReport;
  /** Error message if validation failed to run */
  error?: string;
}

/**
 * Array of validation results returned by bulk validation
 * Output from: `spectr validate --all --strict --json`
 */
export type ValidationOutput = BulkResult[];

/**
 * Type guard to check if validation result is valid
 */
export const isValid = (result: BulkResult): boolean => result.valid;

/**
 * Type guard to check if result has a validation report
 */
export const hasReport = (
  result: BulkResult,
): result is BulkResult & { report: ValidationReport } =>
  result.report !== undefined;

/**
 * Type guard to check if result has an error
 */
export const hasError = (
  result: BulkResult,
): result is BulkResult & { error: string } => result.error !== undefined;

/**
 * Count total issues across all validation results
 */
export const getTotalIssueCount = (output: ValidationOutput): number =>
  output.reduce((sum, r) => sum + (r.report?.issues.length || 0), 0);

/**
 * Get all errors from all validation results
 */
export const getAllErrors = (output: ValidationOutput): ValidationIssue[] =>
  output
    .filter(hasReport)
    .flatMap((r) => r.report.issues.filter((i) => i.level === "ERROR"));

/**
 * Get all warnings from all validation results
 */
export const getAllWarnings = (output: ValidationOutput): ValidationIssue[] =>
  output
    .filter(hasReport)
    .flatMap((r) => r.report.issues.filter((i) => i.level === "WARNING"));

/**
 * Get all info messages from all validation results
 */
export const getAllInfo = (output: ValidationOutput): ValidationIssue[] =>
  output
    .filter(hasReport)
    .flatMap((r) => r.report.issues.filter((i) => i.level === "INFO"));

/**
 * Count total errors across all validation results
 */
export const getTotalErrorCount = (output: ValidationOutput): number =>
  output.reduce((sum, r) => sum + (r.report?.summary.errors || 0), 0);

/**
 * Count total warnings across all validation results
 */
export const getTotalWarningCount = (output: ValidationOutput): number =>
  output.reduce((sum, r) => sum + (r.report?.summary.warnings || 0), 0);

/**
 * Count total info messages across all validation results
 */
export const getTotalInfoCount = (output: ValidationOutput): number =>
  output.reduce((sum, r) => sum + (r.report?.summary.info || 0), 0);

/**
 * Check if any validation result has errors
 */
export const hasAnyErrors = (output: ValidationOutput): boolean =>
  output.some((r) => r.report && r.report.summary.errors > 0);

/**
 * Check if all validation results are valid
 */
export const allValid = (output: ValidationOutput): boolean =>
  output.every(isValid);

/**
 * Get all failed validation results
 */
export const getFailedResults = (output: ValidationOutput): BulkResult[] =>
  output.filter((r) => !r.valid);

/**
 * Get all validation results with errors
 */
export const getResultsWithErrors = (output: ValidationOutput): BulkResult[] =>
  output.filter((r) => r.report && r.report.summary.errors > 0);

/**
 * Format validation issue for display
 */
export const formatIssue = (issue: ValidationIssue): string => {
  const location = issue.line ? `${issue.path}:${issue.line}` : issue.path;
  return `[${issue.level}] ${location}: ${issue.message}`;
};

/**
 * Format all issues from a validation result
 */
export const formatAllIssues = (result: BulkResult): string[] => {
  if (!hasReport(result)) {
    return result.error ? [`Error: ${result.error}`] : [];
  }
  return result.report.issues.map(formatIssue);
};
