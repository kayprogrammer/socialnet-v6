package routes

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/pborman/uuid"
)

var (
	postManager = managers.PostManager{}
	validator = utils.Validator()
)

// @Summary Retrieve Latest Posts
// @Description This endpoint retrieves paginated responses of latest posts
// @Tags Feed
// @Param page query int false "Current Page" default(1)
// @Success 200 {object} schemas.PostsResponseSchema
// @Router /feed/posts [get]
func (endpoint Endpoint) RetrievePosts(c *fiber.Ctx) error {
	db := endpoint.DB
	posts := postManager.All(db)

	// Paginate, Convert type and return Posts
	paginatedData, paginatedPosts, err := PaginateQueryset(posts, c)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	response := schemas.PostsResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Posts fetched"}.Init(),
		Data: schemas.PostsResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
			Items:                       *convertedPosts,
		}.Init(),
	}
	return c.Status(200).JSON(response)
}

// // @Summary Create Post
// // @Description This endpoint creates a new post
// // @Tags Feed
// // @Param post body schemas.PostInputSchema true "Post object"
// // @Success 201 {object} schemas.PostInputResponseSchema
// // @Router /feed/posts [post]
// // @Security BearerAuth
// func (endpoint Endpoint) CreatePost(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	user := c.Locals("user").(*ent.User)

// 	postData := schemas.PostInputSchema{}

// 	// Validate request
// 	if errCode, errData := DecodeJSONBody(c, &postData); errData != nil {
// 		return c.Status(errCode).JSON(errData)
// 	}
// 	if err := validator.Validate(postData); err != nil {
// 		return c.Status(422).JSON(err)
// 	}

// 	post := postManager.Create(db, user, postData)

// 	// Convert type and return Post
// 	convertedPost := utils.ConvertStructData(post, schemas.PostInputResponseDataSchema{}).(*schemas.PostInputResponseDataSchema)
// 	response := schemas.PostInputResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Post created"}.Init(),
// 		Data:           convertedPost.Init(postData.FileType),
// 	}
// 	return c.Status(201).JSON(response)
// }

// // @Summary Retrieve Single Post
// // @Description This endpoint retrieves a single post
// // @Tags Feed
// // @Param slug path string true "Post slug"
// // @Success 200 {object} schemas.PostResponseSchema
// // @Router /feed/posts/{slug} [get]
// func (endpoint Endpoint) RetrievePost(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	slug := c.Params("slug")

// 	// Retrieve, Convert type and return Post
// 	post, errCode, errData := postManager.GetBySlug(db, slug, true)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}
// 	convertedPost := utils.ConvertStructData(post, schemas.PostSchema{}).(*schemas.PostSchema)
// 	response := schemas.PostResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Post Detail fetched"}.Init(),
// 		Data:           convertedPost.Init(),
// 	}
// 	return c.Status(200).JSON(response)
// }

// // @Summary Update Post
// // @Description This endpoint updates a post
// // @Tags Feed
// // @Param slug path string true "Post slug"
// // @Param post body schemas.PostInputSchema true "Post object"
// // @Success 200 {object} schemas.PostInputResponseSchema
// // @Router /feed/posts/{slug} [put]
// // @Security BearerAuth
// func (endpoint Endpoint) UpdatePost(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	user := c.Locals("user").(*ent.User)
// 	slug := c.Params("slug")

// 	postData := schemas.PostInputSchema{}

// 	// Validate request
// 	if errCode, errData := DecodeJSONBody(c, &postData); errData != nil {
// 		return c.Status(errCode).JSON(errData)
// 	}
// 	if err := validator.Validate(postData); err != nil {
// 		return c.Status(422).JSON(err)
// 	}

// 	// Retrieve & Validate Post Existence
// 	post, errCode, errData := postManager.GetBySlug(db, slug, true)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}

// 	// Validate Post ownership
// 	if post.AuthorID != user.ID {
// 		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "This Post isn't yours"))
// 	}

