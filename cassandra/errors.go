package cassandra

import berrors "github.com/kalandramo/bald/berrors"

var (
	// ErrInvalidRequest is returned when an input guard rejects a call before
	// it reaches the cluster (e.g. empty statement).
	ErrInvalidRequest = berrors.BadRequest("INVALID_REQUEST")

	// ErrExecQuery is returned when a query or batch execution fails.
	ErrExecQuery = berrors.Internal("EXEC_QUERY_FAILED")

	// ErrSessionClosed is returned when the underlying session has been closed.
	ErrSessionClosed = berrors.Internal("SESSION_CLOSED")
)
