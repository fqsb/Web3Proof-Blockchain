@echo off
setlocal enabledelayedexpansion

echo === Web3Proof 本地启动 ===

cd /d %~dp0\..

echo [1/4] 启动 Hardhat 本地链...
start "hardhat-node" cmd /c "cd contracts && npx hardhat node"
timeout /t 5 /nobreak >nul

echo [2/4] 部署合约到 localhost...
cd contracts
call npx hardhat run scripts/deploy.ts --network localhost
if errorlevel 1 exit /b 1

for /f "delims=" %%i in ('node -pe "const j=require('./deployments/31337.json'); j.ProjectRegistry"') do set PROJECT_REGISTRY=%%i
for /f "delims=" %%i in ('node -pe "const j=require('./deployments/31337.json'); j.DIDRegistry"') do set DID_REGISTRY=%%i
cd ..

echo ProjectRegistry: %PROJECT_REGISTRY%

echo [3/4] 写入前端环境变量...
(
echo VITE_API_BASE_URL=/api/v1
echo VITE_CHAIN_ID=31337
echo VITE_DID_REGISTRY_ADDRESS=%DID_REGISTRY%
echo VITE_PROJECT_REGISTRY_ADDRESS=%PROJECT_REGISTRY%
) > frontend\.env.development.local

echo [4/4] 启动后端与前端...
set DB_DRIVER=sqlite
set DB_FILE=./web3proof.db
set PROJECT_REGISTRY_ADDRESS=%PROJECT_REGISTRY%
set DID_REGISTRY_ADDRESS=%DID_REGISTRY%
set CHAIN_ID=31337
set CORS_ORIGIN=http://localhost:5173
set JWT_SECRET=dev-jwt-secret-web3proof-local-32chars

start "backend" cmd /c "cd backend && go run ./cmd/server"
timeout /t 3 /nobreak >nul
start "frontend" cmd /c "cd frontend && npm run dev"

echo.
echo 完成！请访问 http://localhost:5173
echo MetaMask 请添加本地网络: RPC http://127.0.0.1:8545, Chain ID 31337
echo 使用 Hardhat 默认账户导入私钥进行测试
pause
