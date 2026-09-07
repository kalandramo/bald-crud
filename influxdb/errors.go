package influxdb

import berrors "github.com/kalandramo/bald/berrors"

var (
	ErrInfluxDBClientNotInitialized = berrors.Internal("INFLUXDB_CLIENT_NOT_INITIALIZED")

	ErrInfluxDBConnectFailed = berrors.Internal("INFLUXDB_CONNECT_FAILED")

	ErrInfluxDBCreateDatabaseFailed = berrors.Internal("INFLUXDB_CREATE_DATABASE_FAILED")

	ErrInfluxDBQueryFailed = berrors.Internal("INFLUXDB_QUERY_FAILED")

	ErrClientNotConnected = berrors.Internal("INFLUXDB_CLIENT_NOT_CONNECTED")

	ErrInvalidPoint = berrors.Internal("INFLUXDB_INVALID_POINT")

	ErrNoPointsToInsert = berrors.Internal("INFLUXDB_NO_POINTS_TO_INSERT")

	ErrEmptyData = berrors.Internal("INFLUXDB_EMPTY_DATA")

	ErrBatchInsertFailed = berrors.Internal("INFLUXDB_BATCH_INSERT_FAILED")

	ErrInsertFailed = berrors.Internal("INFLUXDB_INSERT_FAILED")
)
