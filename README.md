Smart home node
---------------

[Project site](https://e154.github.io/smart-home/) |
[Server](https://github.com/e154/smart-home/)

![status](https://img.shields.io/badge/status-beta-yellow.svg)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)


|Branch      |Status   |
|------------|---------|
|master      | ![Build Status](https://github.com/e154/smart-home-node/actions/workflows/test.yml/badge.svg?branch=master)  |
|dev         | ![Build Status](https://github.com/e154/smart-home-node/actions/workflows/test.yml/badge.svg?branch=develop) |


Attention! The project is under active development.
---------

### Installation for development

access to serial port

sudo gpasswd --add ${USER} dialout
    
or
    
sudo usermod -a -G dialout ${USER}
    
You then need to log out and log back in again for it to be effective. 

```bash
git clone https://github.com/e154/smart-home-node $GOPATH/src/github.com/e154/smart-home-node

cd $GOPATH/src/github.com/e154/smart-home-node

go mod vendor

go build
```

### LICENSE

[GPLv3 Public License](https://github.com/e154/smart-home-node/blob/master/LICENSE)
