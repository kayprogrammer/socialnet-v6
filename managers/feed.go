package managers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/models/choices"
	"github.com/kayprogrammer/socialnet-v6/schemas"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"github.com/pborman/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ----------------------------------
// POST MANAGEMENT
// --------------------------------
type PostManager struct {
}

func (obj PostManager) All(db *gorm.DB) []models.Post {
	posts := []models.Post{}
	db.Preload("AuthorObj").Preload("AuthorObj.AvatarObj").Preload("ImageObj").Preload("Reactions").Preload("Comments").Find(&posts).Order("created_at DESC")
	return posts
}

func (obj PostManager) Create(db *gorm.DB, author models.User, postData schemas.PostInputSchema) models.Post {
	id := uuid.Parse(uuid.New())
	// Create slug
	slug := slug.Make(fmt.Sprintf("%s %s %s", author.FirstName, author.LastName, id))
	base := models.BaseModel{ID: id}
	sub_base := models.FeedAbstract{BaseModel: base, Slug: slug, AuthorID: author.ID, Text: postData.Text}

	post := models.Post{FeedAbstract: sub_base}
	if postData.FileType != nil {
		file := models.File{ResourceType: *postData.FileType}
		post.ImageObj = &file
	}
	db.Create(&post)

	// Set related data
	post.AuthorObj = author
	return post
}

func (obj PostManager) GetBySlug(db *gorm.DB, slug string, opts ...bool) (*models.Post, *int, *utils.ErrorResponse) {
	post := models.Post{FeedAbstract: models.FeedAbstract{Slug: slug}}
	q := db
	if len(opts) > 0 { // Detailed param provided.
		q = db.Preload("AuthorObj.AvatarObj").Preload("Reactions").Preload("Comments")
	}
	q.Take(&post, post)
	if post.ID == nil {
		status_code := 404
		errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "Post does not exist")
		return nil, &status_code, &errData
	}
	return &post, nil, nil
}

func (obj PostManager) Update(db *gorm.DB, post *models.Post, postData schemas.PostInputSchema) *models.Post {
	if postData.FileType != nil {
		// Create or Update Image Object
		image := models.File{ResourceType: *postData.FileType}.UpdateOrCreate(db, post.ImageID)
		post.ImageObj = &image
	}
	post.Text = postData.Text
	db.Save(&post)
	return post
}

// func (obj PostManager) DropData(db *gorm.DB) {
// 	climodels.Post.Delete().ExecX(Ctx)
// }

// ----------------------------------
// COMMENT MANAGEMENT
// --------------------------------
type CommentManager struct {
}

func (obj CommentManager) GetBySlug(db *gorm.DB, slug string, opts ...bool) (*models.Comment, *int, *utils.ErrorResponse) {
	comment := models.Comment{FeedAbstract: models.FeedAbstract{Slug: slug}}
	q := db
	if len(opts) > 0 { // Detailed param provided.
		q = q.Preload(clause.Associations)
	}
	q.Take(&comment, comment)
	if comment.ID == nil {
		status_code := 404
		errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "Comment does not exist")
		return nil, &status_code, &errData
	}
	return &comment, nil, nil
}

func (obj CommentManager) GetByPostID(db *gorm.DB, postID uuid.UUID) []models.Comment {
	comments := []models.Comment{}
	db.Preload("Replies").Preload("Reactions").Preload("AuthorObj").Preload("AuthorObj.AvatarObj").Where(models.Comment{PostID: postID}).Find(&comments)
	return comments
}

func (obj CommentManager) Create(db *gorm.DB, author models.User, post models.Post, text string) models.Comment {
	id := uuid.Parse(uuid.New())
	// Create slug
	slug := slug.Make(fmt.Sprintf("%s %s %s", author.FirstName, author.LastName, id))
	base := models.BaseModel{ID: id}
	sub_base := models.FeedAbstract{BaseModel: base, Slug: slug, AuthorID: author.ID, Text: text}
	
	comment := models.Comment{FeedAbstract: sub_base, PostID: post.ID}
	db.Create(&comment)

	// Set related data
	comment.AuthorObj = author
	comment.PostObj = post
	return comment
}

func (obj CommentManager) Update(db *gorm.DB, comment models.Comment, author *models.User, text string) models.Comment {
	comment.Text = text
	db.Save(&comment)
	return comment
}

// func (obj CommentManager) DropData(db *gorm.DB) {
// 	climodels.Comment.Delete().ExecX(Ctx)
// }

// ----------------------------------
// REPLY MANAGEMENT
// --------------------------------
type ReplyManager struct {
}

func (obj ReplyManager) GetBySlug(db *gorm.DB, slug string, opts ...bool) (*models.Reply, *int, *utils.ErrorResponse) {
	reply := models.Reply{FeedAbstract: models.FeedAbstract{Slug: slug}}
	q := db
	if len(opts) > 0 { // Detailed param provided.
		q = q.Preload(clause.Associations)
	}
	q.Take(&reply, reply)
	if reply.ID == nil {
		status_code := 404
		errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "Reply does not exist")
		return nil, &status_code, &errData
	}
	return &reply, nil, nil
}

func (obj ReplyManager) Create(db *gorm.DB, author models.User, comment models.Comment, text string) models.Reply {
	id := uuid.Parse(uuid.New())
	// Create slug
	slug := slug.Make(fmt.Sprintf("%s %s %s", author.FirstName, author.LastName, id))
	base := models.BaseModel{ID: id}
	sub_base := models.FeedAbstract{BaseModel: base, Slug: slug, AuthorID: author.ID, Text: text}

	reply := models.Reply{FeedAbstract: sub_base, CommentID: comment.ID}
	db.Create(&reply)

	// Set related data
	reply.AuthorObj = author
	return reply
}

