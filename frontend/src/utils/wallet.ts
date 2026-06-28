import { BrowserProvider } from "ethers";
import { CHAIN_ID } from "../config";

const SEPOLIA = {
  chainId: "0xaa36a7",
  chainName: "Sepolia",
  nativeCurrency: { name: "Sepolia ETH", symbol: "ETH", decimals: 18 },
  rpcUrls: ["https://rpc.sepolia.org"],
  blockExplorerUrls: ["https://sepolia.etherscan.io"],
};

async function ensureSepoliaNetwork() {
  if (!window.ethereum) return;

  const currentChainId = await window.ethereum.request({ method: "eth_chainId" });
  if (String(currentChainId).toLowerCase() === SEPOLIA.chainId) {
    return;
  }

  try {
    await window.ethereum.request({
      method: "wallet_switchEthereumChain",
      params: [{ chainId: SEPOLIA.chainId }],
    });
  } catch (err: unknown) {
    const code = (err as { code?: number })?.code;
    if (code === 4902) {
      await window.ethereum.request({
        method: "wallet_addEthereumChain",
        params: [SEPOLIA],
      });
      return;
    }
    throw err;
  }
}

export async function connectWallet() {
  if (!window.ethereum) {
    throw new Error("请先安装 MetaMask 钱包");
  }
  await ensureSepoliaNetwork();
  const provider = new BrowserProvider(window.ethereum);
  await provider.send("eth_requestAccounts", []);
  const network = await provider.getNetwork();
  if (Number(network.chainId) !== CHAIN_ID) {
    throw new Error(`请切换到 Sepolia 测试网 (chainId: ${CHAIN_ID})`);
  }
  const signer = await provider.getSigner();
  const address = await signer.getAddress();
  return { provider, signer, address };
}

declare global {
  interface Window {
    ethereum?: {
      request: (args: { method: string; params?: unknown[] }) => Promise<unknown>;
      on?: (event: "accountsChanged" | "chainChanged", handler: () => void) => void;
      removeListener?: (event: "accountsChanged" | "chainChanged", handler: () => void) => void;
      isMetaMask?: boolean;
    };
  }
}
