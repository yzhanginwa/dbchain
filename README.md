# DBChain

A blockchain-based database management system built on Cosmos SDK and Tendermint consensus engine.

## Overview

**DBChain** transforms traditional database operations into blockchain transactions, creating an immutable, distributed database where all schema changes and data operations are recorded on-chain. It provides a decentralized database solution with cryptographic verification, consensus-based validation, and transparent audit trails.

### Key Features

- **Blockchain-Native Database**: Database operations (creating tables, inserting data, managing schemas) are recorded as blockchain transactions
- **Immutability & Transparency**: All changes are cryptographically verified and permanently recorded
- **BFT Consensus**: Uses Tendermint for Byzantine Fault Tolerant consensus
- **Lua Scripting**: Programmable database with filters, triggers, and functions
- **Dual Cryptography Support**: International edition (secp256k1) and SM2 edition (Chinese national cryptography standard)
- **BSN Integration**: Built for China's Blockchain Service Network with specialized features

## Architecture

### Core Components

DBChain is built on a sophisticated blockchain architecture:

```
┌─────────────────────────────────────────────────────┐
│              DBChain Application Layer              │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────┐  │
│  │   Database   │  │     User     │  │    BSN    │  │
│  │  Operations  │  │  Management  │  │Integration│  │
│  └──────────────┘  └──────────────┘  └───────────┘  │
├─────────────────────────────────────────────────────┤
│          Cosmos SDK v0.39.2 (Forked)                │
│  ┌──────┐ ┌──────┐ ┌────────┐ ┌──────────────────┐  │
│  │ Auth │ │ Bank │ │Staking │ │ DBChain Module   │  │
│  └──────┘ └──────┘ └────────┘ └──────────────────┘  │
├─────────────────────────────────────────────────────┤
│      Tendermint v0.33.8 (BFT Consensus)             │
└─────────────────────────────────────────────────────┘
```

### Modules

**DBChain Module** (`/x/dbchain/`)
- Core database functionality
- Schema management (applications, tables, columns, indexes)
- Data operations (insert, query, freeze)
- Lua scripting engine with custom ORM
- User/group management and permissions
- Friend system and social features

**Bank Module** (`/x/bank/`)
- Token transfers and balance management
- Gas fee handling using `dbctoken`
- P2P transfer limits for BSN compliance

**Standard Cosmos Modules**
- Authentication and account management
- Validator staking and delegation
- Fee distribution
- Slashing for misbehavior
- Token supply management

### Two Editions

1. **International Version**: Standard secp256k1 cryptography with "cosmos" address prefix
2. **SM2 Version** (国密版本): Chinese national cryptography standard (SM2/SM3) with "dbchain" address prefix

Both versions are compiled from the same codebase using build tags for conditional compilation.

## Features

### 1. Database Operations

**Application Management**
- Create, drop, and recover applications (databases)
- Set permissions and freeze status
- Admin-controlled or user-controlled modes
- Community edition limited to 2 applications

**Schema Management**
- Create tables with custom fields and data types
- Add, drop, and rename columns
- Create and drop indexes for performance
- Table associations and relationships
- Counter cache fields for optimization

**Data Operations**
- Insert rows with validation
- Query system with custom queriers
- Fuzzy search capabilities
- Row freezing for immutability
- Transaction-based data integrity

### 2. Lua Scripting System

DBChain includes a powerful Lua scripting engine for programmable database logic:

- **Functions**: Reusable Lua functions stored on-chain
- **Insert Filters**: Pre-insertion validation and data transformation
- **Triggers**: Post-insertion actions and side effects
- **ORM Interface**: Database access from Lua scripts
- **Safety**: Loop count limiting to prevent infinite loops
- **Custom BNF Grammar**: Parser for script expressions

### 3. User & Group Management

- User authentication and authorization
- Group creation and membership management
- Database user roles and permissions
- Friend system (add, respond, drop)
- Admin accounts with elevated privileges
- Genesis-level account initialization

### 4. BSN Integration

Specialized features for China's Blockchain Service Network:

- **Key Custody Service** (托管密钥): Managed private key storage
- **Account Creation**: With or without key custody
- **Token Recharge API**: Payment integration
- **P2P Transfer Limits**: Compliance with BSN regulations
- **Token Keeper Roles**: Admin accounts for managed transfers
- **Modified Gas Model**: Read operations cost 0 gas for commercial viability

