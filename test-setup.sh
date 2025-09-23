#!/bin/bash

echo "🔍 Testing BlogChain Setup..."
echo "================================"

# Check if blockchain is running
echo -n "✓ Blockchain Status: "
if curl -s http://localhost:26657/status > /dev/null 2>&1; then
    echo "Running ✅"
    BLOCK_HEIGHT=$(curl -s http://localhost:26657/status | grep -o '"latest_block_height":"[0-9]*"' | grep -o '[0-9]*' || echo "0")
    echo "  Current block height: $BLOCK_HEIGHT"
else
    echo "Not running ❌"
    echo "  Run: ignite chain serve --reset-once"
fi

# Check if API is accessible
echo -n "✓ API Status: "
if curl -s http://localhost:1317/cosmos/base/tendermint/v1beta1/syncing > /dev/null 2>&1; then
    echo "Running ✅"
else
    echo "Not running ❌"
fi

# Check if frontend is running
echo -n "✓ Frontend Status: "
if curl -s http://localhost:5174 > /dev/null 2>&1; then
    echo "Running on http://localhost:5174 ✅"
elif curl -s http://localhost:5173 > /dev/null 2>&1; then
    echo "Running on http://localhost:5173 ✅"
else
    echo "Not running ❌"
    echo "  Run: npm run dev"
fi

echo ""
echo "📝 Test Account:"
echo "================================"
echo "Mnemonic:"
echo "banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass"
echo ""
echo "Address: $(echo 'banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass' | blogchaind keys add test --recover --keyring-backend test --output json 2>/dev/null | grep -o '"address":"[^"]*"' | cut -d'"' -f4 || echo 'Run blockchain first')"

echo ""
echo "🚀 Quick Start:"
echo "================================"
echo "1. Open browser to http://localhost:5174"
echo "2. Click 'Import with Seed Phrase'"
echo "3. Paste the test mnemonic above"
echo "4. Start creating posts!"