# Welcome to our community blockchain
## 1, project structure
> **blockchain-core**: core blockchain code written by GoLang


> **blockchain-web**: user GUI, written by spring boot plus bootstrap. used by end user who will interacted with

## 2, development environment for blockchain-core

### Eclipse + GoClipse

if you choose eclipse + goclipse to develop golang code, great! please keep reading.

* download Eclipse Oxygen V3
* open menu 'help' -> 'Eclipse Marketplace', search 'goclipse'
* install goclipse plugin, restart Eclipse
* checkout code by: git clone https://github.com/yingjiesun/BC-GoLang-Blockchain.git
* from Eclipse choose 'import' -> 'Existing Projects into Workspace'
* choose folder 'blockchain-core', click 'Open'
* if everything goes well, congradulations! you did it.

### Other IDE

* welcome to add steps


## 3, development environment for blockchain-web

### Eclipse

smart choice

* download Eclipse Oxygen
* checkout code by: git clone https://github.com/yingjiesun/BC-GoLang-Blockchain.git
* from Eclipse choose 'import' -> 'Existing Maven Projects'
* choose folder 'blockchain-web', click 'Open'

### Other IDE

* welcome to add steps

## 4, how to run blockchain

after server start, can visit **http://localhost:8881**

## 5, how to run blockchain-web

from folder 'blockchain-web' run command

`mvn clean spring-boot:run`

then can visit **http://localhost:8080**


