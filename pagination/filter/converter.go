package filter

import storev1 "github.com/kalandramo/bald/bconf/gen/go/bald/store/v1"

var (
	queryStringConverter  = NewQueryStringConverter()
	filterStringConverter = NewFilterStringConverter()
)

type filterRequester interface {
	GetFilterExpr() *storev1.FilterExpr
	GetQuery() string
	GetFilter() string
}

// convertFilterRequest converts a filterRequester to a FilterExpr.
func convertFilterRequest(req filterRequester) (*storev1.FilterExpr, error) {
	if req == nil {
		return nil, nil
	}

	if req.GetFilterExpr() != nil {
		return req.GetFilterExpr(), nil
	}

	if q := req.GetQuery(); q != "" {
		return queryStringConverter.Convert(q)
	}

	if f := req.GetFilter(); f != "" {
		return filterStringConverter.Convert(f)
	}

	return nil, nil
}

// ConvertFilterByPagingRequest converts a PagingRequest to a FilterExpr.
func ConvertFilterByPagingRequest(req *storev1.PagingRequest) (*storev1.FilterExpr, error) {
	return convertFilterRequest(req)
}

// ConvertFilterByPaginationRequest converts a PaginationRequest to a FilterExpr.
func ConvertFilterByPaginationRequest(req *storev1.PaginationRequest) (*storev1.FilterExpr, error) {
	return convertFilterRequest(req)
}
