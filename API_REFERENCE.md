# BlogChain API Reference

## 📖 Overview

BlogChain provides both REST and gRPC APIs for interacting with the decentralized blogging platform. All APIs are automatically generated from protobuf definitions, ensuring type safety and consistency.

## 🔗 Base URLs

- **REST API**: `http://localhost:1317/blogchain/blog/v1/`
- **gRPC API**: `localhost:9090`
- **Tendermint RPC**: `http://localhost:26657`

## 🔐 Authentication

BlogChain uses Cosmos SDK's standard transaction signing mechanism. Users authenticate by signing transactions with their private keys.

```javascript
// Example transaction signing
const tx = await client.signAndBroadcast(
  signerAddress,
  [message],
  fee,
  memo
)
```

## 📝 Posts API

### Create Post
**Transaction**: `MsgCreatePost`

```protobuf
message MsgCreatePost {
  string creator = 1;
  string title = 2;
  string body = 3;
  repeated string tags = 4;
}
```

**REST Example**:
```http
POST /cosmos/tx/v1beta1/txs
Content-Type: application/json

{
  "tx": {
    "body": {
      "messages": [{
        "@type": "/blogchain.blog.v1.MsgCreatePost",
        "creator": "cosmos1...",
        "title": "My First Post",
        "body": "This is the content of my post",
        "tags": ["blockchain", "cosmos"]
      }]
    }
  }
}
```

### Update Post
**Transaction**: `MsgUpdatePost`

```protobuf
message MsgUpdatePost {
  string creator = 1;
  uint64 id = 2;
  string title = 3;
  string body = 4;
  repeated string tags = 5;
  repeated string media_blob_ids = 6;
  string content_blob_id = 7;
}
```

### Delete Post
**Transaction**: `MsgDeletePost`

```protobuf
message MsgDeletePost {
  string creator = 1;
  uint64 id = 2;
}
```

### Like Post
**Transaction**: `MsgLikePost`

```protobuf
message MsgLikePost {
  string liker = 1;
  uint64 post_id = 2;
}
```

### Query Posts
**Endpoint**: `GET /blogchain/blog/v1/posts`

**Parameters**:
- `pagination.limit` (uint64): Number of posts per page (max 100)
- `pagination.offset` (uint64): Number of posts to skip
- `pagination.count_total` (bool): Whether to include total count

**Response**:
```json
{
  "posts": [
    {
      "id": "1",
      "creator": "cosmos1...",
      "title": "Example Post",
      "body": "Post content",
      "tags": ["example"],
      "created_at": "1234567890",
      "updated_at": "1234567890",
      "likes": "5",
      "comment_count": "3",
      "deleted": false
    }
  ],
  "pagination": {
    "total": "100",
    "limit": "10",
    "offset": "0"
  }
}
```

### Query Single Post
**Endpoint**: `GET /blogchain/blog/v1/posts/{id}`

**Response**:
```json
{
  "post": {
    "id": "1",
    "creator": "cosmos1...",
    "title": "Example Post",
    "body": "Post content",
    "tags": ["example"],
    "created_at": "1234567890",
    "updated_at": "1234567890",
    "likes": "5",
    "comment_count": "3",
    "deleted": false
  }
}
```

## 💬 Comments API

### Create Comment
**Transaction**: `MsgCreateComment`

```protobuf
message MsgCreateComment {
  string creator = 1;
  uint64 post_id = 2;
  uint64 parent_id = 3; // 0 for root comment
  string content = 4;
}
```

### Update Comment
**Transaction**: `MsgUpdateComment`

```protobuf
message MsgUpdateComment {
  string creator = 1;
  uint64 id = 2;
  string content = 3;
}
```

### Delete Comment
**Transaction**: `MsgDeleteComment`

```protobuf
message MsgDeleteComment {
  string creator = 1;
  uint64 id = 2;
}
```

### Like Comment
**Transaction**: `MsgLikeComment`

```protobuf
message MsgLikeComment {
  string liker = 1;
  uint64 comment_id = 2;
}
```

