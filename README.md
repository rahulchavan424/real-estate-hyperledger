Manual Deployment

Prerequisites: You need a Linux or Mac OS environment with Docker and Docker Compose installed.

    Download this project and place it in any directory, for example: /root/fabric-realty.

    Give the project executable permissions by running sudo chmod -R +x /root/fabric-realty/.

    Go to the network directory and execute ./start.sh to deploy the blockchain network and smart contracts.

    Go to the application directory and execute ./start.sh to start the frontend and backend applications. You can access the frontend page at http://localhost:8000, and the backend API is available at http://localhost:8888.

    (Optional) If you want to set up the blockchain explorer, go to the network/explorer directory and run ./start.sh. You can access it at http://localhost:8080 with the username "admin" and password "123456."

Clean Environment

Please note that this operation will clear all data. Follow these steps in order:

    (If the blockchain explorer is running) Go to the network/explorer directory and run ./stop.sh to close the blockchain explorer.

    Go to the application directory and run ./stop.sh to shut down the blockchain application.

    Finally, go to the network directory and run ./stop.sh to shut down the blockchain network and clean the chaincode containers.

Directory Structure

    application/server: Uses fabric-sdk-go to invoke smart contracts, and gin provides external access interfaces (RESTful API).

    application/web: Combines Vue and Element-UI to provide the frontend presentation pages.

    chaincode: Contains smart contracts written in Go.

    network: Holds the Hyperledger Fabric blockchain network configuration.

Workflow

    The administrator creates real estate for user homeowners.

    Homeowners can view information about the properties they own.

    Homeowners can initiate sales, and anyone can view the sales listings. Buyers can make purchases, deduct funds, and wait for homeowner confirmation. After the transaction is complete, property ownership is updated. Transactions can be canceled at any time during the validity period, and they will automatically close after the expiration of the validity period.

    Homeowners can also initiate donations, specifying the recipients. Before the recipient confirms acceptance, both parties can cancel the donation.