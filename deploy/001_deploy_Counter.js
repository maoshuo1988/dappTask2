const { deployments, getNamedAccounts, upgrades, ethers } = require("hardhat");

const fs = require("fs");
const path = require("path");

module.exports = async function ({ getNamedAccounts, deployments }) {
  const { deploy, save } = deployments;
  const { deployer } = await getNamedAccounts();

  console.log("Deploying Counter with account:", deployer);
  const CounterFactoryV1 = await ethers.getContractFactory("Counter");
  const counterV1 = await upgrades.deployProxy(CounterFactoryV1, [], {
    deployer
  });
  await counterV1.waitForDeployment();
  const proxyAddress = await counterV1.getAddress();
  const implAddress = await upgrades.erc1967.getImplementationAddress(
    proxyAddress
  );
  console.log("CounterFactoryV1 deployed to:", proxyAddress);
  console.log("CounterFactoryV1 implAddress to:", implAddress);

  const storePath = path.resolve(__dirname, "../cache/store/proxyFactory.json");
  fs.writeFileSync(
    storePath,
    JSON.stringify({
      proxyAddress,
      implAddress,
      abi: CounterFactoryV1.interface.format("json"),
    })
  );

  await save("CounterProxy", {
    abi: CounterFactoryV1.interface.format("json"),
    address: proxyAddress,
    implAddress,
    args: [],
    log: true,
  });
};

module.exports.tags = ["CounterFactoryV1"];
