import { Contract, Interface, getAddress } from "ethers";
import { evidenceRegistryAbi } from "../contracts/evidenceRegistryAbi";
import { EvidencePrepareResult } from "../api/works";
import { connectWallet } from "./wallet";

export async function submitEvidenceOnChain(params: EvidencePrepareResult) {
  if (!params.contract_address) {
    throw new Error("未配置存证合约地址");
  }
  const { signer } = await connectWallet();
  const contract = new Contract(getAddress(params.contract_address), evidenceRegistryAbi, signer);
  const tx = await contract.createEvidence(params.evidence_no_hash, params.file_hash, params.metadata_uri);
  const receipt = await tx.wait();
  if (!receipt) throw new Error("交易未确认");

  const iface = new Interface(evidenceRegistryAbi);
  for (const log of receipt.logs) {
    try {
      const parsed = iface.parseLog({ topics: [...log.topics], data: log.data });
      if (parsed?.name === "EvidenceCreated") {
        return {
          txHash: receipt.hash as string,
          chainEvidenceId: Number(parsed.args.evidenceId),
        };
      }
    } catch {
      // Ignore unrelated logs.
    }
  }
  throw new Error("未找到 EvidenceCreated 事件");
}