// 	// Update, Convert type and return Post
// 	post = postManager.Update(db, post, postData)
// 	convertedPost := utils.ConvertStructData(post, schemas.PostInputResponseDataSchema{}).(*schemas.PostInputResponseDataSchema)
// 	response := schemas.PostInputResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Post updated"}.Init(),
// 		Data:           convertedPost.Init(postData.FileType),
// 	}
// 	return c.Status(200).JSON(response)
// }

// // @Summary Delete a Post
// // @Description This endpoint deletes a post
// // @Tags Feed
// // @Param slug path string true "Post slug"
// // @Success 200 {object} schemas.ResponseSchema
// // @Router /feed/posts/{slug} [delete]
// // @Security BearerAuth
// func (endpoint Endpoint) DeletePost(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	slug := c.Params("slug")
// 	user := c.Locals("user").(*ent.User)

// 	// Retrieve & Validate Post Existence
// 	post, errCode, errData := postManager.GetBySlug(db, slug)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}

// 	// Validate Post ownership
// 	if post.AuthorID != user.ID {
// 		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "This Post isn't yours"))
// 	}

// 	// Delete and return response
// 	db.Post.DeleteOne(post).Exec(managers.Ctx)
// 	response := schemas.ResponseSchema{Message: "Post Deleted"}.Init()
// 	return c.Status(200).JSON(response)
// }

// var reactionManager = managers.ReactionManager{}

// // @Summary Retrieve Latest Reactions of a Post, Comment, or Reply
// // @Description This endpoint retrieves paginated responses of reactions of post, comment, reply
// // @Tags Feed
// // @Param focus path string true "Specify the usage. Use any of the three: POST, COMMENT, REPLY"
// // @Param slug path string true "Enter the slug of the post or comment or reply"
// // @Param page query int false "Current Page" default(1)
// // @Param reaction_type query string false "Reaction Type. Must be any of these: LIKE, LOVE, HAHA, WOW, SAD, ANGRY"
// // @Success 200 {object} schemas.ReactionsResponseSchema
// // @Router /feed/reactions/{focus}/{slug} [get]
// func (endpoint Endpoint) RetrieveReactions(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	focus := c.Params("focus")
// 	slug := c.Params("slug")

// 	// Validate Focus
// 	err := ValidateReactionFocus(focus)
// 	if err != nil {
// 		return c.Status(404).JSON(err)
// 	}

// 	// Paginate, Convert type and return Posts
// 	reactions, errCode, errData := reactionManager.GetReactionsQueryset(db, c, focus, slug)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}
// 	// Paginate, Convert type and return Reactions
// 	paginatedData, paginatedReactions, err := PaginateQueryset(reactions, c)
// 	if err != nil {
// 		return c.Status(400).JSON(err)
// 	}

// 	convertedReactions := utils.ConvertStructData(paginatedReactions, []schemas.ReactionSchema{}).(*[]schemas.ReactionSchema)
// 	response := schemas.ReactionsResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Reactions fetched"}.Init(),
// 		Data: schemas.ReactionsResponseDataSchema{
// 			PaginatedResponseDataSchema: *paginatedData,
// 			Items:                       *convertedReactions,
// 		}.Init(),
// 	}
// 	return c.Status(200).JSON(response)
// }

// // @Summary Create Reaction
// // @Description This endpoint creates a new reaction.
// // @Tags Feed
// // @Param focus path string true "Specify the usage. Use any of the three: POST, COMMENT, REPLY"
// // @Param slug path string true "Enter the slug of the post or comment or reply"
// // @Param post body schemas.ReactionInputSchema true "Reaction object. rtype should be any of these: LIKE, LOVE, HAHA, WOW, SAD, ANGRY"
// // @Success 201 {object} schemas.ReactionResponseSchema
// // @Router /feed/reactions/{focus}/{slug} [post]
// // @Security BearerAuth
// func (endpoint Endpoint) CreateReaction(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	focus := c.Params("focus")
// 	slug := c.Params("slug")
// 	user := c.Locals("user").(*ent.User)

