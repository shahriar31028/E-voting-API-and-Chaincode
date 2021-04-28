"use strict";

const { Gateway, Wallets } = require("fabric-network");
const FabricCAServices = require("fabric-ca-client");
const path = require("path");
const {
  buildCAClient,
  registerAndEnrollUser,
  enrollAdmin,
} = require("../../test-application/javascript/CAUtil.js");
const {
  buildCCPOrg1,
  buildWallet,
} = require("../../test-application/javascript/AppUtil.js");

const channelName = "mychannel";
const chaincodeName = "basic";
const mspOrg1 = "Org1MSP";
const walletPath = path.join(__dirname, "wallet");
const org1UserId = "appUser";

function prettyJSONString(inputString) {
  return JSON.stringify(JSON.parse(inputString), null, 2);
}

async function main() {
  try {
    const ccp = buildCCPOrg1();
    const caClient = buildCAClient(
      FabricCAServices,
      ccp,
      "ca.org1.example.com"
    );
    const wallet = await buildWallet(Wallets, walletPath);
    await enrollAdmin(caClient, wallet, mspOrg1);
    await registerAndEnrollUser(
      caClient,
      wallet,
      mspOrg1,
      org1UserId,
      "org1.department1"
    );
    const gateway = new Gateway();
    try {
      await gateway.connect(ccp, {
        wallet,
        identity: org1UserId,
        discovery: { enabled: true, asLocalhost: true },
      });

      const network = await gateway.getNetwork(channelName);
      const contract = network.getContract(chaincodeName);

      /////////////////////////////////////////////////////////////
      const express = require("express");
      const cookieParser = require("cookie-parser");
      var bodyParser = require("body-parser");
      var cors = require("cors");

      const app = express();
      const port = 3000;

      app.use(cookieParser());
      app.use(express.urlencoded({ extended: false }));
      app.use(express.json());
      app.use(cors({ credentials: true, origin: "http://localhost:3001" }));

      app.get("/", async (req, res) => {
        let data = await contract.evaluateTransaction("SayHello");

        res.send(data);
      });

      app.post("/registerUser", async function (req, res) {
        let { name, email, password } = req.body;
        const id = `user_${email}`;

        try {
          await contract.submitTransaction(
            "RegisterUser",
            id,
            name,
            email,
            password
          );
          res.json({
            status: "register user successful",
          });
        } catch (error) {
          res
            .status(400)
            .send(`******** FAILED to run the application: ${error}`);
        }
      });

      app.post("/loginUser", async function (req, res) {
        let { email, password } = req.body;

        try {
          let queryResult = await contract.evaluateTransaction(
            "LoginUser",
            email,
            password
          );

          let users = JSON.parse(queryResult.toString());

          if (users.length === 0) {
            throw "Email or Password incorrect";
          }

          let user = users[0];

          res.cookie("user", JSON.stringify(user), {
            maxAge: 3600_000,
            httpOnly: true,
          });

          res.json(user);
        } catch (error) {
          res.status(400).json({
            error: error.toString(),
          });
        }
      });

      app.get("/logoutUser", async function (req, res) {
        const { email, password } = req.body;

        try {
          // let queryResult = await contract.evaluateTransaction(
          //   "LoginUser",
          //   email,
          //   password
          // );

          res.cookie("user", null, { maxAge: -1, httpOnly: true });
          res.send("You have successfully Logged out");
        } catch (error) {
          res.status(400).send(`logout failed`);
        }
      });

      app.post("/showallCandidate", async function (req, res) {
        let { electionid } = req.body;

        try {
          let queryResult = await contract.evaluateTransaction(
            "ShowAllCandidates",
            electionid
          );

          res.send(queryResult);
        } catch (error) {
          res.status(400).send(`Failed`);
        }
      });

      app.post("/showallElections", async function (req, res) {
        let { doctype } = req.body;

        console.log(req.cookies);

        try {
          let queryResult = await contract.evaluateTransaction(
            "ShowAllElections",
            doctype
          );

          res.send(queryResult);
        } catch (error) {
          res.status(400).send(`Failed`);
        }
      });

      app.post("/createElection", async function (req, res) {
        let { id, name } = req.body;

        try {
          let queryResult = await contract.submitTransaction(
            "CreateElection",
            id,
            name
          );

          res.send("Election Created");
        } catch (error) {
          res.status(400).send(`error: ${error}`);
        }
      });

      app.post("/addCandidate", async function (req, res) {
        let { id, name, marka, electionid } = req.body;

        //console.log(req.body);

        try {
          let queryResult = await contract.submitTransaction(
            "AddCandidate",
            id,
            name,
            marka,
            electionid
          );

          res.send("Candidate Created");
        } catch (error) {
          res.status(400).send(`error: ${error}`);
        }
      });

      //////////////have to test//////////
      app.post("/votecasting", async function (req, res) {
        let { electionid, candidateid } = req.body;

        let user = JSON.parse(req.cookies.user);

        const id = `vote_${electionid}_${user.ID}`;

        try {
          let queryResult = await contract.submitTransaction(
            "VoteCasting",
            id,
            electionid,
            candidateid
          );

          res.json({ status: "Vote Casted" });
        } catch (error) {
          res.status(400).send(`error: ${error}`);
        }
      });

      ////////////
      app.post("/stopelection", async function (req, res) {
        let { id } = req.body;

        try {
          let queryResult = await contract.submitTransaction(
            "StopElection",
            id
          );

          res.send("Election stopped");
        } catch (error) {
          res.status(400).send(`error: ${error}`);
        }
      });

      app.post("/calculateResult", async function (req, res) {
        let { electionid } = req.body;

        try {
          let queryResult = await contract.evaluateTransaction(
            "CalculateResult",
            electionid
          );

          res.send(queryResult);
        } catch (error) {
          res.status(400).send(`error: ${error}`);
        }
      });

      app.listen(port, async () => {
        console.log(`Example app listening at http://localhost:${port}`);
      });
      /////////////////////////////////////////////////////////////
    } finally {
      //   console.log("disconnecting from chaincode");
      //   gateway.disconnect();
    }
  } catch (error) {
    console.error(`******** FAILED to run the application: ${error}`);
  }
}

main();
