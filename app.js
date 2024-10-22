const express = require('express');
const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');

const app = express();
app.use(express.json());

const ccpPath = path.resolve(__dirname, 'connection-org1.json');

async function connectToNetwork() {
    const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);

    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: 'admin',
        discovery: { enabled: true, asLocalhost: true }
    });

    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('basic');
    return contract;
}

app.post('/createAccount', async (req, res) => {
    try {
        const { dealerID, msisdn, mpin, balance, status } = req.body;
        const contract = await connectToNetwork();
        await contract.submitTransaction('CreateAccount', dealerID, msisdn, mpin, balance, status);
        res.status(200).send('Account created successfully');
    } catch (error) {
        res.status(500).send(`Error: ${error.message}`);
    }
});

app.get('/readAccount/:dealerID', async (req, res) => {
    try {
        const dealerID = req.params.dealerID;
        const contract = await connectToNetwork();
        const result = await contract.evaluateTransaction('ReadAccount', dealerID);
        res.status(200).json(JSON.parse(result));
    } catch (error) {
        res.status(500).send(`Error: ${error.message}`);
    }
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
});
