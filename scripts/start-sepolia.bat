@echo off
setlocal enabledelayedexpansion
cd /d %~dp0\..

echo === Web3Proof Sepolia local start ===

for /f "usebackq tokens=1,* delims==" %%A in ("deploy\.env") do (
  set "line=%%A"
  if not "!line:~0,1!"=="#" if not "%%A"=="" set "%%A=%%B"
)

echo Stopping old backend process on port 8080...
for /f "tokens=5" %%P in ('netstat -aon ^| findstr :8080 ^| findstr LISTENING') do (
  taskkill /F /PID %%P >nul 2>&1
)

echo Starting backend...
start "web3proof-backend" cmd /k "cd backend && set DB_DRIVER=sqlite&& set DB_FILE=./web3proof.db&& set CHAIN_ID=11155111&& set SEPOLIA_RPC_URL=%SEPOLIA_RPC_URL%&& set JWT_SECRET=dev-jwt-secret-web3proof-local-32chars&& set CORS_ORIGIN=http://localhost:5173&& set ADMIN_WALLET_ADDRESS=%ADMIN_WALLET_ADDRESS%&& set DID_PROFILE_ADDRESS=%DID_PROFILE_ADDRESS%&& set EVIDENCE_REGISTRY_ADDRESS=%EVIDENCE_REGISTRY_ADDRESS%&& set CREDENTIAL_SBT_ADDRESS=%CREDENTIAL_SBT_ADDRESS%&& set REPUTATION_ADDRESS=%REPUTATION_ADDRESS%&& go run ./cmd/server"

timeout /t 2 /nobreak >nul

echo Starting frontend...
start "web3proof-frontend" cmd /c "cd frontend && npm run dev"

echo.
echo Done.
echo   Frontend: http://localhost:5173
echo   Backend: http://localhost:8080/health
echo   MetaMask: Sepolia test network
echo   Admin wallet: %ADMIN_WALLET_ADDRESS%
echo.
echo Demo order:
echo   1. Connect wallet and sign in
echo   2. Create a work and upload a file
echo   3. Submit evidence to Sepolia
echo   4. Submit a certification application
echo   5. Review the application and mint SBT
echo   6. Verify by file, evidence number, certificate number, or wallet
pause
