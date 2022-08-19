module github.com/qf-tech/qftech-base

go 1.17

require (
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/pkg/errors v0.8.1
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
)

require (
	github.com/BurntSushi/toml v1.0.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

// replace github.com/qf-tech/qftech-base/pkg/wsclient v0.0.0-20211011135246-e08b4ffee600 => ../wsclient