// 	// Validate Focus
// 	err := ValidateReactionFocus(focus)
// 	if err != nil {
// 		return c.Status(404).JSON(err)
// 	}

// 	reactionData := schemas.ReactionInputSchema{}

// 	// Validate request
// 	if errCode, errData := DecodeJSONBody(c, &reactionData); errData != nil {
// 		return c.Status(errCode).JSON(errData)
// 	}
// 	if err := validator.Validate(reactionData); err != nil {
// 		return c.Status(422).JSON(err)
// 	}

// 	// Update Or Create Reaction
// 	reaction, targetedObjAuthor, errCode, errData := reactionManager.UpdateOrCreate(db, user, focus, slug, reactionData.Rtype)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}

// 	// Convert type and return Reactions
// 	convertedReaction := utils.ConvertStructData(reaction, schemas.ReactionSchema{}).(*schemas.ReactionSchema)
// 	response := schemas.ReactionResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Reaction created"}.Init(),
// 		Data:           convertedReaction.Init(),
// 	}

// 	// Create & Send Notifications
// 	if user.ID != targetedObjAuthor.ID {
// 		notification, created := notificationManager.GetOrCreate(
// 			db, user, "REACTION",
// 			[]uuid.UUID{targetedObjAuthor.ID},
// 			reaction.Edges.Post,
// 			reaction.Edges.Comment,
// 			reaction.Edges.Reply,
// 		)

// 		if created {
// 			SendNotificationInSocket(c, notification, nil, nil)
// 		}
// 	}
// 	return c.Status(201).JSON(response)
// }

// // @Summary Remove Reaction
// // @Description This endpoint deletes a reaction
// // @Tags Feed
// // @Param id path string true "Reaction id (uuid)"
// // @Success 200 {object} schemas.ResponseSchema
// // @Router /feed/reactions/{id} [delete]
// // @Security BearerAuth
// func (endpoint Endpoint) DeleteReaction(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	id := c.Params("id")
// 	// Parse the UUID parameter
// 	reactionID, err := utils.ParseUUID(id)
// 	if err != nil {
// 		return c.Status(400).JSON(err)
// 	}
// 	user := c.Locals("user").(*ent.User)

// 	// Retrieve & Validate Reaction Existence & Ownership
// 	reaction, errCode, errData := reactionManager.GetByID(db, *reactionID)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}

// 	// Validate Reaction ownership
// 	if reaction.UserID != user.ID {
// 		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "This Reaction isn't yours"))
// 	}

// 	// Remove Reaction Notifications
// 	notification := notificationManager.Get(
// 		db, user, "REACTION",
// 		reaction.Edges.Post, reaction.Edges.Comment, reaction.Edges.Reply,
// 	)
// 	if notification != nil {
// 		// Send to websocket and delete notification
// 		SendNotificationInSocket(c, notification, nil, nil, "DELETED")
// 	}

// 	// Delete reaction and return response
// 	db.Reaction.DeleteOne(reaction).Exec(managers.Ctx)
// 	response := schemas.ResponseSchema{Message: "Reaction Deleted"}.Init()
// 	return c.Status(200).JSON(response)
// }

// var commentManager = managers.CommentManager{}

// // @Summary Retrieve Post Comments
// // @Description This endpoint retrieves comments of a particular post
// // @Tags Feed
// // @Param slug path string true "Post Slug"
// // @Param page query int false "Current Page" default(1)
// // @Success 200 {object} schemas.CommentsResponseSchema
// // @Router /feed/posts/{slug}/comments [get]
// func (endpoint Endpoint) RetrieveComments(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	slug := c.Params("slug")

// 	// Get Post
// 	post, errCode, errData := postManager.GetBySlug(db, slug)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}

// 	// Get Comments
// 	comments := commentManager.GetByPostID(db, post.ID)

