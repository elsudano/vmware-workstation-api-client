SHELL = /bin/bash

NAME = vmware-workstation-api-client
VERSION = v0.0.1
BINARY = $(NAME)_$(VERSION)
PRIVATEKEYFILE = workstationapi-key.pem
PK_EXISTS=$(shell [ -e $(PRIVATEKEYFILE) ] && echo 1 || echo 0 )
CERTFILE = workstationapi-cert.pem
CRT_EXISTS=$(shell [ -e $(CERTFILE) ] && echo 1 || echo 0 )
IPADDRESS = 127.0.0.1
PORT = 8697
SENSIBLE_DATA_FILE = config.ini
CONFIG_FILE = ~/.vmrestCfg
F1_EXISTS=$(shell [ -e $(CONFIG_FILE) ] && echo 1 || echo 0 )

#-------------------------------------------------------#
#    Public Functions                                   #
#-------------------------------------------------------#
PHONY += help
help:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sort | awk 'BEGIN {FS = ":.*?## "}; \
	{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

start_api_rest: --generateSSL ## Prepare environment for you can use a API REST of VmWare Workstation Pro and generate files for SSL
ifeq ($(F1_EXISTS), 1)
	@vmrest -k $(PRIVATEKEYFILE) -c $(CERTFILE) -i $(IPADDRESS) -p $(PORT) > /dev/null &
else
	@cat $(SENSIBLE_DATA_FILE)
	@vmrest -C $(CONFIG_FILE)
	@vmrest -k $(PRIVATEKEYFILE) -c $(CERTFILE) -i $(IPADDRESS) -p $(PORT) > /dev/null &
endif

stop_api_rest: ## Stop a API REST server of VmWare Workstation
ifeq ($(shell ps aux | grep -v grep | grep -e vmrest -e $(PRIVATEKEYFILE) -e $(CERTFILE) | awk '{print $$2}' | tr -d ' \n\r'),)
	$(error PID not found, try to find manually)
else
	@kill -9 $(shell ps aux | grep -v grep | grep -e vmrest -e $(PRIVATEKEYFILE) -e $(CERTFILE) | awk '{print $$2}' | tr -d ' \n\r')
endif

PHONY += status_api_rest
status_api_rest: ## Check if the API is work ir not
ifeq ($(shell ps aux | grep -v grep | grep -e vmrest -e $(PRIVATEKEYFILE) -e $(CERTFILE) | awk '{print $$2}' | tr -d ' \n\r'),)
	$(error "The API server don't running")
else
	@echo -e "The API server it's running"
endif

test_api_client: ## Test API client and list all virtual machine of VmWare Workstation
	@go run .

build: ## Build the binary of the module
	@go get
	@go build -o $(BINARY)

clean: ## Clean the project, this only remove default config of API REST VmWare Workstation Pro, the cert, private key and binary
	@rm -f $(PRIVATEKEYFILE)
	@rm -f $(CERTFILE)
	@rm -f $(BINARY)
	@rm -f $(CONFIG_FILE)

#-------------------------------------------------------#
#    Private Functions                                  #
#-------------------------------------------------------#
--generateSSL:
ifeq ($(PK_EXISTS),0)
	@openssl req -x509 -newkey rsa:4096 -keyout $(PRIVATEKEYFILE) -out $(CERTFILE) -days 365 -nodes -subj "/C=ES/ST=Granada/L=Granada/O=Internet SL/OU=IT/CN=localhost"
endif
ifeq ($(CRT_EXISTS),0)
	@openssl req -x509 -newkey rsa:4096 -keyout $(PRIVATEKEYFILE) -out $(CERTFILE) -days 365 -nodes -subj "/C=ES/ST=Granada/L=Granada/O=Internet SL/OU=IT/CN=localhost"
endif

--findPID:
	$(eval PID := $(shell ps aux | grep -v grep | grep -e vmrest -e $(PRIVATEKEYFILE) -e $(CERTFILE) | awk '{print $$2}' | tr -d ' \n\r'))

.PHONY = $(PHONY)
