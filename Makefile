TEST_CONFIG_DIR := ./internal/testdata/etc

-include .makefiles/Makefile
-include .makefiles/pkg/protobuf/v2/Makefile
-include .makefiles/pkg/go/v1/Makefile

run: $(GO_DEBUG_DIR)/grit artifacts/grit
	$< --config $(TEST_CONFIG_DIR) $(RUN_ARGS)

serve: $(GO_DEBUG_DIR)/gritd artifacts/grit
	GRIT_CONFIG_DIR="$(TEST_CONFIG_DIR)" $< $(RUN_ARGS)

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

artifacts/grit:
	mkdir -p $@
