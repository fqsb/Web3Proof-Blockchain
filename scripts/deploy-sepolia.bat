@echo off
setlocal
cd /d %~dp0\..

echo === Deploy contracts to Sepolia ===
echo Make sure deploy\.env has BACKEND_WALLET_PRIVATE_KEY configured.
echo The deployer wallet must have Sepolia ETH for gas.
echo.

if not exist deploy\.env (
  echo Missing deploy\.env. Please copy deploy\.env.example first.
  exit /b 1
)

cd contracts
call npx hardhat run scripts/deploy.ts --network sepolia
if errorlevel 1 (
  echo Deploy failed. Please check private key, RPC, and Sepolia ETH balance.
  exit /b 1
)

echo.
echo Deploy finished. Contract addresses were written to:
echo   - contracts\deployments\11155111.json
echo   - frontend\.env.development.local
echo   - deploy\.env
echo.
echo Restart backend and frontend so they load the new env values.
echo Switch MetaMask to Sepolia.
pause
