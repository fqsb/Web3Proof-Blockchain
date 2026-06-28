import { ethers } from "hardhat";
import * as fs from "fs";
import * as path from "path";

const MINTER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("MINTER_ROLE"));
const UPDATER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("UPDATER_ROLE"));

function writeFrontendEnv(filePath: string, addresses: Record<string, string>, chainId: string) {
  const content = [
    "VITE_API_BASE_URL=/api/v1",
    `VITE_CHAIN_ID=${chainId}`,
    `VITE_DID_PROFILE_ADDRESS=${addresses.DIDProfile}`,
    `VITE_EVIDENCE_REGISTRY_ADDRESS=${addresses.EvidenceRegistry}`,
    `VITE_CREDENTIAL_SBT_ADDRESS=${addresses.CredentialSBT}`,
    `VITE_REPUTATION_ADDRESS=${addresses.Reputation}`,
    "",
  ].join("\n");
  fs.writeFileSync(filePath, content);
  console.log("Frontend env written:", filePath);
}

function patchDeployEnv(filePath: string, addresses: Record<string, string>, chainId: string) {
  if (!fs.existsSync(filePath)) return;
  let content = fs.readFileSync(filePath, "utf8");
  const replacements: Record<string, string> = {
    DID_PROFILE_ADDRESS: addresses.DIDProfile,
    EVIDENCE_REGISTRY_ADDRESS: addresses.EvidenceRegistry,
    CREDENTIAL_SBT_ADDRESS: addresses.CredentialSBT,
    REPUTATION_ADDRESS: addresses.Reputation,
    CHAIN_ID: chainId,
    VITE_CHAIN_ID: chainId,
  };
  for (const [key, value] of Object.entries(replacements)) {
    const re = new RegExp(`^${key}=.*$`, "m");
    content = re.test(content) ? content.replace(re, `${key}=${value}`) : `${content}\n${key}=${value}`;
  }
  fs.writeFileSync(filePath, content);
  console.log("Deploy env patched:", filePath);
}

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("Deploying with account:", deployer.address);

  const DIDProfile = await ethers.getContractFactory("DIDProfile");
  const didProfile = await DIDProfile.deploy();
  await didProfile.waitForDeployment();

  const EvidenceRegistry = await ethers.getContractFactory("EvidenceRegistry");
  const evidenceRegistry = await EvidenceRegistry.deploy(deployer.address);
  await evidenceRegistry.waitForDeployment();

  const CredentialSBT = await ethers.getContractFactory("CredentialSBT");
  const credentialSBT = await CredentialSBT.deploy(deployer.address);
  await credentialSBT.waitForDeployment();

  const Reputation = await ethers.getContractFactory("Reputation");
  const reputation = await Reputation.deploy(deployer.address);
  await reputation.waitForDeployment();

  const backendWallet = process.env.BACKEND_WALLET_ADDRESS;
  if (backendWallet && ethers.isAddress(backendWallet)) {
    await (await credentialSBT.grantRole(MINTER_ROLE, backendWallet)).wait();
    await (await reputation.grantRole(UPDATER_ROLE, backendWallet)).wait();
  }

  const chainId = (await ethers.provider.getNetwork()).chainId.toString();
  const addresses = {
    DIDProfile: await didProfile.getAddress(),
    EvidenceRegistry: await evidenceRegistry.getAddress(),
    CredentialSBT: await credentialSBT.getAddress(),
    Reputation: await reputation.getAddress(),
    deployer: deployer.address,
    chainId,
    deployedAt: new Date().toISOString(),
  };

  const root = path.join(__dirname, "..");
  const outDir = path.join(root, "deployments");
  fs.mkdirSync(outDir, { recursive: true });
  fs.writeFileSync(path.join(outDir, `${chainId}.json`), JSON.stringify(addresses, null, 2));

  if (chainId === "11155111") {
    writeFrontendEnv(path.join(root, "..", "frontend", ".env.development.local"), addresses, chainId);
    patchDeployEnv(path.join(root, "..", "deploy", ".env"), addresses, chainId);
  }

  console.log(addresses);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