// 	// Paginate, Convert type and return comments
// 	paginatedData, paginatedComments, err := PaginateQueryset(comments, c)
// 	if err != nil {
// 		return c.Status(400).JSON(err)
// 	}
// 	convertedComments := utils.ConvertStructData(paginatedComments, []schemas.CommentSchema{}).(*[]schemas.CommentSchema)
// 	response := schemas.CommentsResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Comments fetched"}.Init(),
// 		Data: schemas.CommentsResponseDataSchema{
// 			PaginatedResponseDataSchema: *paginatedData,
// 			Items:                       *convertedComments,
// 		}.Init(),
// 	}
// 	return c.Status(200).JSON(response)
// }

// // @Summary Create Comment
// // @Description This endpoint creates a new comment for a particular post
// // @Tags Feed
// // @Param slug path string true "Post Slug"
// // @Param comment body schemas.CommentInputSchema true "Comment object"
// // @Success 201 {object} schemas.CommentResponseSchema
// // @Router /feed/posts/{slug}/comments [post]
// // @Security BearerAuth
// func (endpoint Endpoint) CreateComment(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	slug := c.Params("slug")
// 	user := c.Locals("user").(*ent.User)

// 	// Get Post
// 	post, errCode, errData := postManager.GetBySlug(db, slug)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}

// 	commentData := schemas.CommentInputSchema{}
// 	// Validate request
// 	if errCode, errData := DecodeJSONBody(c, &commentData); errData != nil {
// 		return c.Status(errCode).JSON(errData)
// 	}
// 	if err := validator.Validate(commentData); err != nil {
// 		return c.Status(422).JSON(err)
// 	}

// 	// Create Comment
// 	comment := commentManager.Create(db, user, post.ID, commentData.Text)

// 	// Created & Send Notification
// 	if user.ID != post.AuthorID {
// 		notification := notificationManager.Create(db, user, "COMMENT", []uuid.UUID{post.AuthorID}, nil, comment, nil, nil)
// 		SendNotificationInSocket(c, notification, nil, nil)
// 	}
// 	// Convert type and return comment
// 	convertedComment := utils.ConvertStructData(comment, schemas.CommentSchema{}).(*schemas.CommentSchema)
// 	response := schemas.CommentResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Comment created"}.Init(),
// 		Data:           convertedComment.Init(),
// 	}
// 	return c.Status(201).JSON(response)
// }

// // @Summary Retrieve Comment with replies
// // @Description This endpoint retrieves a comment with replies
// // @Tags Feed
// // @Param slug path string true "Comment Slug"
// // @Param page query int false "Current Page" default(1)
// // @Success 200 {object} schemas.CommentWithRepliesResponseSchema
// // @Router /feed/comments/{slug} [get]
// func (endpoint Endpoint) RetrieveCommentWithReplies(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	slug := c.Params("slug")

// 	// Get Comment
// 	comment, errCode, errData := commentManager.GetBySlug(db, slug, true)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}

// 	// Paginate, Convert type and return replies
// 	convertedComment := utils.ConvertStructData(comment, schemas.CommentSchema{}).(*schemas.CommentSchema)
// 	paginatedData, paginatedReplies, err := PaginateQueryset(convertedComment.Edges.Replies, c)
// 	if err != nil {
// 		return c.Status(400).JSON(err)
// 	}
// 	convertedReplies := utils.ConvertStructData(paginatedReplies, []schemas.ReplySchema{}).(*[]schemas.ReplySchema)
// 	response := schemas.CommentWithRepliesResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Comment with replies fetched"}.Init(),
// 		Data: schemas.CommentWithRepliesSchema{
// 			Comment: convertedComment.Init(),
// 			Replies: schemas.CommentWithRepliesResponseDataSchema{
// 				PaginatedResponseDataSchema: *paginatedData,
// 				Items:                       *convertedReplies,
// 			}.Init(),
// 		},
// 	}
// 	return c.Status(200).JSON(response)
// }

// var replyManager = managers.ReplyManager{}

