# AMIMOTO-CLI

#### Currently in dev. Do not use in production.

### Usage and Installation

- `wget https://github.com/amimoto-ami/go-amimoto-cli/raw/master/amimoto`
- `chmod +x amimoto`
- `sudo mv amimoto /usr/bin/`

#### Example

##### Clear NGINX proxy cache
`sudo amimoto cache --purge`

##### Add virtual host example.com
`sudo amimoto add example.com`
