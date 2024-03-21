package schemas

import (
	"time"

	"github.com/google/uuid"
)

type PostSchema struct {
	Author 				UserDataSchema 		`json:"author"`
	Text 				string 				`json:"text" example:"God Is Good"`
	Slug 				string				`json:"slug" example:"john-doe-d10dde64-a242-4ed0-bd75-4c759644b3a6"`
	ReactionsCount 		uint				`json:"reactions_count" example:"200"`
	CommentsCount 		uint				`json:"comments_count" example:"35"`
	Image 				*string				`json:"image" example:"https://img.url"`
	CreatedAt 			time.Time			`json:"created_at" example:"2024-01-14T19:00:02.613124+01:00"`
	UpdatedAt 			time.Time			`json:"updated_at" example:"2024-01-14T19:00:02.613124+01:00"`
}

func (post PostSchema) Init () PostSchema {
	// Set Author Details.
	post.Author = post.Author.Init(post.Edges.Author)

	// Set ImageUrl
	image := post.Edges.Image
	if image != nil {
		url := utils.GenerateFileUrl(image.ID.String(), "posts", image.ResourceType)
		post.Image = &url
	}

	// Set Reactions & Comments Count
	post.ReactionsCount = uint(len(post.Edges.Reactions))
	post.CommentsCount = uint(len(post.Edges.Comments))
	post.Edges = nil // Omit edges
	return post
}

type PostInputSchema struct {
	Text				string		`json:"text" validate:"required" example:"God is good"`
	FileType			*string		`json:"file_type" example:"image/jpeg" validate:"omitempty,file_type_validator"`
}

// REACTION SCHEMA
type ReactionSchema struct {
	Edges        	*ent.ReactionEdges 	`json:"edges,omitempty" swaggerignore:"true"`
    ID 				uuid.UUID			`json:"id" example:"d10dde64-a242-4ed0-bd75-4c759644b3a6"`
    User 			UserDataSchema		`json:"user"`
    Rtype 			string				`json:"rtype" example:"LIKE"`
}

func (reaction ReactionSchema) Init() ReactionSchema {
	// Set User Details.
	reaction.User = reaction.User.Init(reaction.Edges.User)

	reaction.Edges = nil // Omit edges
	return reaction
}

type ReactionInputSchema struct {
	Rtype		reaction.Rtype 			`json:"rtype" validate:"required,reaction_type_validator" example:"LIKE"`
}

// COMMENTS & REPLIES SCHEMA
type ReplySchema struct {
	Edges        		*ent.ReplyEdges 		`json:"edges,omitempty" swaggerignore:"true"`
	Author				UserDataSchema			`json:"author"`
	Slug				string					`json:"slug" example:"john-doe-d10dde64-a242-4ed0-bd75-4c759644b3a6"`
	Text				string					`json:"text" example:"Jesus Is King"`
	ReactionsCount 		uint					`json:"reactions_count" example:"200"`
}

func (reply ReplySchema) Init() ReplySchema {
	// Set Related Data.
	reply.Author = reply.Author.Init(reply.Edges.Author)
	reply.ReactionsCount = uint(len(reply.Edges.Reactions))
	reply.Edges = nil // Omit edges
	return reply
}

type CommentSchema struct {
	ReplySchema
	Edges        			*ent.CommentEdges 		`json:"edges,omitempty" swaggerignore:"true"`
	RepliesCount			uint					`json:"replies_count" example:"50"`
}

func (comment CommentSchema) Init() CommentSchema {
	// Set Related Data.
	comment.Author = comment.Author.Init(comment.Edges.Author)
	comment.ReactionsCount = uint(len(comment.Edges.Reactions))
	comment.RepliesCount = uint(len(comment.Edges.Replies))
	comment.Edges = nil // Omit edges
	return comment
}

type CommentInputSchema struct {
	Text			string 			`json:"text" example:"Jesus is Lord"`
}

// RESPONSE SCHEMAS
// POSTS
type PostsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items			[]PostSchema		`json:"posts"`
}

func (data PostsResponseDataSchema) Init () PostsResponseDataSchema {
	// Set Initial Data
	items := data.Items
	for i := range items {
		items[i] = items[i].Init()
	}
	data.Items = items
	return data
}

type PostResponseSchema struct {
	ResponseSchema
	Data			PostSchema		`json:"data"`
}

type PostsResponseSchema struct {
	ResponseSchema
	Data			PostsResponseDataSchema		`json:"data"`
}

type PostInputResponseDataSchema struct {
	PostSchema
	Image 				*string					`json:"image,omitempty" swaggerignore:"true"` // Remove image during create & update 
	FileUploadData 		*utils.SignatureFormat 	`json:"file_upload_data"`
}

func (postData PostInputResponseDataSchema) Init(fileType *string) PostInputResponseDataSchema {
	image := postData.PostSchema.Edges.Image
	if fileType != nil && image != nil { // Generate data when file is being uploaded
		fuData := utils.GenerateFileSignature(image.ID.String(), "posts")
		postData.FileUploadData = &fuData
	}
	postData.PostSchema = postData.PostSchema.Init()
	return postData	
}

type PostInputResponseSchema struct {
	ResponseSchema
	Data PostInputResponseDataSchema `json:"data"`
}

// REACTIONS
type ReactionsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items			[]ReactionSchema		`json:"reactions"`
}

func (data ReactionsResponseDataSchema) Init () ReactionsResponseDataSchema {
	// Set Initial Data
	items := data.Items
	for i := range items {
		items[i] = items[i].Init()
	}
	data.Items = items
	return data
}
type ReactionsResponseSchema struct {
	ResponseSchema
	Data			ReactionsResponseDataSchema		`json:"data"`
}

type ReactionResponseSchema struct {
	ResponseSchema
	Data			ReactionSchema		`json:"data"`
}

// COMMENTS & REPLIES
type CommentWithRepliesResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items			[]ReplySchema		`json:"items"`
}

func (data CommentWithRepliesResponseDataSchema) Init () CommentWithRepliesResponseDataSchema {
	// Set Initial Data
	items := data.Items
	for i := range items {
		items[i] = items[i].Init()
	}
	data.Items = items
	return data
}

type CommentWithRepliesSchema struct {
	Comment			CommentSchema								`json:"comment"`
	Replies			CommentWithRepliesResponseDataSchema		`json:"replies"`
}

type CommentsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items		[]CommentSchema				`json:"comments"`
}

func (data CommentsResponseDataSchema) Init () CommentsResponseDataSchema {
	// Set Initial Data
	items := data.Items
	for i := range items {
		items[i] = items[i].Init()
	}
	data.Items = items
	return data
}

type CommentsResponseSchema struct {
	ResponseSchema
	Data			CommentsResponseDataSchema		`json:"data"`
}

type CommentResponseSchema struct {
	ResponseSchema
	Data			CommentSchema			`json:"data"`
}

type CommentWithRepliesResponseSchema struct {
	ResponseSchema
	Data			CommentWithRepliesSchema			`json:"data"`
}

type ReplyResponseSchema struct {
	ResponseSchema
	Data			ReplySchema			`json:"data"`
}