func (obj ReplyManager) Update(db *gorm.DB, reply models.Reply, author *models.User, text string) models.Reply {
	reply.Text = text
	db.Save(&reply)
	return reply
}

// func (obj ReplyManager) DropData(db *gorm.DB) {
// 	climodels.Reply.Delete().ExecX(Ctx)
// }

// ----------------------------------
// REACTIONS MANAGEMENT
// --------------------------------
type ReactionManager struct {
}

func (obj ReactionManager) GetReactionsQueryset(db *gorm.DB, fiberCtx *fiber.Ctx, focus string, slug string) ([]models.Reaction, *int, *utils.ErrorResponse) {
	reactions := []models.Reaction{}
	q := db.Preload("UserObj").Preload("UserObj.AvatarObj")
	if focus == "POST" {
		// Get Post Object and Query reactions for the post
		post, errCode, errData := PostManager{}.GetBySlug(db, slug)
		if errCode != nil {
			return nil, errCode, errData
		}
		q = q.Where(models.Reaction{Post: post})
	} else if focus == "COMMENT" {
		// Get Comment Object and Query reactions for the comment
		comment, errCode, errData := CommentManager{}.GetBySlug(db, slug)
		if errCode != nil {
			return nil, errCode, errData
		}
		q = q.Where(models.Reaction{Comment: comment})
	} else {
		// Get Reply Object and Query reactions for the reply
		reply, errCode, errData := ReplyManager{}.GetBySlug(db, slug)
		if errCode != nil {
			return nil, errCode, errData
		}
		q = q.Where(models.Reaction{Reply: reply})
	}

	// Filter by Reaction type if provided (e.g LIKE, LOVE)
	rtype := choices.ReactionChoice(fiberCtx.Query("reaction_type"))
	if len(rtype) > 0 {
		q = q.Where(models.Reaction{Rtype: rtype})
	}
	q.Find(&reactions)
	return reactions, nil, nil
}

func (obj ReactionManager) Update(db *gorm.DB, reaction models.Reaction, focus string, id uuid.UUID, rtype choices.ReactionChoice) models.Reaction {
	reaction.Rtype = rtype
	if focus == "POST" {
		reaction.PostID = &id
	} else if focus == "COMMENT" {
		reaction.CommentID = &id
	} else {
		reaction.ReplyID = &id
	}
	db.Save(&reaction)
	return reaction
}

func (obj ReactionManager) Create(db *gorm.DB, user models.User, focus string, focusID uuid.UUID, rtype choices.ReactionChoice) models.Reaction {
	reaction := models.Reaction{UserObj: user, Rtype: rtype}
	if focus == "POST" {
		reaction.PostID = &focusID
	} else if focus == "COMMENT" {
		reaction.CommentID = &focusID
	} else {
		reaction.ReplyID = &focusID
	}
	db.Create(&reaction)
	return reaction
}

func (obj ReactionManager) UpdateOrCreate(db *gorm.DB, user models.User, focus string, slug string, rtype choices.ReactionChoice) (*models.Reaction, *models.User, *int, *utils.ErrorResponse) {
	q := db.Preload("UserObj").Preload("UserObj.AvatarObj")
	var focusID *uuid.UUID
	var targetedObjAuthor *models.User
	reaction := models.Reaction{}
	if focus == "POST" {
		// Get Post Object and Query reactions for the post
		post, errCode, errData := PostManager{}.GetBySlug(db, slug, true)
		if errCode != nil {
			return nil, nil, errCode, errData
		}
		focusID = &post.ID
		q = q.Where(models.Reaction{Post: post})
		targetedObjAuthor = &post.AuthorObj
	} else if focus == "COMMENT" {
		// Get Comment Object and Query reactions for the comment
		comment, errCode, errData := CommentManager{}.GetBySlug(db, slug, true)
		if errCode != nil {
			return nil, nil, errCode, errData
		}
		focusID = &comment.ID
		q = q.Where(models.Reaction{Comment: comment})
		targetedObjAuthor = &comment.AuthorObj
	} else {
		// Get Reply Object and Query reactions for the reply
		reply, errCode, errData := ReplyManager{}.GetBySlug(db, slug, true)
		if errCode != nil {
			return nil, nil, errCode, errData
		}
		focusID = &reply.ID
		q = q.Where(models.Reaction{Reply: reply})
		targetedObjAuthor = &reply.AuthorObj
	}
	q.Take(&reaction, reaction)
	if reaction.ID == nil {
		// Create reaction
		reaction = obj.Create(db, user, focus, *focusID, rtype)
	} else {
		// Update
		reaction = obj.Update(db, reaction, focus, *focusID, rtype)
	}

	return &reaction, targetedObjAuthor, nil, nil
}

func (obj ReactionManager) GetByID(db *gorm.DB, id *uuid.UUID) (*models.Reaction, *int, *utils.ErrorResponse) {
	reaction := models.Reaction{}
	db.Preload(clause.Associations).Take(&reaction, models.Reaction{BaseModel: models.BaseModel{ID: *id}})
	if reaction.ID == nil {
		statusCode := 404
		errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "Reaction does not exist")
		return nil, &statusCode, &errData
	}
	return &reaction, nil, nil
}

// func (obj ReactionManager) DropData(db *gorm.DB) {
// 	climodels.Reaction.Delete().ExecX(Ctx)
// }
