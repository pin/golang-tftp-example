# CAVEAT: path to workspace containing spaces does not work, see: http://savannah.gnu.org/bugs/?712

WORKSPACE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
export GOPATH := $(WORKSPACE_DIR)
all:
	@go get ./...

clean:
	@rm -rf $(WORKSPACE_DIR)bin $(WORKSPACE_DIR)pkg
