# Flow Swiss CLI

This command line interface serves as an additional frontend for the [Flow Swiss](https://my.flow.swiss) and [Cloudbit](https://my.cloudbit.ch) Platforms.

![Build](https://github.com/flowswiss/cli/workflows/Build/badge.svg)

## Installation

If you have GoLang installed, you can download and install the CLI with 
```shell script
go get github.com/flowswiss/cli/cmd/flow
```
otherwise, you will need to download the executable for your system from the release section in the github repository.

## Usage

After downloading you first of all need to authenticate the cli with your username and password.
**Warning**: those credentials will be stored in `$HOME/.flow/credentials.json`  
```shell script
flow auth login --username 'USERNAME' --password 'PASSWORD'
```
alternatively you can also pass `--username USERNAME` and `--password PASSWORD` to every other command or set the environment variables `FLOW_USERNAME` and `FLOW_PASSWORD` to avoid the credentials getting stored in your home directory.

Once you have successfully logged in into your account, you can start manipulating things in your organization. As a first step it would be a good idea to upload your personal ssh key onto our platform. You will need this for every linux virtual machine you deploy. 
```shell script
flow compute key-pair create \
    --name 'My first key pair' \
    --public-key ~/.ssh/id_rsa.pub
```

Just to test things out, you can try creating an ubuntu virtual machine using the previously uploaded key pair:
```shell script
flow compute server create \
    --name 'My first virtual machine' \
    --location 'ALP1' \
    --image 'ubuntu-20.04' \
    --product 'b1.1x1' \
    --key-pair 'My first key pair'
```

Further usage manuals can be found in the application itself using the `-h` or `--help` flags.