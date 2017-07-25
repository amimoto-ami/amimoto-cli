# AMIMOTO-CLI

#### Currently in dev. Do not use in production.

### Usage and Installation

##### 1. Download AMIMOTO-CLI
`wget https://github.com/amimoto-ami/amimoto-cli/raw/master/amimoto`

##### 2. Make AMIMOTO-CLI executable
`chmod +x amimoto`

##### 3. Move to a globally available location
`sudo mv amimoto /usr/bin/`

#### Or this one liner

```
curl -L -s https://github.com/amimoto-ami/amimoto-cli/raw/master/install.sh | sudo /bin/bash
```

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

#### Git Clone

`git clone git@github.com:amimoto-ami/amimoto-cli.git`

or

`git clone https://github.com/amimoto-ami/amimoto-cli.git`

#### Build

`go build -o amimoto`
