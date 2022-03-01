[string]${NAME} = "vmware-workstation-api-client"
[string]${DIRELEASES} = "releases/"
[string]${PRIVATEKEYFILE} = "workstationapi-key.pem"
[string]${CERTFILE} = "workstationapi-cert.pem"
[string]${CONFIG_FILE} = "~/.vmrestCfg"



Function help {
    [array]${HELPS} = Select-String -Path .\Makefile.ps1 -Pattern "^(Function?) (?<fun>\w+) (.+##) (?<help>.+)"
    ${HELPS} | ForEach-Object {
        ${var}, ${help} = $_.Matches[0].Groups['fun','help'].Value
        Write-Host "${var} : `t${help}"
    }
}

Function clean { ## Clean the project, this only remove default config of API REST VmWare Workstation Pro, the cert, private key and binary
    if ( Test-Path -Path ${DIRELEASES}/${global:BINARY} ) {
        Remove-Item ${DIRELEASES}/${global:BINARY}
    }
    Write-Host "we have did Cleaning"
}
 
Function prepare_environment {
    [string]${SENSIBLE_DATA_FILE} = ".\config.ini"
    [array]${OPTIONS} = Select-String -Path $SENSIBLE_DATA_FILE -Pattern "^(?<key>\w+)=(?<value>.+$)"
    ${OPTIONS} | ForEach-Object {
        ${key}, ${value} = $_.Matches[0].Groups['key','value'].Value
        Set-Variable -Name ${key} -Value ${value}
    }
    [string]${VERSION} = Select-String -Path ".\wsapiclient\wsapiclient.go" -Pattern "\tlibraryVersion"
    [string]${global:BINARY} = ${NAME}+"_v"+${VERSION}.Split("=")[1]+".exe"
    ${global:BINARY} = ${global:BINARY} -replace '"',''
    ${global:BINARY} = ${global:BINARY} -replace ' ',''
}

Function build { ## Build the binary of the module
    prepare_environment
    if ( -Not (Test-Path -Path ${DIRELEASES}) ) {
        New-Item ${DIRELEASES} -itemtype directory | Out-Null
    }
    & go get
    & go build -o ${DIRELEASES}/${global:BINARY}
    Write-Host "we made the binary"
}

Function api_test { ## With this command we can test our API Rest
    & go get
    & go run .
    Write-Host "we are testing"
}
Function api_status { ## With this command we can see ths status of the API Rest of VmWare WS

}

Function api_start { ## We can start the API Rest of VmWare WS
    # "C:\Program Files (x86)\VMware\VMware Workstation\vmrest.exe" -c workstationapi-cert.pem -k workstationapi-key.pem
}

Function api_stop { ## We can stop the API Rest of VmWare WS

}

# --------------------------
#       Menu
# --------------------------
if ($args.Count -eq 1 ) {
    switch ( $args[0] ) {
        clean { clean }
        build { build }
        api_test { api_test }
        api_status { api_status }
        api_start { api_start }
        api_stop { api_stop }
    }
} else {
    help
}
