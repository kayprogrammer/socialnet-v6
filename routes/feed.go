package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/socialnet-v6/managers"
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/models/choices"
	"github.com/kayprogrammer/socialnet-v6/schemas"
	"github.com/kayprogrammer/socialnet-v6/utils"
)

var postManager = managers.PostManager{}

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
	posts = paginatedPosts.([]models.Post)
	response := schemas.PostsResponseSchema{
		ResponseSchema: SuccessResponse("Posts fetched"),
		Data: schemas.PostsResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
			Items:                       posts,
		}.Init(),
	}
	return c.Status(200).JSON(response)
}

// @Summary Create Post
// @Description This endpoint creates a new post
// @Tags Feed
// @Param post body schemas.PostInputSchema true "Post object"
// @Success 201 {object} schemas.PostInputResponseSchema
// @Router /feed/posts [post]
// @Security BearerAuth
func (endpoint Endpoint) CreatePost(c *fiber.Ctx) error {
	db := endpoint.DB
	user := RequestUser(c)
	data := schemas.PostInputSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	post := postManager.Create(db, *user, data)

	// Convert type and return Post
	response := schemas.PostInputResponseSchema{
		ResponseSchema: SuccessResponse("Post created"),
		Data:           post.InitC(data.FileType),
	}
	return c.Status(201).JSON(response)
}

// @Summary Retrieve Single Post
// @Description This endpoint retrieves a single post
// @Tags Feed
// @Param slug path string true "Post slug"
// @Success 200 {object} schemas.PostResponseSchema
// @Router /feed/posts/{slug} [get]
func (endpoint Endpoint) RetrievePost(c *fiber.Ctx) error {
	db := endpoint.DB
	slug := c.Params("slug")

	// Retrieve, Convert type and return Post
	post, errCode, errData := postManager.GetBySlug(db, slug, true)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}
	response := schemas.PostResponseSchema{
		ResponseSchema: SuccessResponse("Post Detail fetched"),
		Data:           *post,
	}
	return c.Status(200).JSON(response)
}

// @Summary Update Post
// @Description This endpoint updates a post
// @Tags Feed
// @Param slug path string true "Post slug"
// @Param post body schemas.PostInputSchema true "Post object"
// @Success 200 {object} schemas.PostInputResponseSchema
// @Router /feed/posts/{slug} [put]
// @Security BearerAuth
func (endpoint Endpoint) UpdatePost(c *fiber.Ctx) error {
	db := endpoint.DB
	user := RequestUser(c)
	slug := c.Params("slug")

	data := schemas.PostInputSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Retrieve & Validate Post Existence
	post, errCode, errData := postManager.GetBySlug(db, slug, true)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Validate Post ownership
	if post.AuthorID.String() != user.ID.String() {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "This Post isn't yours"))
	}

	// Update, Convert type and return Post
	post = postManager.Update(db, post, data)
	response := schemas.PostInputResponseSchema{
		ResponseSchema: SuccessResponse("Post updated"),
		Data:           post.InitC(data.FileType),
	}
	return c.Status(200).JSON(response)
}

// @Summary Delete a Post
// @Description This endpoint deletes a post
// @Tags Feed
// @Param slug path string true "Post slug"
// @Success 200 {object} schemas.ResponseSchema
// @Router /feed/posts/{slug} [delete]
// @Security BearerAuth
func (endpoint Endpoint) DeletePost(c *fiber.Ctx) error {
	db := endpoint.DB
	slug := c.Params("slug")
	user := RequestUser(c)

	// Retrieve & Validate Post Existence
	post, errCode, errData := postManager.GetBySlug(db, slug)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Validate Post ownership
	if post.AuthorID.String() != user.ID.String() {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "This Post isn't yours"))
	}

	// Delete and return response
	db.Delete(&post)
	return c.Status(200).JSON(SuccessResponse("Post Deleted"))
}

var reactionManager = managers.ReactionManager{}

