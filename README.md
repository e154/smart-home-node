Smart home node
---------------

[Project site](https://e154.github.io/smart-home/) |
[Configurator](https://github.com/e154/smart-home-configurator/) |
[Server](https://github.com/e154/smart-home/) |
[Development Tools](https://github.com/e154/smart-home-tools/) |
[Smart home Socket](https://github.com/e154/smart-home-socket/) |
[Modbus device controller](https://github.com/e154/smart-home-modbus-ctrl-v1/)

[![Build Status](https://travis-ci.org/e154/smart-home-node.svg?branch=master)](https://travis-ci.org/e154/smart-home-node)
![status](https://img.shields.io/badge/status-beta-yellow.svg)

Attention! The project is under active development.
---------

### Installation for development

access to serial port

sudo gpasswd --add ${USER} dialout
    
or
    
sudo usermod -a -G dialout ${USER}
    
You then need to log out and log back in again for it to be effective. 

```bash
go get -u github.com/golang/dep/cmd/dep

git clone https://github.com/e154/smart-home-node $GOPATH/src/github.com/e154/smart-home-node

cd $GOPATH/src/github.com/e154/smart-home-node

dep ensure

go build
```

### LICENSE

[MIT Public License](https://github.com/e154/smart-home-node/blob/master/LICENSE)