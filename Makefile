.PHONY: install-tools
install-tools:
	@./scripts/install-tools.sh

.PHONY: run-catalogs-service
run-catalogs-service:
	@./scripts/run.sh catalogs

.PHONY: build
build:
	@./scripts/build.sh  pkg
	@./scripts/build.sh  catalogs

.PHONY: update-dependencies
update-dependencies:
	@./scripts/update-dependencies.sh  pkg
	@./scripts/update-dependencies.sh  catalogs

.PHONY: install-dependencies
install-dependencies:
	@./scripts/install-dependencies.sh  pkg
	@./scripts/install-dependencies.sh  catalogs

.PHONY: docker-compose-infra-up
docker-compose-infra-up:
	@docker-compose -f deployments/docker-compose/docker-compose.infrastructure.yaml up --build -d

docker-compose-infra-down:
	@docker-compose -f deployments/docker-compose/docker-compose.infrastructure.yaml down

.PHONY: openapi
openapi:
	@./scripts/openapi.sh catalogs

# https://stackoverflow.com/questions/13616033/install-protocol-buffers-on-windows
.PHONY: proto
proto:
	@./scripts/proto.sh catalogs

.PHONY: unit-test
unit-test:
	@./scripts/test.sh catalogs unit

.PHONY: integration-test
integration-test:
	@./scripts/test.sh catalogs integration

.PHONY: e2e-test
e2e-test:
	@./scripts/test.sh catalogs e2e

#.PHONY: load-test
#load-test:
#	@./scripts/test.sh catalogs_write load-test
#	@./scripts/test.sh catalogs_read load-test
#	@./scripts/test.sh  orders load-test

.PHONY: format
format:
	@./scripts/format.sh catalogs
	@./scripts/format.sh pkg

.PHONY: lint
lint:
	@./scripts/lint.sh catalogs
	@./scripts/lint.sh pkg

# https://github.com/golang-migrate/migrate/blob/856ea12df9d230b0145e23d951b7dbd6b86621cb/database/postgres/TUTORIAL.md
# https://github.com/golang-migrate/migrate/blob/856ea12df9d230b0145e23d951b7dbd6b86621cb/GETTING_STARTED.md
# https://github.com/golang-migrate/migrate/blob/856ea12df9d230b0145e23d951b7dbd6b86621cb/MIGRATIONS.md
# https://github.com/golang-migrate/migrate/tree/856ea12df9d230b0145e23d951b7dbd6b86621cb/cmd/migrate#usage
.PHONY: go-migrate
go-migrate:
	@./scripts/go-migrate.sh -p ./internal/services/catalogs/db/migrations/go-migrate -c create -n create_product_table
	@./scripts/go-migrate.sh -p ./internal/services/catalogs/db/migrations/go-migrate -c up -o postgres://postgres:123456@localhost:5432/catalogs?sslmode=disable
	@./scripts/go-migrate.sh -p ./internal/services/catalogs/db/migrations/go-migrate -c down -o postgres://postgres:123456@localhost:5432/catalogs?sslmode=disable

# https://github.com/pressly/goose#usage
.PHONY: goose-migrate
goose-migrate:
	@./scripts/goose-migrate.sh -p ./internal/services/catalogs/db/migrations/goose-migrate -c create -n create_product_table
	@./scripts/goose-migrate.sh -p ./internal/services/catalogs/db/migrations/goose-migrate -c up -o "user=postgres password=123456 dbname=catalogs sslmode=disable"
	@./scripts/goose-migrate.sh -p ./internal/services/catalogs/db/migrations/goose-migrate -c down -o "user=postgres password=123456 dbname=catalogs sslmode=disable"

# https://atlasgo.io/guides/orms/gorm
.PHONY: atlas
atlas:
	@./scripts/atlas-migrate.sh -c gorm-sync -p "./internal/services/catalogs"
	@./scripts/atlas-migrate.sh -c apply -p "./internal/services/catalogs" -o "postgres://postgres:123456@localhost:5432/catalogs?sslmode=disable"

.PHONY: cycle-check
cycle-check:
	cd internal/pkg && goimportcycle -dot imports.dot dot -Tpng -o cycle/pkg.png imports.dot
	cd internal/services/catalogs && goimportcycle -dot imports.dot dot -Tpng -o cycle/catalogs.png imports.dot

#.PHONY: c4
#c4:
#	cd tools/c4 && go mod tidy && sh generate.sh

# https://medium.com/yemeksepeti-teknoloji/mocking-an-interface-using-mockery-in-go-afbcb83cc773
# https://vektra.github.io/mockery/latest/running/
# https://amitshekhar.me/blog/test-with-testify-and-mockery-in-go
.PHONY: pkg-mocks
pkg-mocks:
	cd internal/pkg/es && mockery --output mocks --all
	cd internal/pkg/core/serializer && mockery --output mocks --all
	cd internal/pkg/core/messaging && mockery --output mocks --all

.PHONY: services-mocks
services-mocks:
	cd internal/services/catalogs && mockery --output mocks --all