// @Summary Retrieve Latest Reactions of a Post, Comment, or Reply
// @Description This endpoint retrieves paginated responses of reactions of post, comment, reply
// @Tags Feed
// @Param focus path string true "Specify the usage. Use any of the three: POST, COMMENT, REPLY"
// @Param slug path string true "Enter the slug of the post or comment or reply"
// @Param page query int false "Current Page" default(1)
// @Param reaction_type query string false "Reaction Type. Must be any of these: LIKE, LOVE, HAHA, WOW, SAD, ANGRY"
// @Success 200 {object} schemas.ReactionsResponseSchema
// @Router /feed/reactions/{focus}/{slug} [get]
func (endpoint Endpoint) RetrieveReactions(c *fiber.Ctx) error {
	db := endpoint.DB
	focus := c.Params("focus")
	slug := c.Params("slug")

	// Validate Focus
	err := ValidateReactionFocus(focus)
	if err != nil {
		return c.Status(404).JSON(err)
	}

	// Paginate, Convert type and return Posts
	reactions, errCode, errData := reactionManager.GetReactionsQueryset(db, c, focus, slug)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}
	// Paginate, Convert type and return Reactions
	paginatedData, paginatedReactions, err := PaginateQueryset(reactions, c)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	reactions = paginatedReactions.([]models.Reaction)
	response := schemas.ReactionsResponseSchema{
		ResponseSchema: SuccessResponse("Reactions fetched"),
		Data: schemas.ReactionsResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
			Items:                       reactions,
		},
	}
	return c.Status(200).JSON(response)
}

// @Summary Create Reaction
// @Description This endpoint creates a new reaction.
// @Tags Feed
// @Param focus path string true "Specify the usage. Use any of the three: POST, COMMENT, REPLY"
// @Param slug path string true "Enter the slug of the post or comment or reply"
// @Param post body schemas.ReactionInputSchema true "Reaction object. rtype should be any of these: LIKE, LOVE, HAHA, WOW, SAD, ANGRY"
// @Success 201 {object} schemas.ReactionResponseSchema
// @Router /feed/reactions/{focus}/{slug} [post]
// @Security BearerAuth
func (endpoint Endpoint) CreateReaction(c *fiber.Ctx) error {
	db := endpoint.DB
	focus := c.Params("focus")
	slug := c.Params("slug")
	user := RequestUser(c)

	// Validate Focus
	err := ValidateReactionFocus(focus)
	if err != nil {
		return c.Status(404).JSON(err)
	}

	data := schemas.ReactionInputSchema{}

	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Update Or Create Reaction
	reaction, targetedObjAuthor, errCode, errData := reactionManager.UpdateOrCreate(db, *user, focus, slug, data.Rtype)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Convert type and return Reactions
	response := schemas.ReactionResponseSchema{
		ResponseSchema: SuccessResponse("Reaction created"),
		Data:           *reaction,
	}

	// Create & Send Notifications
	if user.ID.String() != targetedObjAuthor.ID.String() {
		notification, created := notificationManager.GetOrCreate(
			db, user, choices.NREACTION,
			[]models.User{*targetedObjAuthor},
			reaction.Post,
			reaction.Comment,
			reaction.Reply,
		)
		log.Println(created)
		if created {
			SendNotificationInSocket(c, notification, nil, nil)
		}
	}
	return c.Status(201).JSON(response)
}

// @Summary Remove Reaction
// @Description This endpoint deletes a reaction
// @Tags Feed
// @Param id path string true "Reaction id (uuid)"
// @Success 200 {object} schemas.ResponseSchema
// @Router /feed/reactions/{id} [delete]
// @Security BearerAuth
func (endpoint Endpoint) DeleteReaction(c *fiber.Ctx) error {
	db := endpoint.DB
	id := c.Params("id")
	// Parse the UUID parameter
	reactionID, err := utils.ParseUUID(id)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	user := RequestUser(c)

	// Retrieve & Validate Reaction Existence & Ownership
	reaction, errCode, errData := reactionManager.GetByID(db, reactionID)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Validate Reaction ownership
	if reaction.UserID.String() != user.ID.String() {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "This Reaction isn't yours"))
	}

	// Remove Reaction Notifications
	notification := notificationManager.Get(
		db, user, choices.NREACTION,
		reaction.Post, reaction.Comment, reaction.Reply,
	)
	if notification != nil {
		// Send to websocket and delete notification
		SendNotificationInSocket(c, *notification, nil, nil, "DELETED")
	}

	// Delete reaction and return response
	db.Delete(&reaction)
	return c.Status(200).JSON(SuccessResponse("Reaction Deleted"))
}

var commentManager = managers.CommentManager{}

// @Summary Retrieve Post Comments
// @Description This endpoint retrieves comments of a particular post
// @Tags Feed
// @Param slug path string true "Post Slug"
// @Param page query int false "Current Page" default(1)
// @Success 200 {object} schemas.CommentsResponseSchema
// @Router /feed/posts/{slug}/comments [get]
func (endpoint Endpoint) RetrieveComments(c *fiber.Ctx) error {
	db := endpoint.DB
	slug := c.Params("slug")

	// Get Post
	post, errCode, errData := postManager.GetBySlug(db, slug)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Get Comments
	comments := commentManager.GetByPostID(db, post.ID)

	// Paginate, Convert type and return comments
	paginatedData, paginatedComments, err := PaginateQueryset(comments, c)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	comments = paginatedComments.([]models.Comment)
	response := schemas.CommentsResponseSchema{
		ResponseSchema: SuccessResponse("Comments fetched"),
		Data: schemas.CommentsResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
			Items:                       comments,
		}.Init(),
	}
	return c.Status(200).JSON(response)
}

