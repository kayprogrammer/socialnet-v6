package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Endpoint struct {
	DB *gorm.DB
}

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	endpoint := Endpoint{DB: db}

	api := app.Group("/api/v6")

	// HealthCheck Route (1)
	api.Get("/healthcheck", HealthCheck)

	// General Routes (1)
	generalRouter := api.Group("/general")
	generalRouter.Get("/site-detail", endpoint.GetSiteDetails)

	// Auth Routes (8)
	authRouter := api.Group("/auth")
	authRouter.Post("/register", endpoint.Register)
	authRouter.Post("/verify-email", endpoint.VerifyEmail)
	authRouter.Post("/resend-verification-email", endpoint.ResendVerificationEmail)
	authRouter.Post("/send-password-reset-otp", endpoint.SendPasswordResetOtp)
	authRouter.Post("/set-new-password", endpoint.SetNewPassword)
	authRouter.Post("/login", endpoint.Login)
	authRouter.Post("/refresh", endpoint.Refresh)
	authRouter.Get("/logout", endpoint.AuthMiddleware, endpoint.Logout)

	// Profile Routes (12)
	profilesRouter := api.Group("/profiles")
	profilesRouter.Get("/cities", endpoint.RetrieveCities)
	profilesRouter.Get("", endpoint.GuestMiddleware, endpoint.RetrieveUsers)
	profilesRouter.Get("/profile/:username", endpoint.RetrieveUserProfile)
	profilesRouter.Patch("/profile", endpoint.AuthMiddleware, endpoint.UpdateProfile)
	profilesRouter.Post("/profile", endpoint.AuthMiddleware, endpoint.DeleteUser)
	profilesRouter.Get("/friends", endpoint.AuthMiddleware, endpoint.RetrieveFriends)
	profilesRouter.Get("/friends/requests", endpoint.AuthMiddleware, endpoint.RetrieveFriendRequests)
	profilesRouter.Post("/friends/requests", endpoint.AuthMiddleware, endpoint.SendOrDeleteFriendRequest)
	profilesRouter.Put("/friends/requests", endpoint.AuthMiddleware, endpoint.AcceptOrRejectFriendRequest)
	profilesRouter.Get("/notifications", endpoint.AuthMiddleware, endpoint.RetrieveUserNotifications)
	profilesRouter.Post("/notifications", endpoint.AuthMiddleware, endpoint.ReadNotification)

	// Feed Routes (18)
	feedRouter := api.Group("/feed")
	feedRouter.Get("/posts", endpoint.RetrievePosts)
	feedRouter.Post("/posts", endpoint.AuthMiddleware, endpoint.CreatePost)
	feedRouter.Get("/posts/:slug", endpoint.RetrievePost)
	feedRouter.Put("/posts/:slug", endpoint.AuthMiddleware, endpoint.UpdatePost)
	feedRouter.Delete("/posts/:slug", endpoint.AuthMiddleware, endpoint.DeletePost)
	// feedRouter.Get("/reactions/:focus/:slug", endpoint.RetrieveReactions)
	// feedRouter.Post("/reactions/:focus/:slug", endpoint.AuthMiddleware, endpoint.CreateReaction)
	// feedRouter.Delete("/reactions/:id", endpoint.AuthMiddleware, endpoint.DeleteReaction)
	// feedRouter.Get("/posts/:slug/comments", endpoint.RetrieveComments)
	// feedRouter.Post("/posts/:slug/comments", endpoint.AuthMiddleware, endpoint.CreateComment)
	// feedRouter.Get("/comments/:slug", endpoint.RetrieveCommentWithReplies)
	// feedRouter.Post("/comments/:slug", endpoint.AuthMiddleware, endpoint.CreateReply)
	// feedRouter.Put("/comments/:slug", endpoint.AuthMiddleware, endpoint.UpdateComment)
	// feedRouter.Delete("/comments/:slug", endpoint.AuthMiddleware, endpoint.DeleteComment)
	// feedRouter.Get("/replies/:slug", endpoint.RetrieveReply)
	// feedRouter.Put("/replies/:slug", endpoint.AuthMiddleware, endpoint.UpdateReply)
	// feedRouter.Delete("/replies/:slug", endpoint.AuthMiddleware, endpoint.DeleteReply)
}
