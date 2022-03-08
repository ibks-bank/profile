.PHONY: generate
generate:
	mkdir -p "pkg"
	protoc -I/usr/local/include -I. \
		-Ivendor.protogen \
		--grpc-gateway_out=logtostderr=true:./pkg \
		--swagger_out=allow_merge=true,merge_file_name=api:./swagger \
		--go_out=plugins=grpc:./pkg ./api/profile/profile.proto

	statik -src=./swagger -dest ./cmd -p swagger

.PHONY: migrate-up
migrate-up:
	(cd migrations; goose postgres "host=localhost port=5433 user=postgres password=postgres database=profile sslmode=disable timezone=UTC" up)

.PHONY: migrate-down
migrate-down:
	(cd migrations; goose postgres "host=localhost port=5433 user=postgres password=postgres database=profile sslmode=disable timezone=UTC" down)