### Query Comments
**Endpoint**: `GET /blogchain/blog/v1/posts/{post_id}/comments`

**Parameters**:
- `parent_id` (uint64): Parent comment ID (0 for root comments)
- `pagination.limit` (uint64): Number of comments per page
- `pagination.offset` (uint64): Number of comments to skip

**Response**:
```json
{
  "comments": [
    {
      "id": "1",
      "post_id": "1",
      "parent_id": "0",
      "creator": "cosmos1...",
      "content": "Great post!",
      "created_at": "1234567890",
      "updated_at": "1234567890",
      "likes": "2",
      "depth": 0,
      "deleted": false
    }
  ],
  "pagination": {
    "total": "5",
    "limit": "10",
    "offset": "0"
  }
}
```

### Query Comment Thread
**Endpoint**: `GET /blogchain/blog/v1/comments/{id}/thread`

**Parameters**:
- `max_depth` (uint32): Maximum depth to fetch (0 = unlimited, max 10)

**Response**:
```json
{
  "thread": {
    "comment": {
      "id": "1",
      "post_id": "1",
      "parent_id": "0",
      "creator": "cosmos1...",
      "content": "Root comment",
      "created_at": "1234567890",
      "likes": "3",
      "depth": 0
    },
    "replies": [
      {
        "comment": {
          "id": "2",
          "post_id": "1",
          "parent_id": "1",
          "creator": "cosmos1...",
          "content": "Reply to root",
          "depth": 1
        },
        "replies": []
      }
    ]
  }
}
```

## 👤 Profiles API

### Create Profile
**Transaction**: `MsgCreateProfile`

```protobuf
message MsgCreateProfile {
  string creator = 1;
  string username = 2;
  string display_name = 3;
  string bio = 4;
  string avatar_url = 5;
  string website = 6;
}
```

**Validation Rules**:
- `username`: 3-20 characters, alphanumeric and underscore only
- `display_name`: max 50 characters
- `bio`: max 500 characters
- `avatar_url`: max 500 characters
- `website`: max 200 characters, valid URL format

### Update Profile
**Transaction**: `MsgUpdateProfile`

```protobuf
message MsgUpdateProfile {
  string creator = 1;
  string display_name = 2;
  string bio = 3;
  string avatar_url = 4;
  string website = 5;
}
```

### Follow User
**Transaction**: `MsgFollow`

```protobuf
message MsgFollow {
  string follower = 1;
  string following = 2;
}
```

### Unfollow User
**Transaction**: `MsgUnfollow`

```protobuf
message MsgUnfollow {
  string follower = 1;
  string following = 2;
}
```

### Query Profile by Address
**Endpoint**: `GET /blogchain/blog/v1/profiles/{address}`

**Response**:
```json
{
  "profile": {
    "address": "cosmos1...",
    "username": "alice_crypto",
    "display_name": "Alice",
    "bio": "Blockchain enthusiast",
    "avatar_url": "https://example.com/avatar.png",
    "website": "https://alice.example.com",
    "created_at": "1234567890",
    "updated_at": "1234567890",
    "followers": "150",
    "following": "75",
    "post_count": "42",
    "verified": false
  }
}
```

### Query Profile by Username
**Endpoint**: `GET /blogchain/blog/v1/profiles/username/{username}`

### Query All Profiles
**Endpoint**: `GET /blogchain/blog/v1/profiles`

**Parameters**:
- `pagination.limit` (uint64): Number of profiles per page
- `pagination.offset` (uint64): Number of profiles to skip

### Query Followers
**Endpoint**: `GET /blogchain/blog/v1/profiles/{address}/followers`

**Response**:
```json
{
  "followers": [
    "cosmos1...",
    "cosmos1...",
    "cosmos1..."
  ],
  "pagination": {
    "total": "150"
  }
}
```

### Query Following
**Endpoint**: `GET /blogchain/blog/v1/profiles/{address}/following`

### Check Follow Status
**Endpoint**: `GET /blogchain/blog/v1/profiles/{follower}/following/{following}`

