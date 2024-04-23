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

func AuthorAvatarScope(db *gorm.DB) *gorm.DB {
	return db.Joins("AuthorObj").Joins("AuthorObj.AvatarObj")
}

func AuthorReactionScope(db *gorm.DB) *gorm.DB {
	return db.Scopes(AuthorAvatarScope).Preload("Reactions")
}

// ----------------------------------
// POST MANAGEMENT
// --------------------------------
type PostManager struct {
}

func (obj PostManager) All(db *gorm.DB) []models.Post {
	posts := []models.Post{}
	db.Scopes(AuthorReactionScope).Joins("ImageObj").Preload("Comments").Find(&posts).Order("created_at DESC")
	return posts
}

func (obj PostManager) Create(db *gorm.DB, author models.User, postData schemas.PostInputSchema) models.Post {
	id := uuid.Parse(uuid.New())
	// Create slug
	slug := slug.Make(fmt.Sprintf("%s %s %s", author.FirstName, author.LastName, id))
	base := models.BaseModel{ID: id}
	sub_base := models.FeedAbstract{BaseModel: base, Slug: slug, AuthorObj: author, AuthorID: author.ID, Text: postData.Text}

	post := models.Post{FeedAbstract: sub_base}
	if postData.FileType != nil {
		file := models.File{ResourceType: *postData.FileType}
		post.ImageObj = &file
	}
	db.Create(&post)
	return post
}

