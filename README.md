Smart home node
---------------

[Project site](https://e154.github.io/smart-home/) |
[Configurator](https://github.com/e154/smart-home-configurator/) |
[Server](https://github.com/e154/smart-home/) |
[Development Tools](https://github.com/e154/smart-home-tools/) |
[Smart home Socket](https://github.com/e154/smart-home-socket/)

[![Build Status](https://travis-ci.org/e154/smart-home-node.svg?branch=master)](https://travis-ci.org/e154/smart-home-node)
[![Coverage Status](https://coveralls.io/repos/github/e154/smart-home-node/badge.svg?branch=cover)](https://coveralls.io/github/e154/smart-home-node?branch=cover)

##### Installation

access to serial port

sudo gpasswd --add ${USER} dialout
    
or
    
sudo usermod -a -G dialout ${USER}
    
You then need to log out and log back in again for it to be effective. 

##### Error codes
    
    1 serial port errors 
    2 modbus line errors
    3 tcp read bytes errors
    4 unmarshal bytes to json from tcp errors

##### TODO

* работа в качестве демона https://github.com/takama/daemon
* доступ по сертификату
* shell console ?

##### Протокол основанный но modbus

* ASCII
* проверка целостности пакета по контрольной сумме LRC
* ограничение по времени ожидания ответа 2сек

### LICENSE

[MIT Public License](https://github.com/e154/smart-home-node/blob/master/LICENSE)