package doris

import berrors "github.com/kalandramo/bald/berrors"

// Sentinel errors for the doris client. These are defined for API parity with
// the other go-crud client modules; the CRUD methods currently still return
// the underlying sqlx/database errors (or fmt.Errorf) unchanged. Callers may
// match on these sentinels if/when the methods are migrated to return them.
var (
	ErrClientNotInitialized = berrors.Internal("CLIENT_NOT_INITIALIZED")
	ErrQueryFailed          = berrors.Internal("QUERY_FAILED")
	ErrInsertFailed         = berrors.Internal("INSERT_FAILED")
	ErrBatchInsertFailed    = berrors.Internal("BATCH_INSERT_FAILED")
)
