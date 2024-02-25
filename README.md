# API client to VmWare Workstation

In this project we created a basic client to interact with the VmWare Workstation API in its version of Linux

# Use quickly:

For the fastest users, this is the way to start:

* Clone the repository

```bash
git clone https://github.com/elsudano/vmware-workstation-api-client.git
mv config.example config.ini
```

* Set the necessary values inside the config.ini file

```bash
User: 
Password: 
ParentId: F76OJJ66M052TJE435D3AB69ORP6SAFR
BaseURL: https://localhost:8697/api
Debug: false
```

* For normal use, We have created a Makefile with some commands to help us, with `make` show in screen the usage help.

```bash
$ make
api_start                      Prepare environment for you can use a API REST of VmWare Workstation Pro and generate files for SSL
api_status                     Check if the API is work ir not
api_stop                       Stop a API REST server of VmWare Workstation
api_test                       Test API client and list all virtual machine of VmWare Workstation
build                          Build the binary of the module
clean                          Clean the project, this only remove default config of API REST VmWare Workstation Pro, the cert, private key and binary
publish                        Build and Publish a new TAG in GitHub
```

* The follow step it's run the next command to create a SSL certificates and start the API with the values that you setting in config.ini

```bash
$ make start_api_rest
Generating a RSA private key
..........++++
writing new private key to 'workstationapi-key.pem'
-----
User: xxxxxxx
Password: xxxxxxx
BaseURL: https://localhost:8697/api
Debug: false
VMware Workstation REST API
Copyright (C) 2018-2020 VMware Inc.
All Rights Reserved

vmrest 1.2.1 build-16341506
Username:xxxxxxx
New password:xxxxxxx
Retype new password:
Processing...
Credential updated successfully
```

* Finally to testing if the client works, run the follow command:

```bash
$ make api_test
```

* The reult of the previous command it's a list of vms that you have in the VmWare Workstation

#### 1. Pre-Requisites

VmWare Workstation 14 or higher is required