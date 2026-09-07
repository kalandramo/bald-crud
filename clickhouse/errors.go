package clickhouse

import berrors "github.com/kalandramo/bald/berrors"

var (
	// ErrInvalidColumnName is returned when an invalid column name is used.
	ErrInvalidColumnName = berrors.Internal("INVALID_COLUMN_NAME")

	// ErrInvalidTableName is returned when an invalid table name is used.
	ErrInvalidTableName = berrors.Internal("INVALID_TABLE_NAME")

	// ErrInvalidCondition is returned when an invalid condition is used in a query.
	ErrInvalidCondition = berrors.Internal("INVALID_CONDITION")

	// ErrQueryExecutionFailed is returned when a query execution fails.
	ErrQueryExecutionFailed = berrors.Internal("QUERY_EXECUTION_FAILED")

	// ErrExecutionFailed is returned when a general execution fails.
	ErrExecutionFailed = berrors.Internal("EXECUTION_FAILED")

	// ErrAsyncInsertFailed is returned when an asynchronous insert operation fails.
	ErrAsyncInsertFailed = berrors.Internal("ASYNC_INSERT_FAILED")

	// ErrRowScanFailed is returned when scanning rows from a query result fails.
	ErrRowScanFailed = berrors.Internal("ROW_SCAN_FAILED")

	// ErrRowsIterationError is returned when there is an error iterating over rows.
	ErrRowsIterationError = berrors.Internal("ROWS_ITERATION_ERROR")

	// ErrRowNotFound is returned when a specific row is not found in the result set.
	ErrRowNotFound = berrors.Internal("ROW_NOT_FOUND")

	// ErrConnectionFailed is returned when the connection to ClickHouse fails.
	ErrConnectionFailed = berrors.Internal("CONNECTION_FAILED")

	// ErrDatabaseNotFound is returned when the specified database is not found.
	ErrDatabaseNotFound = berrors.Internal("DATABASE_NOT_FOUND")

	// ErrTableNotFound is returned when the specified table is not found.
	ErrTableNotFound = berrors.Internal("TABLE_NOT_FOUND")

	// ErrInsertFailed is returned when an insert operation fails.
	ErrInsertFailed = berrors.Internal("INSERT_FAILED")

	// ErrUpdateFailed is returned when an update operation fails.
	ErrUpdateFailed = berrors.Internal("UPDATE_FAILED")

	// ErrDeleteFailed is returned when a delete operation fails.
	ErrDeleteFailed = berrors.Internal("DELETE_FAILED")

	// ErrTransactionFailed is returned when a transaction fails.
	ErrTransactionFailed = berrors.Internal("TRANSACTION_FAILED")

	// ErrClientNotInitialized is returned when the ClickHouse client is not initialized.
	ErrClientNotInitialized = berrors.Internal("CLIENT_NOT_INITIALIZED")

	// ErrGetServerVersionFailed is returned when getting the server version fails.
	ErrGetServerVersionFailed = berrors.Internal("GET_SERVER_VERSION_FAILED")

	// ErrPingFailed is returned when a ping to the ClickHouse server fails.
	ErrPingFailed = berrors.Internal("PING_FAILED")

	// ErrCreatorFunctionNil is returned when the creator function is nil.
	ErrCreatorFunctionNil = berrors.Internal("CREATOR_FUNCTION_NIL")

	// ErrBatchPrepareFailed is returned when a batch prepare operation fails.
	ErrBatchPrepareFailed = berrors.Internal("BATCH_PREPARE_FAILED")

	// ErrBatchSendFailed is returned when a batch send operation fails.
	ErrBatchSendFailed = berrors.Internal("BATCH_SEND_FAILED")

	// ErrBatchAppendFailed is returned when appending to a batch fails.
	ErrBatchAppendFailed = berrors.Internal("BATCH_APPEND_FAILED")

	// ErrBatchInsertFailed is returned when a batch insert operation fails.
	ErrBatchInsertFailed = berrors.Internal("BATCH_INSERT_FAILED")

	// ErrInvalidDSN is returned when the data source name (DSN) is invalid.
	ErrInvalidDSN = berrors.Internal("INVALID_DSN")

	// ErrInvalidProxyURL is returned when the proxy URL is invalid.
	ErrInvalidProxyURL = berrors.Internal("INVALID_PROXY_URL")

	// ErrPrepareInsertDataFailed is returned when preparing insert data fails.
	ErrPrepareInsertDataFailed = berrors.Internal("PREPARE_INSERT_DATA_FAILED")

	// ErrInvalidColumnData is returned when the column data type is invalid.
	ErrInvalidColumnData = berrors.Internal("INVALID_COLUMN_DATA")

	ErrInvalidArgument = berrors.BadRequest("INVALID_ARGUMENT")
)