// @Summary Create Comment
// @Description This endpoint creates a new comment for a particular post
// @Tags Feed
// @Param slug path string true "Post Slug"
// @Param comment body schemas.CommentInputSchema true "Comment object"
// @Success 201 {object} schemas.CommentResponseSchema
// @Router /feed/posts/{slug}/comments [post]
// @Security BearerAuth
func (endpoint Endpoint) CreateComment(c *fiber.Ctx) error {
	db := endpoint.DB
	slug := c.Params("slug")
	user := RequestUser(c)

	// Get Post
	post, errCode, errData := postManager.GetBySlug(db, slug)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}

	data := schemas.CommentInputSchema{}
	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Create Comment
	comment := commentManager.Create(db, *user, *post, data.Text)

	// Created & Send Notification
	if user.ID.String() != post.AuthorID.String() {
		notification := notificationManager.Create(db, user, choices.NCOMMENT, []models.User{post.AuthorObj}, nil, &comment, nil, nil)
		SendNotificationInSocket(c, notification, nil, nil)
	}

	response := schemas.CommentResponseSchema{
		ResponseSchema: SuccessResponse("Comment created"),
		Data:           comment.Init(),
	}
	return c.Status(201).JSON(response)
}

// @Summary Retrieve Comment with replies
// @Description This endpoint retrieves a comment with replies
// @Tags Feed
// @Param slug path string true "Comment Slug"
// @Param page query int false "Current Page" default(1)
// @Success 200 {object} schemas.CommentWithRepliesResponseSchema
// @Router /feed/comments/{slug} [get]
func (endpoint Endpoint) RetrieveCommentWithReplies(c *fiber.Ctx) error {
	db := endpoint.DB
	slug := c.Params("slug")

	// Get Comment
	comment, errCode, errData := commentManager.GetBySlug(db, slug, true)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Paginate, Convert type and return replies
	paginatedData, paginatedReplies, err := PaginateQueryset(comment.Replies, c)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	replies := paginatedReplies.([]models.Reply)
	response := schemas.CommentWithRepliesResponseSchema{
		ResponseSchema: SuccessResponse("Comment with replies fetched"),
		Data: schemas.CommentWithRepliesSchema{
			Comment: comment.Init(),
			Replies: schemas.CommentWithRepliesResponseDataSchema{
				PaginatedResponseDataSchema: *paginatedData,
				Items:                       replies,
			}.Init(),
		},
	}
	return c.Status(200).JSON(response)
}

var replyManager = managers.ReplyManager{}

// @Summary Create Reply
// @Description This endpoint creates a reply for a comment
// @Tags Feed
// @Param slug path string true "Comment Slug"
// @Param reply body schemas.CommentInputSchema true "Reply object"
// @Success 201 {object} schemas.ReplyResponseSchema
// @Router /feed/comments/{slug} [post]
// @Security BearerAuth
func (endpoint Endpoint) CreateReply(c *fiber.Ctx) error {
	db := endpoint.DB
	slug := c.Params("slug")
	user := RequestUser(c)

	// Get Comment
	comment, errCode, errData := commentManager.GetBySlug(db, slug)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}

	data := schemas.CommentInputSchema{}
	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Create reply
	reply := replyManager.Create(db, *user, *comment, data.Text)

	// Created & Send Notification
	if user.ID.String() != comment.AuthorID.String() {
		notification := notificationManager.Create(db, user, choices.NREPLY, []models.User{comment.AuthorObj}, nil, nil, &reply, nil)
		SendNotificationInSocket(c, notification, nil, nil)
	}

	// Convert type and return reply
	response := schemas.ReplyResponseSchema{
		ResponseSchema: SuccessResponse("Reply created"),
		Data:           reply.Init(),
	}
	return c.Status(201).JSON(response)
}

