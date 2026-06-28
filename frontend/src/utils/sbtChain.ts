import { Contract, Interface, getAddress } from "ethers";
import { credentialSbtAbi } from "../contracts/credentialSbtAbi";
import { connectWallet } from "./wallet";

export interface CredentialMintPrepareResult {
  contract_address: string;
  to_address: string;
  evidence_id: number;
  token_uri: string;
  chain_id: number;
}

export async function mintCredentialOnChain(params: CredentialMintPrepareResult) {
  const { signer } = await connectWallet();
  const contract = new Contract(getAddress(params.contract_address), credentialSbtAbi, signer);
  const tx = await contract.mintCredential(getAddress(params.to_address), params.evidence_id, params.token_uri);
  const receipt = await tx.wait();
  if (!receipt) throw new Error("凭证交易未确认");

  const iface = new Interface(credentialSbtAbi);
  for (const log of receipt.logs) {
    try {
      const parsed = iface.parseLog({ topics: [...log.topics], data: log.data });
      if (parsed?.name === "CredentialMinted") {
        return { txHash: receipt.hash as string, tokenId: Number(parsed.args.tokenId) };
      }
    } catch {
      // Ignore unrelated logs.
    }
  }
  throw new Error("未找到 CredentialMinted 事件");
}