**Response**:
```json
{
  "is_following": true
}
```

## 🔍 Error Responses

All APIs return standardized error responses:

```json
{
  "code": 5,
  "message": "post not found",
  "details": []
}
```

### Common Error Codes
- `2` - Unknown
- `3` - Invalid Argument
- `5` - Not Found
- `7` - Permission Denied
- `13` - Internal Error
- `16` - Unauthenticated

### Validation Errors
```json
{
  "code": 3,
  "message": "invalid request: comment content cannot be empty",
  "details": []
}
```

## 📡 Real-time Updates (Planned)

### WebSocket Events
Connection: `ws://localhost:26657/websocket`

#### Event Types
- `post_created` - New post created
- `post_updated` - Post updated
- `post_deleted` - Post deleted
- `post_liked` - Post liked
- `comment_created` - New comment
- `comment_updated` - Comment updated
- `comment_deleted` - Comment deleted
- `comment_liked` - Comment liked
- `profile_created` - New profile
- `profile_updated` - Profile updated
- `user_followed` - User followed
- `user_unfollowed` - User unfollowed

#### Event Structure
```json
{
  "type": "post_created",
  "attributes": {
    "id": "123",
    "creator": "cosmos1...",
    "title": "New Post"
  },
  "height": "1000",
  "tx_hash": "ABC123..."
}
```

## 🛠️ SDK Integration

### JavaScript/TypeScript
```bash
npm install @cosmjs/stargate @cosmjs/proto-signing
```

```javascript
import { SigningStargateClient } from "@cosmjs/stargate"

// Initialize client
const client = await SigningStargateClient.connectWithSigner(
  "http://localhost:26657",
  signer
)

// Create post
const msg = {
  typeUrl: "/blogchain.blog.v1.MsgCreatePost",
  value: {
    creator: address,
    title: "My Post",
    body: "Content",
    tags: ["blockchain"]
  }
}

const result = await client.signAndBroadcast(address, [msg], "auto")
```

### Go Client
```go
import (
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/tx"
    "blogchain/x/blog/types"
)

// Create post message
msg := &types.MsgCreatePost{
    Creator: creator,
    Title:   "My Post",
    Body:    "Content",
    Tags:    []string{"blockchain"},
}

// Build and broadcast transaction
txBuilder := clientCtx.TxConfig.NewTxBuilder()
txBuilder.SetMsgs(msg)

txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
result, err := clientCtx.BroadcastTx(txBytes)
```

## 📊 Rate Limits

Current rate limits (subject to blockchain capacity):
- **Transactions**: Limited by block time and gas
- **Queries**: No explicit limits (consider implementing)
- **WebSocket**: Connection limits based on node capacity

## 🔒 Security Considerations

### Transaction Security
- All transactions require valid signatures
- Replay protection through sequence numbers
- Gas fees prevent spam attacks

### Query Security
- Read-only operations
- Input validation on all parameters
- No sensitive data exposure in public queries

### Best Practices
1. Always validate inputs on client side
2. Use HTTPS in production
3. Implement proper error handling
4. Cache responses appropriately
5. Monitor for unusual activity patterns

## 📚 Additional Resources

- **Cosmos SDK Documentation**: https://docs.cosmos.network/
- **gRPC Documentation**: https://grpc.io/docs/
- **Protocol Buffers**: https://developers.google.com/protocol-buffers
- **Tendermint Core**: https://docs.tendermint.com/

## 🚀 Future API Enhancements

### Planned Features
1. **GraphQL API**: Flexible query language
2. **Batch Operations**: Multiple operations in single request
3. **Subscriptions**: Real-time data subscriptions
4. **Advanced Filtering**: Complex query filters
5. **Analytics Endpoints**: Platform statistics and metrics

### API Versioning
The API follows semantic versioning. Current version is `v1`. Future versions will maintain backward compatibility where possible.

### Deprecation Policy
- 6 months notice for breaking changes
- Parallel support for old and new versions
- Clear migration documentation