### 5. Oracle Service

External integration layer for enterprise applications:

- REST API for blockchain interactions
- Key generation and custody
- Payment integration (Apple Pay, Alipay)
- IPFS file upload support
- Block browser functionality
- Certificate management
- Authenticator support (Google Authenticator, HMAC)

### 6. Gas & Fee System

- Configurable gas consumption per operation type
- Fee refund mechanism for unused gas
- Transaction cost tracking per account
- Minimum gas price enforcement
- Special gas rules for BSN compliance

### 7. Security Features

- BFT consensus (Byzantine Fault Tolerant)
- Transaction signing and cryptographic verification
- Access control via groups and permissions
- Script safety (loop count limits)
- Schema and data freezing
- Row-level access control

## Quick Start

### Prerequisites

- **Go**: Version 1.13 or higher
- **Make**: For building the project
- **Git**: For cloning the repository

### Building from Source

#### International Version (Standard Cryptography)

```bash
make install
```

This builds binaries with secp256k1 cryptography and "cosmos" address prefix:
- `dbchaind` - Blockchain node daemon
- `dbchaincli` - Command-line client
- `dbchainoracle` - Oracle service

#### SM2 Version (Chinese National Cryptography)

```bash
make install pkc=sm2
```

This builds binaries with SM2/SM3 cryptography and "dbchain" address prefix:
- `dbchaind_sm2` - Blockchain node daemon
- `dbchaincli_sm2` - Command-line client
- `dbchainoracle_sm2` - Oracle service

#### Community Edition

```bash
make installc
```

Limited edition with restrictions on the number of applications.

### Configuration

After building, initialize the node:

```bash
# Initialize node configuration
dbchaind init <moniker> --chain-id testnet

# Configure persistent peers (for multi-node deployment)
# Edit ~/.dbchaind/config/config.toml

# Set minimum gas prices
# Edit ~/.dbchaind/config/app.toml
minimum-gas-prices = "0.000001dbctoken"
```

### Running a Node

```bash
# Start the blockchain node
dbchaind start

# In another terminal, start the REST API server
dbchaincli rest-server --chain-id testnet --trust-node

# Start the oracle service (optional)
dbchainoracle serve
```

### Basic CLI Usage

```bash
# Create a new account
dbchaincli keys add myaccount

# Check account balance
dbchaincli query account $(dbchaincli keys show myaccount -a)

# Create a database application
dbchaincli tx dbchain create-application myapp --from myaccount

# Create a table
dbchaincli tx dbchain create-table myapp mytable \
  --fields "name:string,age:int,email:string" \
  --from myaccount

# Insert data
dbchaincli tx dbchain insert-row myapp mytable \
  --data '{"name":"Alice","age":30,"email":"alice@example.com"}' \
  --from myaccount

# Query data
dbchaincli query dbchain get-table myapp mytable
```

### Multi-Node Deployment

For production deployment with BFT consensus (requires 4+ nodes):

1. Generate genesis transactions on each validator node
2. Collect all genesis transactions
3. Create final genesis.json
4. Configure persistent peers
5. Start all nodes simultaneously

Refer to `/doc/dbchain 部署及初始化文档.md` for detailed deployment instructions (Chinese).

## Directory Structure

```
dbchain/
├── app.go                      # Main application setup
├── Makefile                    # Build system
├── go.mod, go.sum             # Go dependencies
├── README.md                   # This file
├── readme.txt                  # Build instructions (Chinese)
│
├── address/                    # Address prefix configuration
│   ├── bech32_main_prefix_default.go
│   └── bech32_main_prefix_sm2.go
│
├── cmd/                        # Command-line executables
│   ├── dbchaind/              # Node daemon
│   ├── dbchaincli/            # CLI client
│   └── dbchainoracle/         # Oracle service
│
├── doc/                        # Documentation (Chinese)
│   ├── DBChain技术架构.docx
│   ├── DBChain白皮书1.0.docx
│   ├── dbchain 部署及初始化文档.md
│   ├── gas 消耗说明文档.md
│   ├── bsn开放接口对应库链接口.md
│   └── dbchain nginx 配置.md
│
└── x/                          # Custom modules
    ├── dbchain/               # Main database module
    │   ├── handler.go         # Transaction handlers
    │   ├── client/            # Client interfaces
    │   │   ├── cli/          # CLI commands
    │   │   ├── rest/         # REST API
    │   │   └── oracle/       # Oracle service
    │   └── internal/
    │       ├── keeper/        # Business logic
    │       ├── types/         # Message types
    │       ├── super_script/  # Script parser
    │       └── utils/         # Utilities
    └── bank/                  # Token/balance module
```

