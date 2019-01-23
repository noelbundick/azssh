# azssh

Connect to [Azure Cloud Shell](https://docs.microsoft.com/en-us/azure/cloud-shell/overview) from your terminal

[![Build Status](https://dev.azure.com/noelbundick/noelbundick/_apis/build/status/azssh?branchName=master)](https://dev.azure.com/noelbundick/noelbundick/_build/latest?definitionId=27?branchName=master)

## Usage

* Launch with `azssh`
* Exit by typing `exit`

> Note: This app uses the clientId from the [vscode-azure-account](https://github.com/microsoft/vscode-azure-account) Visual Studio Code extension in order to call the necessary APIs. You will be prompted to allow access to "Visual Studio Code". This is expected behavior.

## Development

```
go get -u github.com/noelbundick/azssh
cd $GOPATH/src/github.com/noelbundick/azssh
go build
```