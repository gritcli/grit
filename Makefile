-include .makefiles/Makefile
-include .makefiles/pkg/protobuf/v2/Makefile
-include .makefiles/pkg/go/v1/Makefile

run: $(GO_DEBUG_DIR)/grit2 artifacts/grit
	$< --config etc/ $(RUN_ARGS)

serve: $(GO_DEBUG_DIR)/gritd artifacts/grit
	GRIT_CONFIG_DIR="etc/" $< $(RUN_ARGS)

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

artifacts/grit:
	mkdir -p $@
