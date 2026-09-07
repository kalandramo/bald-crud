package opensearch

import (
	"fmt"

	berrors "github.com/kalandramo/bald/berrors"
)

var (
	// ErrRequestFailed is returned when a request to Elasticsearch fails.
	ErrRequestFailed = berrors.Internal("REQUEST_FAILED")

	// ErrIndexNotFound is returned when the specified index does not exist.
	ErrIndexNotFound = berrors.Internal("INDEX_NOT_FOUND")

	// ErrIndexAlreadyExists is returned when trying to create an index that already exists.
	ErrIndexAlreadyExists = berrors.Internal("INDEX_ALREADY_EXISTS")

	ErrCreateIndex = berrors.Internal("CREATE_INDEX_FAILED")

	ErrDeleteIndex = berrors.Internal("DELETE_INDEX_FAILED")

	// ErrDocumentNotFound is returned when a document is not found in the index.
	ErrDocumentNotFound = berrors.Internal("DOCUMENT_NOT_FOUND")

	// ErrUnmarshalResponse is returned when the response from Elasticsearch cannot be unmarshalled.
	ErrUnmarshalResponse = berrors.Internal("UNMARSHAL_RESPONSE_FAILED")

	ErrInsertDocument = berrors.Internal("INSERT_DOCUMENT_FAILED")

	ErrBatchInsertDocument = berrors.Internal("BATCH_INSERT_DOCUMENT_FAILED")

	ErrGetDocument = berrors.Internal("GET_DOCUMENT_FAILED")

	ErrSearchDocument = berrors.Internal("SEARCH_DOCUMENT_FAILED")

	ErrCreateISMPolicy = berrors.Internal("CREATE_ISM_POLICY_FAILED")

	ErrDocumentConflict = berrors.Internal("DOCUMENT_CONFLICT")

	ErrDeleteDocument = berrors.Internal("DELETE_DOCUMENT_FAILED")

	ErrUpdateDocument = berrors.Internal("UPDATE_DOCUMENT_FAILED")

	ErrDeleteISMPolicy = berrors.Internal("DELETE_ISM_POLICY_FAILED")

	ErrInvalidRequest = berrors.BadRequest("INVALID_REQUEST")

	ErrInvalidFilter = berrors.BadRequest("INVALID_FILTER")

	ErrCreateTemplate = berrors.Internal("CREATE_TEMPLATE_FAILED")

	ErrDeleteTemplate = berrors.Internal("DELETE_TEMPLATE_FAILED")
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
