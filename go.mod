module github.com/hashicorp/boundary

go 1.15

replace github.com/hashicorp/boundary/api => ./api

require (
	github.com/armon/go-metrics v0.3.6
	github.com/bufbuild/buf v0.37.0
	github.com/dhui/dktest v0.3.4
	github.com/fatih/color v1.10.0
	github.com/favadi/protoc-go-inject-tag v1.1.0
	github.com/go-bindata/go-bindata/v3 v3.1.3
	github.com/go-swagger/go-swagger v0.26.1
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.4
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.2.0
	github.com/hashicorp/cap v0.0.0-20210223021404-4b74380f71bd
	github.com/hashicorp/boundary/api v0.0.8
	github.com/hashicorp/boundary/sdk v0.0.4
	github.com/hashicorp/dbassert v0.0.0-20200930125617-6218396928df
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-bexpr v0.1.7
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-hclog v0.15.0
	github.com/hashicorp/go-kms-wrapping v0.6.1
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/go-retryablehttp v0.6.8
	github.com/hashicorp/go-uuid v1.0.2
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/shared-secure-libs v0.0.4
	github.com/hashicorp/vault/sdk v0.1.14-0.20200916184745-5576096032f8
	github.com/iancoleman/strcase v0.1.3
	github.com/jefferai/keyring v1.1.7-0.20210105022822-8749b3d9ce79
	github.com/jinzhu/gorm v1.9.16
	github.com/jinzhu/now v1.1.1 // indirect
	github.com/kr/pretty v0.2.1
	github.com/kr/text v0.2.0
	github.com/lib/pq v1.9.0
	github.com/mattn/go-colorable v0.1.8
	github.com/mitchellh/cli v1.1.2
	github.com/mitchellh/go-wordwrap v1.0.1
	github.com/mitchellh/gox v1.0.1
	github.com/mitchellh/mapstructure v1.4.1
	github.com/mitchellh/pointerstructure v1.2.0
	github.com/mr-tron/base58 v1.2.0
	github.com/oligot/go-mod-upgrade v0.4.0
	github.com/ory/dockertest/v3 v3.6.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pires/go-proxyproto v0.5.0
	github.com/pkg/errors v0.9.1
	github.com/posener/complete v1.2.3
	github.com/stretchr/testify v1.7.0
	github.com/zalando/go-keyring v0.1.1
	go.uber.org/atomic v1.7.0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/sys v0.0.0-20210220050731-9a76102bfb43
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d
	golang.org/x/tools v0.1.0
	google.golang.org/genproto v0.0.0-20210222152913-aa3ee6e6a81c
	google.golang.org/grpc v1.35.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.25.1-0.20201208041424-160c7477e0e8
	mvdan.cc/gofumpt v0.1.0
	nhooyr.io/websocket v1.8.6
)
