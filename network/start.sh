#!/bin/bash

# Check the operating system type
if [[ `uname` == 'Darwin' ]]; then
  echo "The current operating system is macOS"
  export PATH=${PWD}/hyperledger-fabric-darwin-amd64-1.4.12/bin:$PATH
elif [[ `uname` == 'Linux' ]]; then
  echo "The current operating system is Linux"
  export PATH=${PWD}/hyperledger-fabric-linux-amd64-1.4.12/bin:$PATH
else
  echo "The current operating system is neither macOS nor Linux. The script cannot continue!"
  exit 1
fi

echo -e "Note: If you have previously deployed the network, running this script will result in data loss!\nIf you just want to restart the network, you can execute the 'docker-compose restart' command."
read -p "Are you sure you want to continue? Please enter Y or y to proceed: " confirm

if [[ "$confirm" != "Y" && "$confirm" != "y" ]]; then
  echo "You have canceled the execution of the script."
  exit 1
fi

echo "1. Clean the environment"
./stop.sh

echo "2. Generate certificates and keys (MSP materials); the results will be saved in the crypto-config folder"
cryptogen generate --config=./crypto-config.yaml

echo "3. Create the orderer genesis block"
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./config/genesis.block -channelID firstchannel

echo "4. Generate the channel configuration transaction 'appchannel.tx'"
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./config/appchannel.tx -channelID appchannel

echo "5. Define an anchor node for Taobao"
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./config/TaobaoAnchor.tx -channelID appchannel -asOrg Taobao

echo "6. Define an anchor node for JD"
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./config/JDAnchor.tx -channelID appchannel -asOrg JD

echo "Blockchain: Start"
docker-compose up -d
echo "Waiting for the nodes to start, waiting for 10 seconds"
sleep 10

TaobaoPeer0Cli="CORE_PEER_ADDRESS=peer0.taobao.com:7051 CORE_PEER_LOCALMSPID=TaobaoMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/taobao.com/users/Admin@taobao.com/msp"
TaobaoPeer1Cli="CORE_PEER_ADDRESS=peer1.taobao.com:7051 CORE_PEER_LOCALMSPID=TaobaoMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/taobao.com/users/Admin@taobao.com/msp"
JDPeer0Cli="CORE_PEER_ADDRESS=peer0.jd.com:7051 CORE_PEER_LOCALMSPID=JDMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/jd.com/users/Admin@jd.com/msp"
JDPeer1Cli="CORE_PEER_ADDRESS=peer1.jd.com:7051 CORE_PEER_LOCALMSPID=JDMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/jd.com/users/Admin@jd.com/msp"

echo "7. Create a channel"
docker exec cli bash -c "$TaobaoPeer0Cli peer channel create -o orderer.qq.com:7050 -c appchannel -f /etc/hyperledger/config/appchannel.tx"

echo "8. Join all nodes to the channel"
docker exec cli bash -c "$TaobaoPeer0Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$TaobaoPeer1Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$JDPeer0Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$JDPeer1Cli peer channel join -b appchannel.block"

echo "9. Update anchor nodes"
docker exec cli bash -c "$TaobaoPeer0Cli peer channel update -o orderer.qq.com:7050 -c appchannel -f /etc/hyperledger/config/TaobaoAnchor.tx"
docker exec cli bash -c "$JDPeer0Cli peer channel update -o orderer.qq.com:7050 -c appchannel -f /etc/hyperledger/config/JDAnchor.tx"

# -n Chaincode name, can be set as needed
# -v Version
# -p Chaincode directory, in the /opt/gopath/src/ directory
echo "10. Install chaincode"
docker exec cli bash -c "$TaobaoPeer0Cli peer chaincode install -n fabric-realty -v 1.0.0 -l golang -p chaincode"
docker exec cli bash -c "$JDPeer0Cli peer chaincode install -n fabric-realty -v 1.0.0 -l golang -p chaincode"

# Only one node needs to instantiate the chaincode
# -n corresponds to the name of the chaincode installed in the previous step
# -v Version
# -C is the channel; in the fabric world, a channel is a different chain
# -c is for passing parameters, passing the init parameter
echo "11. Instantiate chaincode"
docker exec cli bash -c "$TaobaoPeer0Cli peer chaincode instantiate -o orderer.qq.com:7050 -C appchannel -n fabric-realty -l golang -v 1.0.0 -c '{\"Args\":[\"init\"]}' -P \"AND ('TaobaoMSP.member','JDMSP.member')\""

echo "Waiting for chaincode instantiation to complete, waiting for 5 seconds"
sleep 5

# Interact with the chaincode to verify if it is correctly installed and if the blockchain network is working
echo "12. Verify the chaincode"
docker exec cli bash -c "$TaobaoPeer0Cli peer chaincode invoke -C appchannel -n fabric-realty -c '{\"Args\":[\"hello\"]}'"

if docker exec cli bash -c "$JDPeer0Cli peer chaincode invoke -C appchannel -n fabric-realty -c '{\"Args\":[\"hello\"]}'" 2>&1 | grep "Chaincode invoke successful"; then
  echo "Congratulations! The network deployment was successful. If you need to temporarily stop the network, you can execute the 'docker-compose stop' command (data will not be lost)."
  exit 0
fi

echo "Warning: Network deployment was not successful. Please check each step to identify the specific issue."