// @Summary Update Comment
// @Description This endpoint updates a comment
// @Tags Feed
// @Param slug path string true "Comment Slug"
// @Param comment body schemas.CommentInputSchema true "Comment object"
// @Success 200 {object} schemas.CommentResponseSchema
// @Router /feed/comments/{slug} [put]
// @Security BearerAuth
func (endpoint Endpoint) UpdateComment(c *fiber.Ctx) error {
	db := endpoint.DB
	slug := c.Params("slug")
	user := RequestUser(c)

	// Get Comment
	comment, errCode, errData := commentManager.GetBySlug(db, slug, true)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}
	if comment.AuthorID.String() != user.ID.String() {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "Not yours to edit"))
	}

	data := schemas.CommentInputSchema{}
	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Update Comment
	updatedComment := commentManager.Update(db, *comment, user, data.Text)

	// Convert type and return comment
	response := schemas.CommentResponseSchema{
		ResponseSchema: SuccessResponse("Comment updated"),
		Data:           updatedComment.Init(),
	}
	return c.Status(200).JSON(response)
}

// @Summary Delete Comment
// @Description This endpoint deletes a comment
// @Tags Feed
// @Param slug path string true "Comment Slug"
// @Success 200 {object} schemas.ResponseSchema
// @Router /feed/comments/{slug} [delete]
// @Security BearerAuth
func (endpoint Endpoint) DeleteComment(c *fiber.Ctx) error {
	db := endpoint.DB
	slug := c.Params("slug")
	user := RequestUser(c)

	// Retrieve & Validate Comment Existence & Ownership
	comment, errCode, errData := commentManager.GetBySlug(db, slug)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}
	if comment.AuthorID.String() != user.ID.String() {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "Not yours to delete"))
	}

	// Remove Comment Notifications
	notification := notificationManager.Get(
		db, user, choices.NCOMMENT,
		nil, comment, nil,
	)
	if notification != nil {
		// Send to websocket and delete notification & comment
		SendNotificationInSocket(c, *notification, &comment.Slug, nil, "DELETED")
	}

	// Return response
	return c.Status(200).JSON(SuccessResponse("Comment Deleted"))
}

// @Summary Retrieve Reply
// @Description This endpoint retrieves a reply
// @Tags Feed
// @Param slug path string true "Reply Slug"
// @Success 200 {object} schemas.ReplyResponseSchema
// @Router /feed/replies/{slug} [get]
func (endpoint Endpoint) RetrieveReply(c *fiber.Ctx) error {
	db := endpoint.DB
	slug := c.Params("slug")

	// Get Reply
	reply, errCode, errData := replyManager.GetBySlug(db, slug, true)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Convert type and return reply
	response := schemas.ReplyResponseSchema{
		ResponseSchema: SuccessResponse("Reply Fetched"),
		Data:           reply.Init(),
	}
	return c.Status(200).JSON(response)
}

// @Summary Update Reply
// @Description This endpoint updates a reply
// @Tags Feed
// @Param slug path string true "Reply Slug"
// @Param reply body schemas.CommentInputSchema true "Reply object"
// @Success 200 {object} schemas.ReplyResponseSchema
// @Router /feed/replies/{slug} [put]
// @Security BearerAuth
func (endpoint Endpoint) UpdateReply(c *fiber.Ctx) error {
	db := endpoint.DB
	slug := c.Params("slug")
	user := RequestUser(c)

	// Get Reply
	reply, errCode, errData := replyManager.GetBySlug(db, slug, true)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}
	if reply.AuthorID.String() != user.ID.String() {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "Not yours to edit"))
	}

	data := schemas.CommentInputSchema{}
	// Validate request
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	// Update Reply
	updatedReply := replyManager.Update(db, *reply, user, data.Text)

	// Convert type and return reply
	response := schemas.ReplyResponseSchema{
		ResponseSchema: SuccessResponse("Reply updated"),
		Data:           updatedReply.Init(),
	}
	return c.Status(200).JSON(response)
}

// @Summary Delete Reply
// @Description This endpoint deletes a reply
// @Tags Feed
// @Param slug path string true "Reply Slug"
// @Success 200 {object} schemas.ResponseSchema
// @Router /feed/replies/{slug} [delete]
// @Security BearerAuth
func (endpoint Endpoint) DeleteReply(c *fiber.Ctx) error {
	db := endpoint.DB
	slug := c.Params("slug")
	user := RequestUser(c)

	// Retrieve & Validate Reply Existence & Ownership
	reply, errCode, errData := replyManager.GetBySlug(db, slug)
	if errCode != nil {
		return c.Status(*errCode).JSON(errData)
	}
	if reply.AuthorID.String() != user.ID.String() {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_OWNER, "Not yours to delete"))
	}

	// Remove Reply Notifications
	notification := notificationManager.Get(
		db, user, choices.NREPLY,
		nil, nil, reply,
	)
	if notification != nil {
		// Send to websocket and delete notification
		SendNotificationInSocket(c, *notification, nil, &reply.Slug, "DELETED")
	}

	// Return response
	return c.Status(200).JSON(SuccessResponse("Reply Deleted"))
}