// // @Summary Create Reply
// // @Description This endpoint creates a reply for a comment
// // @Tags Feed
// // @Param slug path string true "Comment Slug"
// // @Param reply body schemas.CommentInputSchema true "Reply object"
// // @Success 201 {object} schemas.ReplyResponseSchema
// // @Router /feed/comments/{slug} [post]
// // @Security BearerAuth
// func (endpoint Endpoint) CreateReply(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	slug := c.Params("slug")
// 	user := c.Locals("user").(*ent.User)

// 	// Get Comment
// 	comment, errCode, errData := commentManager.GetBySlug(db, slug)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}

// 	replyData := schemas.CommentInputSchema{}
// 	// Validate request
// 	if errCode, errData := DecodeJSONBody(c, &replyData); errData != nil {
// 		return c.Status(errCode).JSON(errData)
// 	}
// 	if err := validator.Validate(replyData); err != nil {
// 		return c.Status(422).JSON(err)
// 	}

// 	// Create reply
// 	reply := replyManager.Create(db, user, comment.ID, replyData.Text)

// 	// Created & Send Notification
// 	if user.ID != comment.AuthorID {
// 		notification := notificationManager.Create(db, user, "REPLY", []uuid.UUID{comment.AuthorID}, nil, nil, reply, nil)
// 		SendNotificationInSocket(c, notification, nil, nil)
// 	}

// 	// Convert type and return reply
// 	convertedReply := utils.ConvertStructData(reply, schemas.ReplySchema{}).(*schemas.ReplySchema)
// 	response := schemas.ReplyResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Reply created"}.Init(),
// 		Data:           convertedReply.Init(),
// 	}
// 	return c.Status(201).JSON(response)
// }

// // @Summary Update Comment
// // @Description This endpoint updates a comment
// // @Tags Feed
// // @Param slug path string true "Comment Slug"
// // @Param comment body schemas.CommentInputSchema true "Comment object"
// // @Success 200 {object} schemas.CommentResponseSchema
// // @Router /feed/comments/{slug} [put]
// // @Security BearerAuth
// func (endpoint Endpoint) UpdateComment(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	slug := c.Params("slug")
// 	user := c.Locals("user").(*ent.User)

// 	// Get Comment
// 	comment, errCode, errData := commentManager.GetBySlug(db, slug, true)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}
// 	if comment.AuthorID != user.ID {
// 		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "Not yours to edit"))
// 	}

// 	commentData := schemas.CommentInputSchema{}
// 	// Validate request
// 	if errCode, errData := DecodeJSONBody(c, &commentData); errData != nil {
// 		return c.Status(errCode).JSON(errData)
// 	}
// 	if err := validator.Validate(commentData); err != nil {
// 		return c.Status(422).JSON(err)
// 	}

// 	// Update Comment
// 	comment = commentManager.Update(comment, user, commentData.Text)

// 	// Convert type and return comment
// 	convertedComment := utils.ConvertStructData(comment, schemas.CommentSchema{}).(*schemas.CommentSchema)
// 	response := schemas.CommentResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Comment updated"}.Init(),
// 		Data:           convertedComment.Init(),
// 	}
// 	return c.Status(200).JSON(response)
// }

// // @Summary Delete Comment
// // @Description This endpoint deletes a comment
// // @Tags Feed
// // @Param slug path string true "Comment Slug"
// // @Success 200 {object} schemas.ResponseSchema
// // @Router /feed/comments/{slug} [delete]
// // @Security BearerAuth
// func (endpoint Endpoint) DeleteComment(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	slug := c.Params("slug")
// 	user := c.Locals("user").(*ent.User)

// 	// Retrieve & Validate Comment Existence & Ownership
// 	comment, errCode, errData := commentManager.GetBySlug(db, slug)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}
// 	if comment.AuthorID != user.ID {
// 		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "Not yours to delete"))
// 	}

// 	// Remove Comment Notifications
// 	notification := notificationManager.Get(
// 		db, user, "COMMENT",
// 		nil, comment, nil,
// 	)
// 	if notification != nil {
// 		// Send to websocket and delete notification & comment
// 		SendNotificationInSocket(c, notification, &comment.Slug, nil, "DELETED")
// 	}

