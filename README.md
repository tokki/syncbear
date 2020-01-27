# syncbear

sidecar to sync user info and traffic 

## install

install golang	

make

scp ./syncbear to your server

## options

check --help

## run

./syncbear --help

## run shadowsocks

support [shadowsocks-libev](https://github.com/shadowsocksr-backup/shadowsocksr-libev)

ss-manager --manager-address 127.0.0.1:8080 -s your.domain(ip).com -u

## run v2ray

check the v2ray.json

dev.json for debug v2ray on local ,ignore it

run caddy 

check caddy.conf

put some index.html to /www/ folder to make there is real site

