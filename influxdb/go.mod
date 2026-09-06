module github.com/kalandramo/bald-crud/influxdb

go 1.26.5

replace github.com/kalandramo/bald-crud/pagination => ../pagination

require (
	github.com/InfluxCommunity/influxdb3-go/v2 v2.17.0
	github.com/kalandramo/bald-crud/pagination v0.0.15
	github.com/kalandramo/bald-utils v0.0.0-00010101000000-000000000000
	github.com/kalandramo/bald/bconf v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
	github.com/tx7do/go-wind v0.0.2
	github.com/tx7do/go-wind-plugins/encoding v0.0.1
	github.com/tx7do/go-wind-plugins/encoding/json v0.0.1
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/apache/arrow-go/v18 v18.7.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/frankban/quicktest v1.14.0 // indirect
	github.com/goccy/go-json v0.10.6 // indirect
	github.com/google/flatbuffers v25.12.19+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/influxdata/line-protocol/v2 v2.2.1 // indirect
	github.com/klauspost/compress v1.19.2 // indirect
	github.com/klauspost/cpuid/v2 v2.4.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.28 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/zeebo/xxh3 v1.1.0 // indirect
	go.einride.tech/aip v0.86.3 // indirect
	golang.org/x/exp v0.0.0-20260813180055-c1d0aacb2297 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260810153831-ec0a7760b754 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260810153831-ec0a7760b754 // indirect
	google.golang.org/grpc v1.83.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/kalandramo/bald-crud => ../

replace github.com/kalandramo/bald/bconf => ../../bald/bconf

replace github.com/kalandramo/bald-utils => ../../bald-utils
