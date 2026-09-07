module github.com/kalandramo/bald-crud/elasticsearch

go 1.27.1

replace github.com/kalandramo/bald-crud/pagination => ../pagination

require (
	github.com/elastic/elastic-transport-go/v8 v8.11.0
	github.com/elastic/go-elasticsearch/v9 v9.5.0
	github.com/kalandramo/bald/bconf v0.0.0-00010101000000-000000000000
	github.com/kalandramo/bald/berrors v0.0.0-00010101000000-000000000000
	github.com/kalandramo/bald/encoding/json v0.0.0-00010101000000-000000000000
	github.com/kalandramo/bald/log v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/kalandramo/bald/encoding v0.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.45.0 // indirect
	go.opentelemetry.io/otel/metric v1.45.0 // indirect
	go.opentelemetry.io/otel/sdk v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.45.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/kalandramo/bald-crud => ../

replace github.com/kalandramo/bald/bconf => ../../bald/bconf

replace github.com/kalandramo/bald/log => ../../bald/log

replace github.com/kalandramo/bald/berrors => ../../bald/berrors

replace github.com/kalandramo/bald/encoding => ../../bald/encoding

replace github.com/kalandramo/bald/encoding/json => ../../bald/encoding/json
