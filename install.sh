#!/bin/bash
wget -O amimoto https://github.com/amimoto-ami/amimoto-cli/releases/download/v0.0.1/amimoto-cli_linux_amd64
chmod +x amimoto
chown root:root amimoto
mv amimoto /usr/bin/
