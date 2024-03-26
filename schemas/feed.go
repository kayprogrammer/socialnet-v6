package schemas

import (
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/models/choices"
)

type PostInputSchema struct {
	Text				string		`json:"text" validate:"required" example:"God is good"`
	FileType			*string		`json:"file_type" example:"image/jpeg" validate:"omitempty,file_type_validator"`
}

// // REACTION SCHEMA
type ReactionInputSchema struct {
	Rtype		choices.ReactionChoice 			`json:"rtype" validate:"required,reaction_type_validator" example:"LIKE"`
}

// // COMMENTS & REPLIES SCHEMA
// type ReplySchema struct {
// 	Edges        		*ent.ReplyEdges 		`json:"edges,omitempty" swaggerignore:"true"`
// 	Author				UserDataSchema			`json:"author"`
// 	Slug				string					`json:"slug" example:"john-doe-d10dde64-a242-4ed0-bd75-4c759644b3a6"`
// 	Text				string					`json:"text" example:"Jesus Is King"`
// 	ReactionsCount 		uint					`json:"reactions_count" example:"200"`
// }

// func (reply ReplySchema) Init() ReplySchema {
// 	// Set Related Data.
// 	reply.Author = reply.Author.Init(reply.Edges.Author)
// 	reply.ReactionsCount = uint(len(reply.Edges.Reactions))
// 	reply.Edges = nil // Omit edges
// 	return reply
// }
type CommentInputSchema struct {
	Text			string 			`json:"text" example:"Jesus is Lord"`
}

// RESPONSE SCHEMAS
// POSTS
type PostsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items			[]models.Post		`json:"posts"`
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
	Data			models.Post		`json:"data"`
}

type PostsResponseSchema struct {
	ResponseSchema
	Data			PostsResponseDataSchema		`json:"data"`
}

type PostInputResponseSchema struct {
	ResponseSchema
	Data models.Post `json:"data"`
}

// REACTIONS
type ReactionsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items			[]models.Reaction		`json:"reactions"`
}

type ReactionsResponseSchema struct {
	ResponseSchema
	Data			ReactionsResponseDataSchema		`json:"data"`
}

type ReactionResponseSchema struct {
	ResponseSchema
	Data			models.Reaction		`json:"data"`
}

// COMMENTS & REPLIES
type CommentWithRepliesResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items			[]models.Reply		`json:"items"`
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
	Comment			models.Comment								`json:"comment"`
	Replies			CommentWithRepliesResponseDataSchema		`json:"replies"`
}

type CommentsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items		[]models.Comment				`json:"comments"`
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
	Data			models.Comment			`json:"data"`
}

type CommentWithRepliesResponseSchema struct {
	ResponseSchema
	Data			CommentWithRepliesSchema			`json:"data"`
}

type ReplyResponseSchema struct {
	ResponseSchema
	Data			models.Reply			`json:"data"`
}