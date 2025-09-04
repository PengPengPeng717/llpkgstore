package actions

import (
	"fmt"
	"time"
)

// ErrorCode represents different types of errors
type ErrorCode string

const (
	// Configuration errors
	ErrorCodeConfigNotFound    ErrorCode = "CONFIG_NOT_FOUND"
	ErrorCodeConfigInvalid     ErrorCode = "CONFIG_INVALID"
	ErrorCodeConfigParseFailed ErrorCode = "CONFIG_PARSE_FAILED"

	// Version errors
	ErrorCodeVersionNotFound ErrorCode = "VERSION_NOT_FOUND"
	ErrorCodeVersionInvalid  ErrorCode = "VERSION_INVALID"
	ErrorCodeVersionExists   ErrorCode = "VERSION_EXISTS"

	// GitHub API errors
	ErrorCodeGitHubAuthFailed ErrorCode = "GITHUB_AUTH_FAILED"
	ErrorCodeGitHubAPIError   ErrorCode = "GITHUB_API_ERROR"
	ErrorCodeReleaseExists    ErrorCode = "RELEASE_EXISTS"
	ErrorCodeTagExists        ErrorCode = "TAG_EXISTS"

	// Package processing errors
	ErrorCodePackageNotFound    ErrorCode = "PACKAGE_NOT_FOUND"
	ErrorCodePackageInvalid     ErrorCode = "PACKAGE_INVALID"
	ErrorCodePackageTypeUnknown ErrorCode = "PACKAGE_TYPE_UNKNOWN"

	// System errors
	ErrorCodeSystemError     ErrorCode = "SYSTEM_ERROR"
	ErrorCodeNetworkError    ErrorCode = "NETWORK_ERROR"
	ErrorCodeFileSystemError ErrorCode = "FILESYSTEM_ERROR"
)

// ProcessingError represents a structured error with context
type ProcessingError struct {
	Code        ErrorCode              `json:"code"`
	Message     string                 `json:"message"`
	PackageType string                 `json:"package_type,omitempty"`
	PackageName string                 `json:"package_name,omitempty"`
	Version     string                 `json:"version,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Cause       error                  `json:"cause,omitempty"`
}

// Error implements the error interface
func (e *ProcessingError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *ProcessingError) Unwrap() error {
	return e.Cause
}

// NewProcessingError creates a new processing error
func NewProcessingError(code ErrorCode, message string, cause error) *ProcessingError {
	return &ProcessingError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Cause:     cause,
	}
}

// WithPackageInfo adds package information to the error
func (e *ProcessingError) WithPackageInfo(packageType, packageName, version string) *ProcessingError {
	e.PackageType = packageType
	e.PackageName = packageName
	e.Version = version
	return e
}

// WithDetails adds additional details to the error
func (e *ProcessingError) WithDetails(details map[string]interface{}) *ProcessingError {
	e.Details = details
	return e
}

// Error handling helper functions

// HandleConfigError handles configuration-related errors
func HandleConfigError(err error, configPath string) *ProcessingError {
	return NewProcessingError(ErrorCodeConfigParseFailed,
		fmt.Sprintf("failed to parse configuration file: %s", configPath), err)
}

// HandleVersionError handles version-related errors
func HandleVersionError(err error, version string) *ProcessingError {
	return NewProcessingError(ErrorCodeVersionInvalid,
		fmt.Sprintf("invalid version: %s", version), err)
}

// HandleGitHubError handles GitHub API errors
func HandleGitHubError(err error, operation string) *ProcessingError {
	return NewProcessingError(ErrorCodeGitHubAPIError,
		fmt.Sprintf("GitHub API error during %s", operation), err)
}

// HandlePackageError handles package processing errors
func HandlePackageError(err error, packageType, packageName string) *ProcessingError {
	return NewProcessingError(ErrorCodePackageInvalid,
		fmt.Sprintf("failed to process %s package: %s", packageType, packageName), err).
		WithPackageInfo(packageType, packageName, "")
}

// HandleSystemError handles system-level errors
func HandleSystemError(err error, operation string) *ProcessingError {
	return NewProcessingError(ErrorCodeSystemError,
		fmt.Sprintf("system error during %s", operation), err)
}

// Error recovery and retry mechanisms

// RetryableError represents an error that can be retried
type RetryableError struct {
	*ProcessingError
	MaxRetries int
	RetryDelay time.Duration
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	if retryableErr, ok := err.(*RetryableError); ok {
		return retryableErr.MaxRetries > 0
	}

	// Check for specific retryable error codes
	if processingErr, ok := err.(*ProcessingError); ok {
		switch processingErr.Code {
		case ErrorCodeNetworkError, ErrorCodeGitHubAPIError:
			return true
		}
	}

	return false
}

// NewRetryableError creates a new retryable error
func NewRetryableError(code ErrorCode, message string, cause error, maxRetries int, retryDelay time.Duration) *RetryableError {
	return &RetryableError{
		ProcessingError: NewProcessingError(code, message, cause),
		MaxRetries:      maxRetries,
		RetryDelay:      retryDelay,
	}
}

// Error context and logging

// ErrorContext provides context for error logging
type ErrorContext struct {
	Operation   string                 `json:"operation"`
	PackageType string                 `json:"package_type,omitempty"`
	PackageName string                 `json:"package_name,omitempty"`
	Version     string                 `json:"version,omitempty"`
	Repository  string                 `json:"repository,omitempty"`
	Branch      string                 `json:"branch,omitempty"`
	CommitSHA   string                 `json:"commit_sha,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// NewErrorContext creates a new error context
func NewErrorContext(operation string) *ErrorContext {
	return &ErrorContext{
		Operation: operation,
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
	}
}

// WithPackageInfo adds package information to the context
func (ec *ErrorContext) WithPackageInfo(packageType, packageName, version string) *ErrorContext {
	ec.PackageType = packageType
	ec.PackageName = packageName
	ec.Version = version
	return ec
}

// WithRepositoryInfo adds repository information to the context
func (ec *ErrorContext) WithRepositoryInfo(repository, branch, commitSHA string) *ErrorContext {
	ec.Repository = repository
	ec.Branch = branch
	ec.CommitSHA = commitSHA
	return ec
}

// WithDetails adds additional details to the context
func (ec *ErrorContext) WithDetails(details map[string]interface{}) *ErrorContext {
	for k, v := range details {
		ec.Details[k] = v
	}
	return ec
}

// LogError logs an error with context
func LogError(err error, context *ErrorContext) {
	// This would integrate with your logging system
	// For now, we'll just print to stdout
	fmt.Printf("ERROR [%s] %s: %v\n", context.Operation, context.Timestamp.Format(time.RFC3339), err)

	if context.PackageName != "" {
		fmt.Printf("  Package: %s (%s) %s\n", context.PackageName, context.PackageType, context.Version)
	}

	if context.Repository != "" {
		fmt.Printf("  Repository: %s (branch: %s, commit: %s)\n",
			context.Repository, context.Branch, context.CommitSHA)
	}

	if len(context.Details) > 0 {
		fmt.Printf("  Details: %+v\n", context.Details)
	}
}

