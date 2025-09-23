# BlogChain Enhancement Implementation Plan

## Overview
This document outlines the phased implementation of key improvements to the BlogChain decentralized blogging platform, including integration with Walrus for decentralized media storage.

## Phase 1: Core Post Management (Week 1-2)

### 1.1 Post Edit/Update Functionality
**Backend Changes:**
```proto
// Add to proto/blogchain/blog/v1/tx.proto
message MsgUpdatePost {
  string creator = 1;
  uint64 id = 2;
  string title = 3;
  string body = 4;
  repeated string tags = 5;
}
```

**Implementation:**
- [ ] Add `MsgUpdatePost` to tx.proto
- [ ] Implement keeper method `UpdatePost` with ownership validation
- [ ] Add `updated_at` timestamp field to Post
- [ ] Create tests for update functionality
- [ ] Update frontend to show edit button for post creators

### 1.2 Post Deletion (Soft Delete)
**Backend Changes:**
```proto
message MsgDeletePost {
  string creator = 1;
  uint64 id = 2;
}

// Add to Post message
bool deleted = 8;
int64 deleted_at = 9;
```

**Implementation:**
- [ ] Add `MsgDeletePost` to tx.proto
- [ ] Implement soft delete in keeper
- [ ] Filter deleted posts from queries
- [ ] Add restore capability for future moderation
- [ ] Update frontend to handle deleted posts

## Phase 2: Comments System (Week 2-3)

### 2.1 Create Comments Module
**Structure:**
```
x/comments/
├── keeper/
│   ├── keeper.go
│   ├── msg_server.go
│   ├── query.go
│   └── comment.go
├── types/
│   ├── comment.proto
│   ├── tx.proto
│   └── query.proto
└── module/
    └── module.go
```

**Proto Definition:**
```proto
message Comment {
  uint64 id = 1;
  uint64 post_id = 2;
  uint64 parent_id = 3; // For nested comments
  string creator = 4;
  string content = 5;
  int64 created_at = 6;
  uint64 likes = 7;
  bool deleted = 8;
}

message MsgCreateComment {
  string creator = 1;
  uint64 post_id = 2;
  uint64 parent_id = 3;
  string content = 4;
}
```

**Implementation:**
- [ ] Create comments module structure
- [ ] Define proto messages
- [ ] Implement keeper logic
- [ ] Add comment count to posts
- [ ] Create pagination for comments
- [ ] Add frontend comment components

## Phase 3: User Profiles (Week 3-4)

### 3.1 Profile Module
**Proto Definition:**
```proto
message UserProfile {
  string address = 1;
  string username = 2;
  string bio = 3;
  string avatar_url = 4; // Walrus blob ID
  string website = 5;
  int64 created_at = 6;
  int64 updated_at = 7;
  uint64 followers = 8;
  uint64 following = 9;
}

message MsgCreateProfile {
  string creator = 1;
  string username = 2;
  string bio = 3;
  string avatar_blob_id = 4; // Walrus blob ID
}
```

**Implementation:**
- [ ] Create profiles module
- [ ] Username uniqueness validation
- [ ] Follow/unfollow functionality
- [ ] Profile query endpoints
- [ ] Frontend profile pages

## Phase 4: Walrus Integration (Week 4-5)

### 4.1 Walrus Storage Service
**Architecture:**
```
Frontend → Upload to Walrus → Store blob ID on chain

Components:
1. Frontend Walrus service
2. Backend blob ID validation
3. HTTP gateway for retrieval
```

**Frontend Service:**
```typescript
// services/walrus.ts
class WalrusService {
  private aggregatorUrl = 'https://walrus-testnet.sui.io'
  private publisherUrl = 'https://walrus-publisher.testnet.sui.io'
  
  async uploadBlob(file: File): Promise<string> {
    // Upload to Walrus
    // Return blob ID
  }
  
  getBlobUrl(blobId: string): string {
    return `${this.aggregatorUrl}/v1/${blobId}`
  }
}
```

**Backend Integration:**
```proto
message Post {
  // ... existing fields
  repeated string media_blob_ids = 10; // Walrus blob IDs
  string content_blob_id = 11; // For large content
}
```

