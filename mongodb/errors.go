package mongodb

import berrors "github.com/kalandramo/bald/berrors"

// Sentinel errors for the mongodb client. These are defined for API parity
// with the other go-crud client modules; the CRUD methods currently still
// return the underlying mongo-driver errors unchanged. Callers may match on
// these sentinels if/when the methods are migrated to return them.
var (
	ErrClientNotInitialized = berrors.Internal("CLIENT_NOT_INITIALIZED")
	ErrQueryFailed          = berrors.Internal("QUERY_FAILED")
	ErrInsertFailed         = berrors.Internal("INSERT_FAILED")
	ErrUpdateFailed         = berrors.Internal("UPDATE_FAILED")
	ErrDeleteFailed         = berrors.Internal("DELETE_FAILED")
	ErrCountFailed          = berrors.Internal("COUNT_FAILED")
)
