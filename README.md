# AMIMOTO-CLI

#### Beta release. Use at your own risk.

### Usage and Installation

##### 1. Download AMIMOTO-CLI
`wget -O amimoto https://github.com/amimoto-ami/amimoto-cli/releases/download/v0.0.1/amimoto-cli_linux_amd64`

##### 2. Make AMIMOTO-CLI executable
`chmod +x amimoto`

##### 3. Move to a globally available location
`sudo mv amimoto /usr/bin/`

#### Or this one liner

`curl -L -s https://github.com/amimoto-ami/amimoto-cli/raw/master/install.sh | sudo /bin/bash`

#### Examples

##### Clear NGINX proxy cache
`sudo amimoto cache --purge`

##### Add virtual host example.com
`sudo amimoto site --add example.com`

##### Disable virtual host example.com
`sudo amimoto site --disable example.com`

##### Enable virtual host example.com
`sudo amimoto site --enable example.com`

##### Remove virtual host example.com
`sudo amimoto site --remove example.com`

### Developing New Features

## Requirements

- The [Go](https://github.com/golang/go) Programming Language

#### Download

Get binary from here.
- https://github.com/amimoto-ami/amimoto-cli/releases


#### Git Clone

`git clone git@github.com:amimoto-ami/amimoto-cli.git`

or

`git clone https://github.com/amimoto-ami/amimoto-cli.git`

#### Build

```
$ go get github.com/cloudbuy/go-pkg-optarg
$ go get github.com/koron/go-dproxy
$ go get github.com/go-sql-driver/mysql
```

Build single binary for local os.
```
$ go build amimoto.go
```

Build for multi os(linux 386, amd64).
```
$ go get github.com/mitchellh/gox
$ gox -os="linux" -arch="386 amd64" -output "build/amimoto-cli_{{.OS}}_{{.Arch}}"
```

upload releases to github. (for maintaner information)
```
$ go get github.com/tcnksm/ghr
$ ghr --replace -u amimoto-ami v0.0.1 build/
```
