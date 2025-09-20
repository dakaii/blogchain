# BlogChain - Decentralized Blogging Platform

**BlogChain** is a blockchain-powered blogging platform built using Cosmos SDK and Vue 3. It enables users to create, share, and interact with blog posts on a decentralized network.

## 🚀 Features

### Blockchain Features
- **Create Blog Posts**: Publish content directly on the blockchain
- **Like System**: Interact with posts through on-chain likes
- **Tag Support**: Organize posts with tags
- **Immutable Content**: All posts are permanently stored on-chain

### Wallet & Financial Features
- **Keplr Wallet Integration**: Connect with Keplr browser extension
- **Mnemonic Support**: Import wallets using seed phrases
- **Token Transfers**: Send and receive STAKE tokens
- **Transaction History**: View all your blockchain transactions
- **Asset Management**: Track your token balances

### Frontend Features
- **Modern Vue 3**: Built with Composition API and TypeScript
- **Tailwind CSS**: Beautiful, responsive UI design
- **Real-time Updates**: Live blockchain data synchronization
- **Search & Filter**: Find posts by tags or content
- **Wallet Dashboard**: Complete wallet management interface

## 📋 Prerequisites

- Node.js v20.19.0 or higher
- Go 1.24.0 or higher
- [Ignite CLI](https://docs.ignite.com/welcome/install)
- Git

## 🛠️ Installation

1. **Clone the repository**
```bash
git clone https://github.com/dakaii/blogchain.git
cd blogchain
```

2. **Install Ignite CLI** (if not already installed)
```bash
curl https://get.ignite.com/cli! | bash
```

## 🎮 Quick Start

### Start the Blockchain

```bash
# Start the blockchain with auto-reset
ignite chain serve --reset-once

# Or start without reset to preserve data
ignite chain serve
```

The blockchain will be available at:
- **Tendermint RPC**: http://localhost:26657
- **REST API**: http://localhost:1317
- **gRPC**: localhost:9090

### Start the Frontend

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The Vue application will be available at: http://localhost:5173

## 🧪 Test Accounts

Use these test accounts for development:

### Alice (Main test account)
```
Mnemonic: banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass
Address: blogchain1...
Balance: 200,000,000 STAKE, 20,000 TOKEN
```

### Bob
```
Balance: 100,000,000 STAKE, 10,000 TOKEN
```

## 📱 Using the Application

### Connect Wallet
1. Open http://localhost:5173
2. Click "Connect Wallet"
3. Choose either:
   - **Keplr**: Use Keplr browser extension
   - **Mnemonic**: Paste the test mnemonic above

### Create a Blog Post
1. Navigate to "Create Post" in the menu
2. Fill in the title, content, and tags
3. Click "Publish Post"
4. Confirm the transaction in your wallet

### Explore Posts
1. Go to "Explore" to browse all posts
2. Use search and tag filters to find content
3. Click the heart icon to like posts
4. Click "Read more" to view full posts

### Manage Wallet
1. Go to "Wallet" in the menu
2. View your asset balances
3. Send tokens to other addresses
4. Review transaction history

## 🏗️ Project Structure

```
blogchain/
├── app/               # Blockchain application setup
├── x/
│   └── blog/         # Blog module
│       ├── keeper/   # Business logic
│       ├── types/    # Types and messages
│       └── proto/    # Protocol buffer definitions
├── proto/            # Proto files
│   └── blogchain/
│       └── blog/
│           └── v1/   # API definitions
├── frontend/         # Vue 3 application
│   ├── src/
│   │   ├── components/   # Vue components
│   │   ├── views/        # Page views
│   │   ├── services/     # Blockchain services
│   │   └── stores/       # Pinia stores
│   └── package.json
└── config.yml        # Ignite configuration
```

## 🔧 Development

### Build the Blockchain
```bash
ignite chain build
```

### Generate Proto Files
```bash
ignite generate proto-go --yes
```

### Run Tests
```bash
go test ./...
```

### Frontend Development
```bash
cd frontend
npm run dev      # Development server
npm run build    # Production build
npm run preview  # Preview production build
```

## 📝 Available Commands

### Blockchain Commands
```bash
# Query posts
blogchaind query blog posts

# Create a post (CLI)
blogchaind tx blog create-post "Title" "Content" "tag1,tag2" --from alice

# Like a post
blogchaind tx blog like-post [post-id] --from alice

# Check balance
blogchaind query bank balances [address]
```

### Frontend Scripts
```bash
npm run dev        # Start dev server
npm run build      # Build for production
npm run preview    # Preview production build
npm run type-check # Run TypeScript checks
```

## 🌐 API Endpoints

### REST API
- `GET /blogchain/blog/v1/posts` - List all posts
- `GET /blogchain/blog/v1/posts/{id}` - Get specific post
- `POST /blogchain/blog/v1/tx` - Submit transactions

### WebSocket
- `ws://localhost:26657/websocket` - Real-time updates

## 🚢 Deployment

### Build for Production
```bash
# Build blockchain binary
ignite chain build --release

# Build frontend
cd frontend
npm run build
```

### Docker (Coming Soon)
```bash
docker-compose up
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 🔗 Resources

- [Cosmos SDK Documentation](https://docs.cosmos.network)
- [Ignite CLI Documentation](https://docs.ignite.com)
- [Vue 3 Documentation](https://vuejs.org)
- [Tailwind CSS Documentation](https://tailwindcss.com)

## 🐛 Troubleshooting

### Blockchain won't start
```bash
# Reset the blockchain data
ignite chain serve --reset-once
```

### Frontend connection issues
- Ensure blockchain is running on port 26657
- Check API proxy settings in `vite.config.ts`
- Verify CORS settings

### Keplr connection fails
- Install Keplr extension from Chrome Web Store
- Ensure blockchain is running
- Try using mnemonic connection instead

## 📮 Support

For issues and questions:
- Open an issue on [GitHub](https://github.com/dakaii/blogchain/issues)
- Join our Discord community (coming soon)

---

Built with ❤️ using Cosmos SDK and Vue 3