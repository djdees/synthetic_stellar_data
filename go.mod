module djdees/synthetic_steller_data

go 1.23

toolchain go1.24.10

require (
	github.com/apache/cassandra-gocql-driver/v2 v2.0.0
	github.com/google/uuid v1.2.0
	github.com/xitongsys/parquet-go v1.6.2
	github.com/xitongsys/parquet-go-source v0.0.0-20220315005136-aec0fe3e777c
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/apache/arrow/go/arrow v0.0.0-20200730104253-651201b0f516 // indirect
	github.com/apache/thrift v0.14.2 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/klauspost/compress v1.13.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.8 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

replace github.com/hailocab/go-hostpool => github.com/bitly/go-hostpool v0.1.0
