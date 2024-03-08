// package shared contains the shared interface definition
package iExec

// CdExecutor represents the interface for executing commands in a plugin.
type CdExecutor interface {
	// CdExec is a method that executes a command and returns the result.
	CdExec(jsonInput string) (string, error)
}