// 	// Delete and return response
// 	response := schemas.ResponseSchema{Message: "Comment Deleted"}.Init()
// 	return c.Status(200).JSON(response)
// }

// // @Summary Retrieve Reply
// // @Description This endpoint retrieves a reply
// // @Tags Feed
// // @Param slug path string true "Reply Slug"
// // @Success 200 {object} schemas.ReplyResponseSchema
// // @Router /feed/replies/{slug} [get]
// func (endpoint Endpoint) RetrieveReply(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	slug := c.Params("slug")

// 	// Get Reply
// 	reply, errCode, errData := replyManager.GetBySlug(db, slug, true)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}

// 	// Convert type and return reply
// 	convertedReply := utils.ConvertStructData(reply, schemas.ReplySchema{}).(*schemas.ReplySchema)
// 	response := schemas.ReplyResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Reply Fetched"}.Init(),
// 		Data:           convertedReply.Init(),
// 	}
// 	return c.Status(200).JSON(response)
// }

// // @Summary Update Reply
// // @Description This endpoint updates a reply
// // @Tags Feed
// // @Param slug path string true "Reply Slug"
// // @Param reply body schemas.CommentInputSchema true "Reply object"
// // @Success 200 {object} schemas.ReplyResponseSchema
// // @Router /feed/replies/{slug} [put]
// // @Security BearerAuth
// func (endpoint Endpoint) UpdateReply(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	slug := c.Params("slug")
// 	user := c.Locals("user").(*ent.User)

// 	// Get Reply
// 	reply, errCode, errData := replyManager.GetBySlug(db, slug, true)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}
// 	if reply.AuthorID != user.ID {
// 		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "Not yours to edit"))
// 	}

// 	replyData := schemas.CommentInputSchema{}
// 	// Validate request
// 	if errCode, errData := DecodeJSONBody(c, &replyData); errData != nil {
// 		return c.Status(errCode).JSON(errData)
// 	}
// 	if err := validator.Validate(replyData); err != nil {
// 		return c.Status(422).JSON(err)
// 	}

// 	// Update Reply
// 	reply = replyManager.Update(reply, user, replyData.Text)

// 	// Convert type and return reply
// 	convertedReply := utils.ConvertStructData(reply, schemas.ReplySchema{}).(*schemas.ReplySchema)
// 	response := schemas.ReplyResponseSchema{
// 		ResponseSchema: schemas.ResponseSchema{Message: "Reply updated"}.Init(),
// 		Data:           convertedReply.Init(),
// 	}
// 	return c.Status(200).JSON(response)
// }

// // @Summary Delete Reply
// // @Description This endpoint deletes a reply
// // @Tags Feed
// // @Param slug path string true "Reply Slug"
// // @Success 200 {object} schemas.ResponseSchema
// // @Router /feed/replies/{slug} [delete]
// // @Security BearerAuth
// func (endpoint Endpoint) DeleteReply(c *fiber.Ctx) error {
// 	db := endpoint.DB
// 	slug := c.Params("slug")
// 	user := c.Locals("user").(*ent.User)

// 	// Retrieve & Validate Reply Existence & Ownership
// 	reply, errCode, errData := replyManager.GetBySlug(db, slug)
// 	if errCode != nil {
// 		return c.Status(*errCode).JSON(errData)
// 	}
// 	if reply.AuthorID != user.ID {
// 		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "Not yours to delete"))
// 	}

// 	// Remove Reply Notifications
// 	notification := notificationManager.Get(
// 		db, user, "REPLY",
// 		nil, nil, reply,
// 	)
// 	if notification != nil {
// 		// Send to websocket and delete notification
// 		SendNotificationInSocket(c, notification, nil, &reply.Slug, "DELETED")
// 	}

// 	// Delete and return response
// 	response := schemas.ResponseSchema{Message: "Reply Deleted"}.Init()
// 	return c.Status(200).JSON(response)
// }
