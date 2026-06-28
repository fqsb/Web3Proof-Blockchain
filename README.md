# Web3Proof

Web3Proof 是一个基于区块链的数字作品与成果存证、认证、核验平台。系统面向创作者、核验方、认证审核员和管理员，支持文件上传、SHA-256 哈希计算、Sepolia 链上存证、PDF 存证证书、SBT 认证凭证、公开可信档案和第三方核验报告。

本项目不发行代币，不提供 NFT 交易、充值提现或二级市场功能。链上只保存文件哈希、作者地址、存证编号哈希和元数据 URI，原始文件和业务数据保存在链下服务器。

## 项目结构

```text
proof/
├── frontend/   React + Vite + TypeScript + Ant Design + Ethers.js
├── backend/    Go + Gin + Gorm + JWT
├── contracts/  Hardhat + Solidity + OpenZeppelin
├── deploy/     Docker Compose + Nginx + MySQL + Redis + MinIO
└── docs/       系统文档与数据库脚本
```

## 核心功能

- 钱包签名登录，JWT 中携带当前角色。
- 多角色切换：`creator`、`verifier`、`auditor`、`admin`。
- 作品/成果上传，链下保存文件，链上保存哈希摘要。
- `EvidenceRegistry` 存证合约记录证据编号、文件哈希、作者地址和时间。
- 认证审核通过后由 `CredentialSBT` 发放不可转让凭证。
- 支持按文件、存证编号、证书编号、钱包地址进行核验。
- 公开可信档案聚合用户作品、存证记录、认证凭证和可信评分。

## 本地启动

### 1. 合约

```bash
cd contracts
npm install
npm test
npm run deploy:sepolia
```

部署后会生成 `contracts/deployments/<chainId>.json`，Sepolia 部署会写入 `frontend/.env.development.local` 和 `deploy/.env`。

### 2. 后端

```bash
cd backend
go mod tidy
set DB_DRIVER=sqlite
set DB_FILE=./web3proof.db
set CHAIN_ID=11155111
go run ./cmd/server
```

生产 MySQL 初始化脚本在 `docs/database/schema.sql`。

### 3. 前端

```bash
cd frontend
npm install
npm run dev
```

访问 `http://localhost:5173`。

## Docker 部署

```bash
copy deploy\.env.example deploy\.env
cd deploy
docker compose up -d --build
```

部署前需要在 `deploy/.env` 中填写：

- `JWT_SECRET`
- `SEPOLIA_RPC_URL`
- `ADMIN_WALLET_ADDRESS`
- `DID_PROFILE_ADDRESS`
- `EVIDENCE_REGISTRY_ADDRESS`
- `CREDENTIAL_SBT_ADDRESS`
- `REPUTATION_ADDRESS`

## 验证

```bash
cd contracts && npm test
cd backend && go test ./...
cd frontend && npm run build
```
