# collectjs Makefile

GOCMD = GO111MODULE=on go
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GOTEST = $(GOCMD) test
export GO111MODULE ?= on
CGO			?= 0



LOCAL_OS := $(shell uname)
ifeq ($(LOCAL_OS),Linux)
   TARGET_OS_LOCAL = linux
   GOLANGCI_LINT:=golangci-lint
else ifeq ($(LOCAL_OS),Darwin)
   TARGET_OS_LOCAL = darwin
   GOLANGCI_LINT:=golangci-lint
else
   TARGET_OS_LOCAL ?= windows
   GOLANGCI_LINT:=golangci-lint.exe
endif


export GOOS ?= $(TARGET_OS_LOCAL)

ifeq ($(origin DEBUG), undefined)
  BUILDTYPE_DIR:=release
else ifeq ($(DEBUG),0)
  BUILDTYPE_DIR:=release
else
  BUILDTYPE_DIR:=debug
  GCFLAGS:=-gcflags="all=-N -l"
  $(info $(H) Build with debugger information)
endif
ifeq ($(GOOS),windows)
GOLANGCI_LINT:=golangci-lint.exe
else
GOLANGCI_LINT:=golangci-lint
endif

################################################################################
# Target: tests                                                                #
################################################################################

test:
ifeq ($(GOOS), windows)
	@go test -v -cover -gcflags=all=-l .\...
else
	@go test -v -cover -gcflags=all=-l -coverprofile=coverage.out ./...
endif


################################################################################
# Target: lint                                                                 #
################################################################################
# Due to https://github.com/golangci/golangci-lint/issues/580, we need to add --fix for windows
.PHONY: lint
lint:
	$(GOLANGCI_LINT) run --timeout=20m


.PHONY: install generate

