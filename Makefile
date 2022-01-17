TOPDIR=.

include $(TOPDIR)/Makefile.common.includes

export VERSION=$(Version)

BINPATH=bin$(SLASH)
SRCFiles=$(wildcard *$(SLASH)*.go *.go)
COVERAGE=cov.out
BINNAME=bratconverter

setup: go_get

srcFile:
	echo "SRCFiles ... $(SRCFiles)"

go_get: go.mod
	echo "Setting up go modules"
	$(GO) get -v

all: setup test productionbuild

productionbuild: productionbuildWindows productionbuildLinux productionbuildMac productionbuildMacM1

productionbuildWindows:
	$(eval BIN_OS:=$(GOOSWIN))
	$(eval ARCH_OS=$(GOARCHx64))
	$(eval OUTEXT=$(WINEXT))
	$(eval OUTBIN=$(BINNAME).$(BIN_OS).$(ARCH_OS)$(OUTEXT))
	@echo "building binary for $(BIN_OS) and $(ARCH_OS)"
	CGO_ENABLED=0 GOOS=$(BIN_OS) GOARCH=$(ARCH_OS) $(GO) build -v -tags "production" -ldflags="-s -w -X 'main.Version=$(VERSION)'" -o $(BINPATH)$(OUTBIN)

productionbuildMac:
	$(eval BIN_OS=$(GOOSMAC))
	$(eval ARCH_OS=$(GOARCHx64))
	$(eval OUTEXT=$(MACEXT))
	$(eval OUTBIN=$(BINNAME).$(BIN_OS).$(ARCH_OS)$(OUTEXT))
	@echo "building binary for $(BIN_OS) and $(ARCH_OS)"
	CGO_ENABLED=0 GOOS=$(BIN_OS) GOARCH=$(ARCH_OS) $(GO) build -v -tags "production" -ldflags="-s -w -X 'main.Version=$(VERSION)'" -o $(BINPATH)$(OUTBIN)

productionbuildMacM1:
	$(eval BIN_OS=$(GOOSMAC))
	$(eval ARCH_OS=$(GOARCHM1))
	$(eval OUTEXT=$(MACEXT))
	$(eval OUTBIN=$(BINNAME).$(BIN_OS).$(ARCH_OS)$(OUTEXT))
	@echo "building binary for $(BIN_OS) and $(ARCH_OS)"
	CGO_ENABLED=0 GOOS=$(BIN_OS) GOARCH=$(ARCH_OS) $(GO) build -v -tags "production" -ldflags="-s -w -X 'main.Version=$(VERSION)'" -o $(BINPATH)$(OUTBIN)

productionbuildLinux:
	$(eval BIN_OS=$(GOOSLINUX))
	$(eval ARCH_OS=$(GOARCHx64))
	$(eval OUTEXT=$(LINEXT))
	$(eval OUTBIN=$(BINNAME).$(BIN_OS).$(ARCH_OS)$(OUTEXT))
	@echo "building binary for $(BIN_OS) and $(ARCH_OS)"
	CGO_ENABLED=0 GOOS=$(BIN_OS) GOARCH=$(ARCH_OS) $(GO) build -v -tags "production" -ldflags="-s -w -X 'main.Version=$(VERSION)'" -o $(BINPATH)$(OUTBIN)

test: unittest 

unittest: $(COVERAGE) 

$(COVERAGE): $(SRCFiles) 
	echo "initiating test..."
	$(GO) test -p 1 -tags "development" -coverprofile=$(COVERAGE) .$(SLASH)...

clean: packClean
	@-$(RM) $(COVERAGE)
	@-$(RM) $(BINPATH)