func (obj PostManager) GetBySlug(db *gorm.DB, slug string, opts ...bool) (*models.Post, *int, *utils.ErrorResponse) {
	post := models.Post{FeedAbstract: models.FeedAbstract{Slug: slug}}
	q := db.Scopes(AuthorReactionScope)
	if len(opts) > 0 { // Detailed param provided.
		q = q.Preload("Comments")
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
	db.Omit(clause.Associations).Save(&post)
	return post
}

func (obj PostManager) DropData(db *gorm.DB) {
	db.Delete(&models.Post{})
}

// ----------------------------------
// COMMENT MANAGEMENT
// --------------------------------
type CommentManager struct {
}

func (obj CommentManager) GetBySlug(db *gorm.DB, slug string, opts ...bool) (*models.Comment, *int, *utils.ErrorResponse) {
	comment := models.Comment{FeedAbstract: models.FeedAbstract{Slug: slug}}
	q := db.Scopes(AuthorAvatarScope)
	if len(opts) > 0 { // Detailed param provided.
		q = q.Preload("Reactions").Preload("Replies").Preload("Replies.AuthorObj").Preload("Replies.AuthorObj.AvatarObj")
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
	db.Preload("Replies").Scopes(AuthorReactionScope).Where(models.Comment{PostID: postID}).Find(&comments)
	return comments
}

func (obj CommentManager) Create(db *gorm.DB, author models.User, post models.Post, text string) models.Comment {
	id := uuid.Parse(uuid.New())
	// Create slug
	slug := slug.Make(fmt.Sprintf("%s %s %s", author.FirstName, author.LastName, id))
	base := models.BaseModel{ID: id}
	sub_base := models.FeedAbstract{BaseModel: base, Slug: slug, AuthorID: author.ID, AuthorObj: author, Text: text}

	comment := models.Comment{FeedAbstract: sub_base, PostID: post.ID, PostObj: post}
	db.Create(&comment)
	return comment
}

func (obj CommentManager) Update(db *gorm.DB, comment models.Comment, author *models.User, text string) models.Comment {
	comment.Text = text
	db.Omit(clause.Associations).Save(&comment)
	return comment
}

func (obj CommentManager) DropData(db *gorm.DB) {
	db.Delete(&models.Comment{})
}

// ----------------------------------
// REPLY MANAGEMENT
// --------------------------------
type ReplyManager struct {
}

func (obj ReplyManager) GetBySlug(db *gorm.DB, slug string, opts ...bool) (*models.Reply, *int, *utils.ErrorResponse) {
	reply := models.Reply{FeedAbstract: models.FeedAbstract{Slug: slug}}
	q := db
	if len(opts) > 0 { // Detailed param provided.
		q = q.Scopes(AuthorReactionScope)
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
	sub_base := models.FeedAbstract{BaseModel: base, Slug: slug, AuthorID: author.ID, AuthorObj: author, Text: text}

	reply := models.Reply{FeedAbstract: sub_base, CommentID: comment.ID, CommentObj: comment}
	db.Create(&reply)
	return reply
}

func (obj ReplyManager) Update(db *gorm.DB, reply models.Reply, author *models.User, text string) models.Reply {
	reply.Text = text
	db.Save(&reply)
	return reply
}

func (obj ReplyManager) DropData(db *gorm.DB) {
	db.Delete(&models.Reply{})
}

// ----------------------------------
// REACTIONS MANAGEMENT
// --------------------------------
func UserAvatarReactionScope(db *gorm.DB) *gorm.DB {
	return db.Joins("UserObj").Joins("UserObj.AvatarObj")
}

type ReactionManager struct {
}

func (obj ReactionManager) GetReactionsQueryset(db *gorm.DB, fiberCtx *fiber.Ctx, focus choices.FocusTypeChoice, slug string) ([]models.Reaction, *int, *utils.ErrorResponse) {
	reactions := []models.Reaction{}
	q := db.Scopes(UserAvatarReactionScope)
	if focus == choices.FTPOST {
		// Get Post Object and Query reactions for the post
		post, errCode, errData := PostManager{}.GetBySlug(db, slug)
		if errCode != nil {
			return nil, errCode, errData
		}
		q = q.Where(models.Reaction{Post: post})
	} else if focus == choices.FTCOMMENT {
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

func (obj ReactionManager) Update(db *gorm.DB, reaction models.Reaction, focus choices.FocusTypeChoice, post *models.Post, comment *models.Comment, reply *models.Reply, rtype choices.ReactionChoice) models.Reaction {
	reaction.Rtype = rtype
	if focus == choices.FTPOST {
		reaction.PostID = &post.ID
		reaction.Post = post
	} else if focus == choices.FTCOMMENT {
		reaction.CommentID = &comment.ID
		reaction.Comment = comment
	} else {
		reaction.ReplyID = &reply.ID
		reaction.Reply = reply
	}
	db.Save(&reaction)
	return reaction
}

func (obj ReactionManager) Create(db *gorm.DB, user models.User, focus choices.FocusTypeChoice, post *models.Post, comment *models.Comment, reply *models.Reply, rtype choices.ReactionChoice) models.Reaction {
	reaction := models.Reaction{UserObj: user, UserID: user.ID, Rtype: rtype}
	if focus == choices.FTPOST {
		reaction.PostID = &post.ID
		reaction.Post = post
	} else if focus == choices.FTCOMMENT {
		reaction.CommentID = &comment.ID
		reaction.Comment = comment
	} else {
		reaction.ReplyID = &reply.ID
		reaction.Reply = reply
	}
	db.Create(&reaction)
	return reaction
}

func (obj ReactionManager) UpdateOrCreate(db *gorm.DB, user models.User, focus choices.FocusTypeChoice, slug string, rtype choices.ReactionChoice) (*models.Reaction, *models.User, *int, *utils.ErrorResponse) {
	q := db.Scopes(UserAvatarReactionScope)
	var post *models.Post
	var comment *models.Comment
	var reply *models.Reply

	var targetedObjAuthor *models.User
	reaction := models.Reaction{}
	if focus == choices.FTPOST {
		// Get Post Object and Query reactions for the post
		postObj, errCode, errData := PostManager{}.GetBySlug(db, slug, true)
		if errCode != nil {
			return nil, nil, errCode, errData
		}
		post = postObj
		q = q.Where(models.Reaction{PostID: &post.ID})
		targetedObjAuthor = &post.AuthorObj
	} else if focus == choices.FTCOMMENT {
		// Get Comment Object and Query reactions for the comment
		commentObj, errCode, errData := CommentManager{}.GetBySlug(db, slug, true)
		if errCode != nil {
			return nil, nil, errCode, errData
		}
		comment = commentObj
		q = q.Where(models.Reaction{CommentID: &comment.ID})
		targetedObjAuthor = &comment.AuthorObj
	} else {
		// Get Reply Object and Query reactions for the reply
		replyObj, errCode, errData := ReplyManager{}.GetBySlug(db, slug, true)
		if errCode != nil {
			return nil, nil, errCode, errData
		}
		reply = replyObj
		q = q.Where(models.Reaction{ReplyID: &reply.ID})
		targetedObjAuthor = &reply.AuthorObj
	}
	q.Take(&reaction, reaction)
	if reaction.ID == nil {
		// Create reaction
		reaction = obj.Create(db, user, focus, post, comment, reply, rtype)
	} else {
		// Update
		reaction = obj.Update(db, reaction, focus, post, comment, reply, rtype)
	}

	return &reaction, targetedObjAuthor, nil, nil
}

func (obj ReactionManager) GetByID(db *gorm.DB, id *uuid.UUID) (*models.Reaction, *int, *utils.ErrorResponse) {
	reaction := models.Reaction{}
	db.Scopes(UserAvatarReactionScope).Joins("Post").Joins("Comment").Joins("Reply").Take(&reaction, *id)
	if reaction.ID == nil {
		statusCode := 404
		errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "Reaction does not exist")
		return nil, &statusCode, &errData
	}
	return &reaction, nil, nil
}

func (obj ReactionManager) DropData(db *gorm.DB) {
	db.Delete(&models.Reaction{})
}
