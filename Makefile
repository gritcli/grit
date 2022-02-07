TEST_CONFIG_DIR := ./internal/testdata/etc
GO_EMBEDDED_FILES += $(shell PATH="$(PATH)" git-find cli/internal/commands -name 'help.txt')

-include .makefiles/Makefile
-include .makefiles/pkg/protobuf/v2/Makefile
-include .makefiles/pkg/go/v1/Makefile

run: $(GO_DEBUG_DIR)/grit artifacts/grit
	$< --socket artifacts/grit/daemon.sock $(RUN_ARGS)

serve: $(GO_DEBUG_DIR)/gritd artifacts/grit
	GRIT_CONFIG_DIR="internal/testdata/etc" $< $(RUN_ARGS)

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

artifacts/grit:
	mkdir -p $@
