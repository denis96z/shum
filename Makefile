PWD := $(shell pwd)

CMD_PATH := $(PWD)/cmd/shum
PKG_PATH := $(PWD)/pkg/shum

MAIN_PATH := $(CMD_PATH)/main.go
TARGET_PATH := $(PWD)/bin/shum

.PHONY: all
all: shum

.PHONY: fmt
fmt:
	go fmt $(PWD)/...

.PHONY: dep
dep:
	go mod tidy && go mod vendor && go mod verify

.PHONY: shum
shum:
	go build -o $(TARGET_PATH) $(MAIN_PATH)

CONTRIB_PATH := $(PWD)/contrib

CONTRIB_CONF_PATH := $(CONTRIB_PATH)/shum.yml
CONTRIB_SERVICE_PATH := $(CONTRIB_PATH)/shum.service

BIN_PATH := /usr/bin/shum
CONF_PATH := /etc/shum.yml
SERVICE_PATH := /etc/systemd/system/shum.service

.PHONY: install
install:
	cp $(TARGET_PATH) $(BIN_PATH)
	cp $(CONTRIB_CONF_PATH) $(CONF_PATH)
	cp $(CONTRIB_SERVICE_PATH) $(SERVICE_PATH)

.PHONY: remove
remove:
	rm -f $(SERVICE_PATH)
	rm -f $(CONF_PATH)
	rm -f $(BIN_PATH)
