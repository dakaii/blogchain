package keeper_test

import (
	"testing"

	"blogchain/x/blog/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestCommentOperations(t *testing.T) {
	f := initFixture(t)
	k := f.keeper
	ctx := sdk.UnwrapSDKContext(f.ctx)
	msgServer := NewMsgServerImpl(k)
	
	// Create a post first
	creator := GetTestAddress(t)
	createPostMsg := &types.MsgCreatePost{
		Creator: creator,
		Title:   "Post with Comments",
		Body:    "This post will receive comments",
		Tags:    []string{"test"},
	}
	postResp, err := msgServer.CreatePost(f.ctx, createPostMsg)
	require.NoError(t, err)
	postId := postResp.Id
	
	commenter1 := GetTestAddress(t)
	commenter2 := GetTestAddress(t)
	
	t.Run("Create root comment", func(t *testing.T) {
		msg := &types.MsgCreateComment{
			Creator:  commenter1,
			PostId:   postId,
			ParentId: 0, // Root comment
			Content:  "This is a great post!",
		}
		
		resp, err := msgServer.CreateComment(f.ctx, msg)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, uint64(0), resp.Id) // First comment ID
		
		// Verify comment was created
		comment, err := k.GetComment(ctx, resp.Id)
		require.NoError(t, err)
		require.Equal(t, commenter1, comment.Creator)
		require.Equal(t, postId, comment.PostId)
		require.Equal(t, uint64(0), comment.ParentId)
		require.Equal(t, "This is a great post!", comment.Content)
		require.Equal(t, uint32(0), comment.Depth)
		
		// Verify post comment count increased
		post, err := k.GetPost(ctx, postId)
		require.NoError(t, err)
		require.Equal(t, uint64(1), post.CommentCount)
	})
	
	t.Run("Create reply comment", func(t *testing.T) {
		// Reply to first comment
		msg := &types.MsgCreateComment{
			Creator:  commenter2,
			PostId:   postId,
			ParentId: 0, // Reply to first comment
			Content:  "I agree with this!",
		}
		
		resp, err := msgServer.CreateComment(f.ctx, msg)
		require.NoError(t, err)
		require.NotNil(t, resp)
		
		comment, err := k.GetComment(ctx, resp.Id)
		require.NoError(t, err)
		require.Equal(t, uint32(0), comment.Depth) // Still root since parent is 0
	})
	
	t.Run("Create nested reply", func(t *testing.T) {
		// First create a root comment
		rootMsg := &types.MsgCreateComment{
			Creator:  commenter1,
			PostId:   postId,
			ParentId: 0,
			Content:  "Root comment for nesting test",
		}
		rootResp, err := msgServer.CreateComment(f.ctx, rootMsg)
		require.NoError(t, err)
		rootId := rootResp.Id
		
		// Reply to root comment
		replyMsg := &types.MsgCreateComment{
			Creator:  commenter2,
			PostId:   postId,
			ParentId: rootId,
			Content:  "First level reply",
		}
		replyResp, err := msgServer.CreateComment(f.ctx, replyMsg)
		require.NoError(t, err)
		
		replyComment, err := k.GetComment(ctx, replyResp.Id)
		require.NoError(t, err)
		require.Equal(t, uint32(1), replyComment.Depth)
		
		// Reply to reply (second level)
		nestedMsg := &types.MsgCreateComment{
			Creator:  commenter1,
			PostId:   postId,
			ParentId: replyResp.Id,
			Content:  "Second level reply",
		}
		nestedResp, err := msgServer.CreateComment(f.ctx, nestedMsg)
		require.NoError(t, err)
		
		nestedComment, err := k.GetComment(ctx, nestedResp.Id)
		require.NoError(t, err)
		require.Equal(t, uint32(2), nestedComment.Depth)
	})
	
	t.Run("Cannot comment on deleted post", func(t *testing.T) {
		// Create and delete a post
		createMsg := &types.MsgCreatePost{
			Creator: creator,
			Title:   "Post to Delete",
			Body:    "This will be deleted",
		}
		postResp, err := msgServer.CreatePost(f.ctx, createMsg)
		require.NoError(t, err)
		
		deleteMsg := &types.MsgDeletePost{
			Creator: creator,
			Id:      postResp.Id,
		}
		_, err = msgServer.DeletePost(f.ctx, deleteMsg)
		require.NoError(t, err)
		
		// Try to comment on deleted post
		commentMsg := &types.MsgCreateComment{
			Creator:  commenter1,
			PostId:   postResp.Id,
			ParentId: 0,
			Content:  "Comment on deleted post",
		}
		_, err = msgServer.CreateComment(f.ctx, commentMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot comment on deleted post")
	})
	
	t.Run("Max depth enforcement", func(t *testing.T) {
		// Create a new post for depth testing
		createMsg := &types.MsgCreatePost{
			Creator: creator,
			Title:   "Depth Test Post",
			Body:    "Testing comment depth limits",
		}
		postResp, err := msgServer.CreatePost(f.ctx, createMsg)
		require.NoError(t, err)
		depthPostId := postResp.Id
		
		// Create a chain of comments up to max depth
		parentId := uint64(0)
		var lastCommentId uint64
		
		for i := 0; i <= 5; i++ { // MAX_COMMENT_DEPTH is 5
			msg := &types.MsgCreateComment{
				Creator:  commenter1,
				PostId:   depthPostId,
				ParentId: parentId,
				Content:  "Depth test comment",
			}
			
			if i <= 5 { // Should succeed up to depth 5
				resp, err := msgServer.CreateComment(f.ctx, msg)
				if i == 0 || parentId == 0 {
					// Root comment or when parentId is 0
					require.NoError(t, err)
					lastCommentId = resp.Id
					parentId = resp.Id
				} else if i <= 5 {
					require.NoError(t, err)
					lastCommentId = resp.Id
					parentId = resp.Id
				}
			}
		}
		
		// Try to exceed max depth
		msg := &types.MsgCreateComment{
			Creator:  commenter1,
			PostId:   depthPostId,
			ParentId: lastCommentId,
			Content:  "This should fail - too deep",
		}
		_, err = msgServer.CreateComment(f.ctx, msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "maximum comment depth")
	})
	
	t.Run("Update comment", func(t *testing.T) {
		// Create a comment
		createMsg := &types.MsgCreateComment{
			Creator:  commenter1,
			PostId:   postId,
			ParentId: 0,
			Content:  "Original comment content",
		}
		createResp, err := msgServer.CreateComment(f.ctx, createMsg)
		require.NoError(t, err)
		commentId := createResp.Id
		
		// Update the comment
		updateMsg := &types.MsgUpdateComment{
			Creator: commenter1,
			Id:      commentId,
			Content: "Updated comment content with edit",
		}
		_, err = msgServer.UpdateComment(f.ctx, updateMsg)
		require.NoError(t, err)
		
		// Verify update
		comment, err := k.GetComment(ctx, commentId)
		require.NoError(t, err)
		require.Equal(t, "Updated comment content with edit", comment.Content)
		require.GreaterOrEqual(t, comment.UpdatedAt, comment.CreatedAt) // May be same in tests
		
		// Non-creator cannot update
		updateMsg2 := &types.MsgUpdateComment{
			Creator: commenter2,
			Id:      commentId,
			Content: "Hacker trying to edit",
		}
		_, err = msgServer.UpdateComment(f.ctx, updateMsg2)
		require.Error(t, err)
		require.Contains(t, err.Error(), "only comment creator can update")
	})
	
	t.Run("Delete comment", func(t *testing.T) {
		// Create a comment
		createMsg := &types.MsgCreateComment{
			Creator:  commenter1,
			PostId:   postId,
			ParentId: 0,
			Content:  "Comment to delete",
		}
		createResp, err := msgServer.CreateComment(f.ctx, createMsg)
		require.NoError(t, err)
		commentId := createResp.Id
		
		// Get post comment count before deletion
		post, err := k.GetPost(ctx, postId)
		require.NoError(t, err)
		countBefore := post.CommentCount
		
		// Delete the comment
		deleteMsg := &types.MsgDeleteComment{
			Creator: commenter1,
			Id:      commentId,
		}
		_, err = msgServer.DeleteComment(f.ctx, deleteMsg)
		require.NoError(t, err)
		
		// Verify soft delete
		comment, err := k.GetComment(ctx, commentId)
		require.NoError(t, err)
		require.True(t, comment.Deleted)
		require.Greater(t, comment.DeletedAt, int64(0))
		
		// Verify post comment count decreased
		post, err = k.GetPost(ctx, postId)
		require.NoError(t, err)
		require.Equal(t, countBefore-1, post.CommentCount)
		
		// Cannot delete again
		_, err = msgServer.DeleteComment(f.ctx, deleteMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "comment already deleted")
		
		// Non-creator cannot delete
		createMsg2 := &types.MsgCreateComment{
			Creator:  commenter2,
			PostId:   postId,
			ParentId: 0,
			Content:  "Another comment",
		}
		createResp2, err := msgServer.CreateComment(f.ctx, createMsg2)
		require.NoError(t, err)
		
		deleteMsg2 := &types.MsgDeleteComment{
			Creator: commenter1, // Wrong creator
			Id:      createResp2.Id,
		}
		_, err = msgServer.DeleteComment(f.ctx, deleteMsg2)
		require.Error(t, err)
		require.Contains(t, err.Error(), "only comment creator can delete")
	})
	
	t.Run("Like comment", func(t *testing.T) {
		// Create a comment
		createMsg := &types.MsgCreateComment{
			Creator:  commenter1,
			PostId:   postId,
			ParentId: 0,
			Content:  "Likeable comment",
		}
		createResp, err := msgServer.CreateComment(f.ctx, createMsg)
		require.NoError(t, err)
		commentId := createResp.Id
		
		// Like the comment
		likeMsg := &types.MsgLikeComment{
			Liker:     commenter2,
			CommentId: commentId,
		}
		_, err = msgServer.LikeComment(f.ctx, likeMsg)
		require.NoError(t, err)
		
		// Verify like
		comment, err := k.GetComment(ctx, commentId)
		require.NoError(t, err)
		require.Equal(t, uint64(1), comment.Likes)
		
		// Cannot like twice
		_, err = msgServer.LikeComment(f.ctx, likeMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "comment already liked")
		
		// Another user can like
		likeMsg2 := &types.MsgLikeComment{
			Liker:     creator,
			CommentId: commentId,
		}
		_, err = msgServer.LikeComment(f.ctx, likeMsg2)
		require.NoError(t, err)
		
		comment, err = k.GetComment(ctx, commentId)
		require.NoError(t, err)
		require.Equal(t, uint64(2), comment.Likes)
	})
	
	t.Run("Get post comments", func(t *testing.T) {
		// Create a fresh post
		createPostMsg := &types.MsgCreatePost{
			Creator: creator,
			Title:   "Post for comment listing",
			Body:    "Testing comment retrieval",
		}
		postResp, err := msgServer.CreatePost(f.ctx, createPostMsg)
		require.NoError(t, err)
		testPostId := postResp.Id
		
		// Create multiple root comments
		for i := 0; i < 3; i++ {
			msg := &types.MsgCreateComment{
				Creator:  commenter1,
				PostId:   testPostId,
				ParentId: 0,
				Content:  "Root comment",
			}
			_, err := msgServer.CreateComment(f.ctx, msg)
			require.NoError(t, err)
		}
		
		// Get root comments
		comments, err := k.GetPostComments(ctx, testPostId, 0)
		require.NoError(t, err)
		require.Len(t, comments, 3)
		
		// Create reply to first comment
		if len(comments) > 0 {
			replyMsg := &types.MsgCreateComment{
				Creator:  commenter2,
				PostId:   testPostId,
				ParentId: comments[0].Id,
				Content:  "Reply to first",
			}
			_, err = msgServer.CreateComment(f.ctx, replyMsg)
			require.NoError(t, err)
			
			// Get replies
			replies, err := k.GetPostComments(ctx, testPostId, comments[0].Id)
			require.NoError(t, err)
			require.Len(t, replies, 1)
			require.Equal(t, "Reply to first", replies[0].Content)
		}
	})
}