package elasticsearch

import (
	"fmt"

	berrors "github.com/kalandramo/bald/berrors"
)

var (
	// ErrRequestFailed is returned when a request to Elasticsearch fails.
	ErrRequestFailed = berrors.Internal("REQUEST_FAILED")

	ErrInvalidRequest = berrors.BadRequest("INVALID_REQUEST")

	// ErrIndexNotFound is returned when the specified index does not exist.
	ErrIndexNotFound = berrors.Internal("INDEX_NOT_FOUND")

	// ErrIndexAlreadyExists is returned when trying to create an index that already exists.
	ErrIndexAlreadyExists = berrors.Internal("INDEX_ALREADY_EXISTS")

	ErrCreateIndex = berrors.Internal("CREATE_INDEX_FAILED")

	ErrDeleteIndex = berrors.Internal("DELETE_INDEX_FAILED")

	// ErrDocumentNotFound is returned when a document is not found in the index.
	ErrDocumentNotFound = berrors.Internal("DOCUMENT_NOT_FOUND")

	// ErrDocumentAlreadyExists is returned when trying to create a document that already exists.
	ErrDocumentAlreadyExists = berrors.Internal("DOCUMENT_ALREADY_EXISTS")

	// ErrInvalidQuery is returned when the query provided to Elasticsearch is invalid.
	ErrInvalidQuery = berrors.Internal("INVALID_QUERY")

	// ErrUnmarshalResponse is returned when the response from Elasticsearch cannot be unmarshalled.
	ErrUnmarshalResponse = berrors.Internal("UNMARSHAL_RESPONSE_FAILED")

	ErrInsertDocument = berrors.Internal("INSERT_DOCUMENT_FAILED")

	ErrBatchInsertDocument = berrors.Internal("BATCH_INSERT_DOCUMENT_FAILED")

	ErrGetDocument = berrors.Internal("GET_DOCUMENT_FAILED")

	ErrSearchDocument = berrors.Internal("SEARCH_DOCUMENT_FAILED")

	ErrUpdateDocument = berrors.Internal("UPDATE_DOCUMENT_FAILED")

	ErrDeleteDocument = berrors.Internal("DELETE_DOCUMENT_FAILED")

	ErrCreateILMPolicy = berrors.Internal("CREATE_ILM_POLICY_FAILED")

	ErrDocumentConflict = berrors.Internal("DOCUMENT_CONFLICT")
)

// PartialFailureError 表示批量操作部分失败
type PartialFailureError struct {
	Total     int
	Failed    int
	FailedIDs []string
}

func (e *PartialFailureError) Error() string {
	return fmt.Sprintf("bulk insert: %d/%d failed, failed IDs: %v", e.Failed, e.Total, e.FailedIDs)
}