**Implementation:**
- [ ] Create Walrus service in frontend
- [ ] Add blob ID fields to Post/Profile protos
- [ ] Implement upload flow in frontend
- [ ] Add blob URL resolution
- [ ] Create fallback for IPFS gateway
- [ ] Add CDN caching layer

### 4.2 Content Size Optimization
**Strategy:**
- Small posts (<10KB): Store on-chain
- Large posts (>10KB): Store on Walrus
- Media files: Always on Walrus

## Phase 5: Search & Discovery (Week 5-6)

### 5.1 Enhanced Queries
**New Query Endpoints:**
```proto
service Query {
  rpc SearchPosts(SearchPostsRequest) returns (SearchPostsResponse);
  rpc PostsByTag(PostsByTagRequest) returns (PostsByTagResponse);
  rpc PostsByAuthor(PostsByAuthorRequest) returns (PostsByAuthorResponse);
  rpc TrendingPosts(TrendingPostsRequest) returns (TrendingPostsResponse);
}
```

**Implementation:**
- [ ] Add tag indexing in keeper
- [ ] Implement search with filters
- [ ] Create trending algorithm (likes + recency)
- [ ] Add related posts query
- [ ] Frontend search interface

## Phase 6: Rich Text Editor (Week 6)

### 6.1 Markdown Support
**Features:**
- Markdown editing with preview
- Code syntax highlighting
- Image embedding (via Walrus)
- Tables and lists support

**Implementation:**
- [ ] Add markdown parser library
- [ ] Create editor component with TipTap
- [ ] Add preview mode
- [ ] Integrate Walrus for image uploads
- [ ] Add code highlighting with Prism.js

## Testing Strategy

### Unit Tests
- [ ] Keeper method tests (100% coverage)
- [ ] Message validation tests
- [ ] Query pagination tests

### Integration Tests
- [ ] End-to-end post CRUD operations
- [ ] Comment threading tests
- [ ] Profile creation and updates
- [ ] Walrus upload/retrieval

### Frontend Tests
- [ ] Component tests with Vitest
- [ ] Service layer tests
- [ ] E2E tests with Playwright

## Deployment Plan

### Testnet Deployment
1. Deploy updated contracts to testnet
2. Run migration scripts
3. Test Walrus integration on testnet
4. Frontend deployment to staging

### Mainnet Deployment
1. Governance proposal for upgrade
2. Coordinated upgrade at block height
3. Post-upgrade validation
4. Frontend production deployment

## Performance Considerations

### Optimization Strategies
1. **Query Optimization**
   - Add pagination to all list queries
   - Implement cursor-based pagination for feeds
   - Cache frequently accessed data

2. **Storage Optimization**
   - Use Walrus for media and large content
   - Compress text before storage
   - Implement blob garbage collection

3. **Frontend Optimization**
   - Lazy loading for posts
   - Virtual scrolling for long lists
   - Service worker for offline support

## Security Measures

### Input Validation
- [ ] Content size limits (100KB for on-chain)
- [ ] HTML sanitization
- [ ] Rate limiting per address
- [ ] Spam detection algorithm

### Access Control
- [ ] Owner-only updates/deletes
- [ ] Profile uniqueness enforcement
- [ ] Comment moderation tools

## Monitoring & Analytics

### Metrics to Track
- Post creation rate
- Active users (daily/monthly)
- Storage usage (on-chain vs Walrus)
- Query performance
- Error rates

### Implementation
- [ ] Add telemetry to keeper methods
- [ ] Frontend analytics integration
- [ ] Grafana dashboard setup
- [ ] Alert configuration

## Timeline Summary

| Phase | Duration | Key Deliverables |
|-------|----------|-----------------|
| Phase 1 | Week 1-2 | Edit/Delete posts |
| Phase 2 | Week 2-3 | Comments system |
| Phase 3 | Week 3-4 | User profiles |
| Phase 4 | Week 4-5 | Walrus integration |
| Phase 5 | Week 5-6 | Search & Discovery |
| Phase 6 | Week 6 | Rich text editor |

## Success Criteria

1. All tests passing with >90% coverage
2. Walrus integration working on testnet
3. <2s page load time
4. Zero critical security issues
5. Positive user feedback from beta testing

## Next Steps

1. Begin Phase 1 implementation
2. Set up Walrus testnet account
3. Create development branch
4. Update project documentation