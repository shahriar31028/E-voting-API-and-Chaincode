'use strict';

const { Gateway, Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const path = require('path');
const { buildCAClient, registerAndEnrollUser, enrollAdmin } = require('../../test-application/javascript/CAUtil.js');
const { buildCCPOrg1, buildWallet } = require('../../test-application/javascript/AppUtil.js');

const channelName = 'mychannel';
const chaincodeName = 'basic4';
const mspOrg1 = 'Org1MSP';
const walletPath = path.join(__dirname, 'wallet');
const org1UserId = 'appUser';

function prettyJSONString(inputString) {
	return JSON.stringify(JSON.parse(inputString), null, 2);
}

async function main() {
	try {
		
		const ccp = buildCCPOrg1();
		const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org1.example.com');
		const wallet = await buildWallet(Wallets, walletPath);
		await enrollAdmin(caClient, wallet, mspOrg1);
		await registerAndEnrollUser(caClient, wallet, mspOrg1, org1UserId, 'org1.department1');
		const gateway = new Gateway();
		try {
			await gateway.connect(ccp, {
				wallet,
				identity: org1UserId,
				discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
			});

			const network = await gateway.getNetwork(channelName);

			const contract = network.getContract(chaincodeName);

			// console.log('\n--> Submit Transaction: RegisterUser, registers new asset with ID, name, email, password');
			// await contract.submitTransaction('RegisterUser', 'user1', 'shahriar', 'shahriar2069@gmail.com', 'abcd1234');
			// console.log('*** Result: committed');

			// console.log('\n--> Evaluate Transaction: SayHello');
			// let result = await contract.evaluateTransaction('SayHello');
			// console.log(`*** Result: ${prettyJSONString(result.toString())}`);

            return {
                query: contract.evaluateTransaction,
                invoke: contract.submitTransaction
            }
		} finally {
			// gateway.disconnect();
		}
	} catch (error) {
		console.error(`******** FAILED to run the application: ${error}`);
	}
}

module.exports = main