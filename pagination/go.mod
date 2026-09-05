module github.com/kalandramo/bald-crud/pagination

go 1.26.3


require (
	github.com/google/go-cmp v0.7.0
	github.com/kalandramo/bald/bconf v0.0.0-00010101000000-000000000000
	github.com/tx7do/go-utils v1.1.40
	github.com/tx7do/go-wind-plugins/encoding v0.0.1
	github.com/tx7do/go-wind-plugins/encoding/json v0.0.1
	go.einride.tech/aip v0.86.3
	google.golang.org/genproto/googleapis/api v0.0.0-20260810153831-ec0a7760b754
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/google/uuid v1.6.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260810153831-ec0a7760b754 // indirect
)

replace github.com/kalandramo/bald/bconf => ../../bald/bconf
