module github.com/gzjjyz/srvlib

go 1.20

require (
	github.com/995933447/confloader v0.0.0-20230314141707-e7b191386ae2
	github.com/995933447/gonetutil v0.0.0-20230517070832-763d0c3b1d7e
	github.com/995933447/log-go v0.0.0-20230420123341-5d684963433b
	github.com/995933447/redisgroup v0.0.0-20230510085956-718f047520a1
	github.com/elliotchance/testify-stats v1.0.3
	github.com/golang/protobuf v1.5.3
	github.com/golang/snappy v0.0.4
	github.com/gorilla/websocket v1.5.0
	github.com/gzjjyz/micro v0.0.5-0.20231016033527-908637979b70
	github.com/gzjjyz/trace v0.0.0-20230831064247-64e42e91ac5d
	github.com/huandu/go-clone v1.6.0
	github.com/huaweicloud/huaweicloud-sdk-go-obs v3.23.9+incompatible
	github.com/json-iterator/go v1.1.12
	github.com/nats-io/nats.go v1.30.2
	github.com/petermattis/goid v0.0.0-20230808133559-b036b712a89b
	github.com/pkg/sftp v1.13.5
	github.com/robfig/cron/v3 v3.0.1
	github.com/spf13/viper v1.17.0
	github.com/stretchr/testify v1.8.4
	golang.org/x/crypto v0.13.0
	golang.org/x/net v0.15.0
	google.golang.org/grpc v1.58.2
	gorm.io/driver/mysql v1.5.0
	gorm.io/gorm v1.25.1
)

require (
	github.com/995933447/simpletrace v0.0.0-20230217061256-c25a914bd376 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/sagikazarmark/locafero v0.3.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.10.0 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230913181813-007df8e322eb // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230920204549-e6e6cdab5c13 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

require (
	github.com/995933447/std-go v0.0.0-20220806175833-ab3496c0b696
	github.com/995933447/stringhelper-go v0.0.0-20221220072216-628db3bc29d8 // indirect
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gzjjyz/logger v1.0.1
	github.com/howeyc/fsnotify v0.9.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nats-io/nats-server/v2 v2.9.16 // indirect
	github.com/nats-io/nkeys v0.4.5 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/satori/go.uuid v1.2.0
	go.etcd.io/etcd/api/v3 v3.5.9 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.9 // indirect
	go.etcd.io/etcd/client/v3 v3.5.9 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto v0.0.0-20230913181813-007df8e322eb // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/coreos/bbolt v1.3.7 => go.etcd.io/bbolt v1.3.7
	github.com/derekparker/delve => github.com/go-delve/delve v1.20.1
	github.com/go-delve/delve => github.com/derekparker/delve v1.4.0
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
