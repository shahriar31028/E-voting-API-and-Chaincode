# E-voting-API-and-Chaincode

This project is built on Hyperledger Fabric for a voting system.

# Prerequisites: (Git,Curl, Docker, Docker-compose)

## installing Git and Curl
    sudo apt install -y git curl
    
## Installing Node.JS
    #install nvm (node version manager)
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.38.0/install.sh | bash

    source ~/.bashrc

    #install lts version of node
    nvm install --lts
    nvm use --lts
    nvm alias default --lts

## Installing Docker 
To install docker on your ubuntu machine, follow the steps from the tutorials given below:

Open (https://docs.docker.com/engine/install/ubuntu/)

## Installing Docker Compose  
To install docker compose on your ubuntu machine, follow the steps from the tutorials given below:

Open (https://docs.docker.com/compose/install/)

## Installing Hyperledger Fabric:

    cd $HOME
    curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/master/scripts/bootstrap.sh | bash -s
    
## Setting up Environmental Variables:

    echo export PATH=\$PATH:\$HOME/fabric-samples/bin | tee -a ~/.bashrc

    echo export FABRIC_CFG_PATH=$HOME/fabric-samples/config | tee -a ~/.bashrc

    source ~/.bashrc
    
## Downloading E-voting-API-and-Chaincode
    
    cd $HOME/fabric-samples

    git clone https://github.com/shahriar31028/E-voting-API-and-Chaincode.git 
    
    cd $HOME/fabric-samples/E-voting-API-and-Chaincode/application-javascript
    
    npm install
    npm install --save express
    npm install body-parser --save
    npm i cookie-parser --save
    npm install cors
    
   
## Installing GO in the system 
open (https://golang.org/doc/install)

## Starting Blockchain Test Network and install E-voting-API-and-Chaincode CC
    
    cd $HOME/fabric-samples/test-network

    # Start Test Network
    ./network.sh down && ./network.sh createChannel -ca -c mychannel -s couchdb

    # Install Chaincode
    ./network.sh deployCC -ccn basic -ccp ~/fabric-samples/E-voting-API-and-Chaincode/chaincode-go/ -ccl go
    
## Starting API

    cd $HOME/fabric-samples/E-voting-API-and-Chaincode/application-javascript

    npm install
    
    # Deleting the existing wallet from previous test network
    rm -rf wallet 
    
    #installing nodemon  
    npm install --save nodemon
    
    npx nodemon index.js
    
# Testing API
### install postman
 
    sudo snap install postman   

You can test the API with POSTMAN passing through the key and value,hitting the https://localhost/3000/*


## Viewing Blockchain state in CouchDB

You can view the current state at http://localhost:5984/_utils/.

#### Username : admin
#### Password : adminpw

## Stopping test network
    
    cd $HOME/fabric-samples/test-network
    ./network.sh down
