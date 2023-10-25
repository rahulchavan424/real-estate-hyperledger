package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// Configuration information
var (
	sdk           *fabsdk.FabricSDK                              // Fabric SDK
	configPath    = "config-local-dev.yaml"                                // Configuration file path
	channelName   = "appchannel"                                 // Channel name
	user          = "Admin"                                      // User
	chainCodeName = "fabric-realty"                              // Chaincode name
	endpoints     = []string{"peer0.jd.com", "peer0.taobao.com"} // Nodes to send transactions to

	//configPath    = "config-local-dev.yaml"                      // Configuration file path (used for local development)
)

// Init initializes the blockchain SDK
func Init() {
	var err error
	// Initialize the SDK with the configuration file
	sdk, err = fabsdk.New(config.FromFile(configPath))
	if err != nil {
		panic(err)
	}
}

// ChannelExecute interacts with the blockchain
func ChannelExecute(fcn string, args [][]byte) (channel.Response, error) {
	// Create a client, indicating the identity on the channel
	ctx := sdk.ChannelContext(channelName, fabsdk.WithUser(user))
	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}
	// Write operation to the blockchain ledger (invokes the chaincode)
	resp, err := cli.Execute(channel.Request{
		ChaincodeID: chainCodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints(endpoints...))
	if err != nil {
		return channel.Response{}, err
	}
	// Return the result of the chaincode execution
	return resp, nil
}

// ChannelQuery performs a blockchain query
func ChannelQuery(fcn string, args [][]byte) (channel.Response, error) {
	// Create a client, indicating the identity on the channel
	ctx := sdk.ChannelContext(channelName, fabsdk.WithUser(user))
	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}
	// Query the blockchain ledger (invokes the chaincode) and only returns the result
	resp, err := cli.Query(channel.Request{
		ChaincodeID: chainCodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints(endpoints...))
	if err != nil {
		return channel.Response{}, err
	}
	// Return the result of the chaincode execution
	return resp, nil
}
