import { Contract, getAddress } from "ethers";
import { didRegistryAbi } from "../contracts/didRegistryAbi";
import { connectWallet } from "./wallet";

export interface PrepareDIDResult {
  contract_address: string;
  did: string;
  github: string;
  metadata_cid: string;
  chain_id: number;
}

export async function registerDIDOnChain(params: PrepareDIDResult) {
  const { signer } = await connectWallet();
  const contract = new Contract(getAddress(params.contract_address), didRegistryAbi, signer);
  const tx = await contract.registerDID(params.did, params.github || "", params.metadata_cid);
  const receipt = await tx.wait();
  return receipt?.hash as string;
}
