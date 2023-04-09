GO_EMBEDDED_FILES += $(shell PATH="$(PATH)" git-find cli/internal/commands -name 'help.txt')
GO_EMBEDDED_FILES += $(shell PATH="$(PATH)" git-find cli/internal/commands/setupshell -name 'install.*')
GO_FERRITE_BINARY = gritd

-include .makefiles/Makefile
-include .makefiles/pkg/protobuf/v2/Makefile
-include .makefiles/pkg/go/v1/Makefile
-include .makefiles/pkg/go/v1/with-ferrite.mk

run: $(GO_DEBUG_DIR)/grit artifacts/grit
	$< --socket artifacts/grit/daemon.sock $(args)

serve: $(GO_DEBUG_DIR)/gritd artifacts/grit
	GRIT_CONFIG_DIR="testdata/etc" $< $(RUN_ARGS)

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

artifacts/grit:
	mkdir -p $@
