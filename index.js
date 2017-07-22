#!/usr/bin/env node

var program = require('commander');
var exec = require('exec');
var child_process = require('child_process');

program
  .version('0.0.1')
  .description('CLI tool for managing AMIMOTO-AMI')
  .command('', 'Display amimoto-cli help')

program
  .command('cache')
  .option('-d, --delete', 'Clear NGINX proxy cache')
  .action(function() {
    child_process.execFile('rm', ['-rf', '/var/cache/nginx/proxy_cache/*'], function(error, stdout, stderr){
      console.log(stdout);
    });
  });

program.parse(process.argv);
