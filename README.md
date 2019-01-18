# azssh

Connect to [Azure Cloud Shell](https://docs.microsoft.com/en-us/azure/cloud-shell/overview) from your terminal

## Usage

* Launch with `azssh`
* Exit with `Ctrl+C`

> Note: This app uses the clientId from the [vscode-azure-account](https://github.com/microsoft/vscode-azure-account) Visual Studio Code extension in order to call the necessary APIs. You will be prompted to allow access to "Visual Studio Code". This is expected behavior.

## Development

```
go get -u github.com/noelbundick/azssh
cd $GOPATH/src/github.com/noelbundick/azssh
go build
```