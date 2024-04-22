package tests

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/socialnet-v6/database"
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/routes"
	"github.com/kayprogrammer/socialnet-v6/schemas"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)


func register(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Register User", func(t *testing.T) {
		url := fmt.Sprintf("%s/register", baseUrl)
		validEmail := "testregisteruser@email.com"
		userData := schemas.RegisterUser{
			FirstName:      "TestRegister",
			LastName:       "User",
			Email:          validEmail,
			Password:       "testregisteruserpassword",
			TermsAgreement: true,
		}
		res := ProcessTestBody(t, app, url, "POST", userData)

		// Assert Status code
		assert.Equal(t, 201, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Registration successful", body["message"])
		expectedData := make(map[string]interface{})
		expectedData["email"] = validEmail
		assert.Equal(t, expectedData, body["data"].(map[string]interface{}))

		// Verify that a user with the same email cannot be registered again
		res = ProcessTestBody(t, app, url, "POST", userData)
		assert.Equal(t, 422, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_INVALID_ENTRY, body["code"])
		assert.Equal(t, "Invalid Entry", body["message"])
		expectedData = make(map[string]interface{})
		expectedData["email"] = "Email already taken!"
		assert.Equal(t, expectedData, body["data"].(map[string]interface{}))
	})
}

func verifyEmail(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Verify Email", func(t *testing.T) {
		user := CreateTestUser(db)
		otp := uint32(1111)

		url := fmt.Sprintf("%s/verify-email", baseUrl)
		emailOtpData := schemas.VerifyEmailRequestSchema{
			EmailRequestSchema: schemas.EmailRequestSchema{Email: user.Email},
			Otp:   otp,
		}

		res := ProcessTestBody(t, app, url, "POST", emailOtpData)

		// Verify that the email verification fails with an invalid otp
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_INCORRECT_OTP, body["code"])
		assert.Equal(t, "Incorrect Otp", body["message"])

		// Verify that the email verification succeeds with a valid otp
		realOtp := models.Otp{UserId: user.ID}
		db.Take(&realOtp, realOtp)
		db.Save(&realOtp) // Create or save
		emailOtpData.Otp = realOtp.Code
		res = ProcessTestBody(t, app, url, "POST", emailOtpData)
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Account verification successful", body["message"])
	})
}

func resendVerificationEmail(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Resend Verification Email", func(t *testing.T) {
		// Drop User Data since the previous test uses it...
		DropAndCreateSingleTable(db, models.User{})
		user := CreateTestUser(db)

		url := fmt.Sprintf("%s/resend-verification-email", baseUrl)
		emailData := schemas.EmailRequestSchema{
			Email: user.Email,
		}

		res := ProcessTestBody(t, app, url, "POST", emailData)

		// Verify that an unverified user can get a new email
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Verification email sent", body["message"])

		// Verify that a verified user cannot get a new email
		user.IsEmailVerified = true
		db.Save(&user)
		res = ProcessTestBody(t, app, url, "POST", emailData)

		assert.Equal(t, 200, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Email already verified", body["message"])

		// Verify that an error is raised when attempting to resend the verification email for a user that doesn't exist
		emailData.Email = "invalid@example.com"
		res = ProcessTestBody(t, app, url, "POST", emailData)

		assert.Equal(t, 404, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Incorrect Email", body["message"])
	})
}

func sendPasswordResetOtp(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Send Password Reset Otp", func(t *testing.T) {

		user := CreateTestVerifiedUser(db)

		url := fmt.Sprintf("%s/send-password-reset-otp", baseUrl)
		emailData := schemas.EmailRequestSchema{
			Email: user.Email,
		}

		res := ProcessTestBody(t, app, url, "POST", emailData)

		// Verify that an unverified user can get a new email
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Password otp sent", body["message"])

		// Verify that an error is raised when attempting to send password reset email for a user that doesn't exist
		emailData.Email = "invalid@example.com"
		res = ProcessTestBody(t, app, url, "POST", emailData)

		assert.Equal(t, 404, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_INCORRECT_EMAIL, body["code"])
		assert.Equal(t, "Incorrect Email", body["message"])
	})
}

func setNewPassword(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	// Drop User data since the previous test uses the verified_user it...
	DropAndCreateSingleTable(db, models.User{})

	t.Run("Set New Password", func(t *testing.T) {
		user := CreateTestVerifiedUser(db)

		url := fmt.Sprintf("%s/set-new-password", baseUrl)
		passwordResetData := schemas.SetNewPasswordSchema{
			VerifyEmailRequestSchema: schemas.VerifyEmailRequestSchema{
				EmailRequestSchema: schemas.EmailRequestSchema{Email: "invalid@example.com"}, // Invalid otp
				Otp:   11111,                 // Invalid otp
			},
			Password: "newpassword",
		}

		res := ProcessTestBody(t, app, url, "POST", passwordResetData)

		// Verify that the request fails with incorrect email
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_INCORRECT_EMAIL, body["code"])
		assert.Equal(t, "Incorrect Email", body["message"])

		// Verify that the request fails with incorrect otp
		passwordResetData.Email = user.Email
		res = ProcessTestBody(t, app, url, "POST", passwordResetData)
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_INCORRECT_OTP, body["code"])
		assert.Equal(t, "Incorrect Otp", body["message"])

		// Verify that password reset succeeds
		realOtp := models.Otp{UserId: user.ID}
		db.Take(&realOtp, realOtp)
		db.Save(&realOtp) // Create or save
		passwordResetData.Otp = realOtp.Code
		res = ProcessTestBody(t, app, url, "POST", passwordResetData)

		// Assert response
		assert.Equal(t, 200, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Password reset successful", body["message"])
	})
}

