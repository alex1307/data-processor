# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
APP_FOLDER = apps
RELEASE_DIR = release
DEST_DIR = /Users/$(USER)/Software/$(RELEASE_DIR)/GO/bin
RESOURCE_DIR = resources
CONFIG_DIR = $(RESOURCE_DIR)/config
DEST_CONFIG_DIR = /Users/$(USER)/Software/$(RELEASE_DIR)/GO/config
APP_DIR = /Users/$(USER)/Software/$(RELEASE_DIR)/GO/apps/data-processor
# Build target


all: clean install copy-apps

install:
	go install ./cmd/processor
	strip $(DEST_DIR)/processor
clean:
	go clean

copy-apps:
	mkdir -p $(DEST_CONFIG_DIR)
	mkdir -p $(APP_DIR)
	cp -r $(DEST_DIR)/processor $(APP_DIR)/
	cp -r $(RESOURCE_DIR) $(APP_DIR)/
	rm -r $(APP_DIR)/$(RESOURCE_DIR)/test
	cp -r $(CONFIG_DIR)/* $(DEST_CONFIG_DIR)/


	