## Technology Stack

- **Blockchain Framework**: Cosmos SDK v0.39.2 (forked with SM2 support)
- **Consensus Engine**: Tendermint v0.33.8 (forked with SM2 support)
- **Programming Language**: Go 1.13+
- **Scripting Engine**: Lua (via gopher-lua)
- **Database**: Tendermint KV store
- **Serialization**: Amino (Cosmos)
- **HTTP Router**: Gorilla Mux
- **CLI Framework**: Cobra + Viper

## Key Dependencies

- `github.com/dbchaincloud/cosmos-sdk` - Forked Cosmos SDK with SM2 cryptography
- `github.com/dbchaincloud/tendermint` - Forked Tendermint with SM2 cryptography
- `github.com/yuin/gopher-lua` - Lua VM for Go
- `github.com/ipfs/go-ipfs-api` - IPFS integration
- `github.com/smartwalle/alipay/v3` - Alipay payment integration
- `github.com/aliyun/alibaba-cloud-sdk-go` - Alibaba Cloud services

## Configuration Files

- `~/.dbchaind/config/config.toml` - Node configuration (peers, RPC, consensus)
- `~/.dbchaind/config/app.toml` - Application configuration (gas prices)
- `~/.dbchaind/config/genesis.json` - Genesis state and initial validators
- `~/.dbchaincli/config/config.toml` - CLI client configuration
- `~/.dbchainoracle/config/config.toml` - Oracle service configuration

## API Access Methods

DBChain provides three ways to interact with the blockchain:

1. **CLI** (`dbchaincli`): Command-line interface for transactions and queries
2. **REST API** (`dbchaincli rest-server`): HTTP API for web applications
3. **Oracle HTTP Server** (`dbchainoracle serve`): Enterprise integration layer

## Tokens

- **dbctoken**: Gas token for transaction fees
- **stake**: Staking token for validators

## Documentation

Comprehensive Chinese documentation is available in the `/doc` directory:

- **DBChain技术架构.docx** - Technical architecture
- **DBChain白皮书1.0.docx** - Project whitepaper
- **dbchain 部署及初始化文档.md** - Deployment and initialization guide
- **gas 消耗说明文档.md** - Gas consumption specification
- **bsn开放接口对应库链接口.md** - BSN API interface mapping
- **dbchain nginx 配置.md** - Nginx configuration guide

## Storage Structure

DBChain uses a key-value store with the following structure:

```
appcode:<appCode>                           → database struct
db:<appId>:mt:tables                        → table list
db:<appId>:mt:tn:<tableName>               → table fields
db:<appId>:mt:idx:<tableName>              → table indexes
db:<appId>:dt:<tableName>:<fieldName>:<id> → field data
db:<appId>:ix:<tableName>:<fieldName>:<val>→ index data
db:<appId>:grp:<groupName>                 → group info
friend:<selfAddress>:<friendAddress>        → friend relationship
sysgrp:<groupName>                          → system groups
```

For detailed key structure documentation, see `/x/dbchain/internal/keeper/db_key/readme.txt`.

## Development

### Building for Development

```bash
# Build with race detection
go build -race ./cmd/dbchaind

# Run tests
go test ./...

# Build specific module tests
go test ./x/dbchain/...
```

### Build Tags

The project uses Go build tags for conditional compilation:

- `sm2` - Enable SM2 cryptography
- `!sm2` - Use standard secp256k1 cryptography

## Known Limitations

- Delete and update operations for rows are currently not provided (see comments in codebase)
- Community edition is limited to 2 applications
- Requires minimum 4 nodes for BFT consensus in production

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

[Contribution guidelines not found in repository]

## Support

For issues and questions:
- Check the documentation in `/doc` (Chinese)
- Review the codebase comments and inline documentation

---

**Built with** [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and [Tendermint](https://github.com/tendermint/tendermint)