func login(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Login", func(t *testing.T) {
		user := CreateTestUser(db)

		url := fmt.Sprintf("%s/login", baseUrl)
		loginData := schemas.LoginSchema{
			Email:    "invalid@example.com", // Invalid email
			Password: "invalidpassword",
		}

		res := ProcessTestBody(t, app, url, "POST", loginData)

		// # Test for invalid credentials
		// Assert Status code
		assert.Equal(t, 401, res.StatusCode)
		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_INVALID_CREDENTIALS, body["code"])
		assert.Equal(t, "Invalid Credentials", body["message"])

		// Test for unverified credentials (email)
		loginData.Email = user.Email
		loginData.Password = "testpassword"
		res = ProcessTestBody(t, app, url, "POST", loginData)
		// Assert Status code
		assert.Equal(t, 401, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_UNVERIFIED_USER, body["code"])
		assert.Equal(t, "Verify your email first", body["message"])

		// Test for valid credentials and verified email address
		user.IsEmailVerified = true
		db.Save(&user)
		res = ProcessTestBody(t, app, url, "POST", loginData)
		// Assert response
		assert.Equal(t, 201, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Login successful", body["message"])
		db.Take(&user, user.ID) // Get updated user 
		expectedData := map[string]string{
			"access":  *user.Access,
			"refresh": *user.Refresh,
		}
		data, _ := json.Marshal(body["data"])
		expectedDataJson, _ := json.Marshal(expectedData)
		assert.Equal(t, expectedDataJson, data)
	})
}

func refresh(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	// Drop User Data since the previous test uses the verified_user it...
	DropAndCreateSingleTable(db, models.User{})

	t.Run("Refresh", func(t *testing.T) {
		user := CreateTestVerifiedUser(db)

		url := fmt.Sprintf("%s/refresh", baseUrl)
		refreshTokenData := schemas.RefreshTokenSchema{
			Refresh: "invalid_token",
		}

		// Test for invalid refresh token (invalid or expired)
		res := ProcessTestBody(t, app, url, "POST", refreshTokenData)
		// Assert Status code
		assert.Equal(t, 401, res.StatusCode)
		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_INVALID_TOKEN, body["code"])
		assert.Equal(t, "Refresh token is invalid or expired", body["message"])

		// Test for valid refresh token
		access := "whatever"
		refresh := routes.GenerateRefreshToken()
		user.Access = &access
		user.Refresh = &refresh
		db.Save(&user)
		refreshTokenData.Refresh = *user.Refresh
		res = ProcessTestBody(t, app, url, "POST", refreshTokenData)
		// Assert response
		assert.Equal(t, 201, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Tokens refresh successful", body["message"])
		db.Take(&user, user.ID) // Get updated user 
		expectedData := map[string]string{
			"access":  *user.Access,
			"refresh": *user.Refresh,
		}
		data, _ := json.Marshal(body["data"])
		expectedDataJson, _ := json.Marshal(expectedData)
		assert.Equal(t, expectedDataJson, data)
	})
}

func logout(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	// Drop and Create User Table since the previous test uses the verified_user it...
	DropAndCreateSingleTable(db, models.User{})
	t.Run("Logout", func(t *testing.T) {
		url := fmt.Sprintf("%s/logout", baseUrl)
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		res, _ := app.Test(req)

		// Ensures an unauthorized user cannot log out
		// Assert Status code
		assert.Equal(t, 401, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_INVALID_TOKEN, body["code"])
		assert.Equal(t, "Auth Token is Invalid or Expired!", body["message"])

		// Ensures an authorized user can log out
		req = httptest.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AccessToken(db)))
		res, _ = app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Logout successful", body["message"])
	})
}

func TestAuth(t *testing.T) {
	os.Setenv("ENVIRONMENT", "TESTING")
	app := fiber.New()
	db := Setup(t, app)
	BASEURL := "/api/v6/auth"

	// Run Auth Endpoint Tests
	register(t, app, db, BASEURL)
	verifyEmail(t, app, db, BASEURL)
	resendVerificationEmail(t, app, db, BASEURL)
	sendPasswordResetOtp(t, app, db, BASEURL)
	setNewPassword(t, app, db, BASEURL)
	login(t, app, db, BASEURL)
	logout(t, app, db, BASEURL)
	refresh(t, app, db, BASEURL)

	// Drop Tables and Close Connectiom
	database.DropTables(db)
	CloseTestDatabase(db)
}
