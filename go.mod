module github.com/ilia-tsyplenkov/klines-stat

go 1.23

toolchain go1.23.6

require (
	github.com/jackc/pgx/v5 v5.7.2
	github.com/sirupsen/logrus v1.9.3
	github.com/wuhewuhe/bybit.go.api v1.0.18
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/bitly/go-simplejson v0.5.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)

replace github.com/wuhewuhe/bybit.go.api v1.0.18 => github.com/ilia-tsyplenkov/bybit.go.api v0.0.2
