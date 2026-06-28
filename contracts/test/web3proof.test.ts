import { expect } from "chai";
import { ethers } from "hardhat";

describe("Web3Proof evidence platform contracts", function () {
  it("registers a DID profile", async function () {
    const [user] = await ethers.getSigners();
    const DIDProfile = await ethers.getContractFactory("DIDProfile");
    const did = await DIDProfile.deploy();
    await did.waitForDeployment();

    await expect(did.connect(user).registerProfile("did:web3proof:demo", "ipfs://profile"))
      .to.emit(did, "ProfileRegistered");

    const profile = await did.getProfile(user.address);
    expect(profile.did).to.equal("did:web3proof:demo");
  });

  it("creates and verifies evidence by file hash", async function () {
    const [admin, user] = await ethers.getSigners();
    const EvidenceRegistry = await ethers.getContractFactory("EvidenceRegistry");
    const registry = await EvidenceRegistry.deploy(admin.address);
    await registry.waitForDeployment();

    const evidenceNoHash = ethers.sha256(ethers.toUtf8Bytes("EV-1"));
    const fileHash = ethers.sha256(ethers.toUtf8Bytes("file"));
    await expect(registry.connect(user).createEvidence(evidenceNoHash, fileHash, "ipfs://meta"))
      .to.emit(registry, "EvidenceCreated");

    const result = await registry.verifyEvidence(fileHash);
    expect(result[0]).to.equal(true);
    expect(result[2]).to.equal(user.address);
  });

  it("prevents duplicate evidence file hashes", async function () {
    const [admin, user] = await ethers.getSigners();
    const EvidenceRegistry = await ethers.getContractFactory("EvidenceRegistry");
    const registry = await EvidenceRegistry.deploy(admin.address);
    await registry.waitForDeployment();

    const evidenceNoHash = ethers.sha256(ethers.toUtf8Bytes("EV-1"));
    const fileHash = ethers.sha256(ethers.toUtf8Bytes("file"));
    await registry.connect(user).createEvidence(evidenceNoHash, fileHash, "ipfs://meta");
    await expect(registry.connect(user).createEvidence(evidenceNoHash, fileHash, "ipfs://meta"))
      .to.be.revertedWith("File hash already exists");
  });

  it("mints non-transferable credential SBTs", async function () {
    const [admin, user] = await ethers.getSigners();
    const CredentialSBT = await ethers.getContractFactory("CredentialSBT");
    const sbt = await CredentialSBT.deploy(admin.address);
    await sbt.waitForDeployment();

    await expect(sbt.connect(admin).mintCredential(user.address, 1, "ipfs://token"))
      .to.emit(sbt, "CredentialMinted");

    expect(await sbt.ownerOf(1)).to.equal(user.address);
    await expect(sbt.connect(user).transferFrom(user.address, admin.address, 1))
      .to.be.revertedWith("SBT: non-transferable");
  });

  it("updates reputation score", async function () {
    const [admin, user] = await ethers.getSigners();
    const Reputation = await ethers.getContractFactory("Reputation");
    const rep = await Reputation.deploy(admin.address);
    await rep.waitForDeployment();

    await rep.connect(admin).updateScore(user.address, 500, 300, 200);
    const score = await rep.getScore(user.address);
    expect(score.total).to.equal(1000n);
  });
});
