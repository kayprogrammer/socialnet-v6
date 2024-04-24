package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/socialnet-v6/database"
	"github.com/kayprogrammer/socialnet-v6/models/choices"
	"github.com/kayprogrammer/socialnet-v6/schemas"
	"github.com/kayprogrammer/socialnet-v6/utils"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func getChats(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	CreateChat(db)
	t.Run("Retrieve Chats", func(t *testing.T) {
		url := baseUrl
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AccessToken(db)))
		res, _ := app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Chats fetched", body["message"])
		data, _ := json.Marshal(body["data"])
		assert.Equal(t, true, (len(data) > 0))
	})
}

func sendMessage(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	chat := CreateChat(db)
	sender := chat.OwnerObj
	token := AccessToken(db)
	t.Run("Send Message", func(t *testing.T) {
		url := baseUrl
		invalidUUID := uuid.NewUUID()
		text := "JESUS is KING"
		messageData := schemas.MessageCreateSchema{ChatID: &invalidUUID, Text: &text}
		// Test for valid response for invalid chat id
		res := ProcessTestBody(t, app, url, "POST", messageData, token)
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)
		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_NON_EXISTENT, body["code"])
		assert.Equal(t, "User has no chat with that ID", body["message"])

		// Test for valid response for valid entry
		messageData.ChatID = &chat.ID
		res = ProcessTestBody(t, app, url, "POST", messageData, token)
		// Assert Status code
		assert.Equal(t, 201, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		data, _ := json.Marshal(body)
		dataMap := body["data"].(map[string]interface{})
		expectedData := map[string]interface{}{
			"status":  "success",
			"message": "Message sent",
			"data": map[string]interface{}{
				"id":         body["data"].(map[string]interface{})["id"],
				"chat_id":    chat.ID,
				"sender":     GetUserMap(sender),
				"text":       messageData.Text,
				"file": nil,
				"created_at": dataMap["created_at"],
				"updated_at": dataMap["updated_at"],
			},
		}
		expectedDataJson, _ := json.Marshal(expectedData)
		assert.JSONEq(t, string(expectedDataJson), string(data))
	})
}

func getChatMessages(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	message := CreateMessage(db)
	chat := message.ChatObj
	owner := chat.OwnerObj
	token := AccessToken(db)
	t.Run("Retrieve Chat Messages", func(t *testing.T) {
		invalidChatID := uuid.New()
		url := fmt.Sprintf("%s/%s", baseUrl, invalidChatID)
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		res, _ := app.Test(req)

		// Verify the request fails with invalid chat ID
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, 404, res.StatusCode)
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_NON_EXISTENT, body["code"])
		assert.Equal(t, "User has no chat with that ID", body["message"])

		// Verify the request succeeds with valid chat ID
		url = fmt.Sprintf("%s/%s", baseUrl, chat.ID)
		req = httptest.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		res, _ = app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		dataMap := body["data"].(map[string]interface{})
		chatMap := dataMap["chat"].(map[string]interface{})
		messagesMap := dataMap["messages"].(map[string]interface{})
		messagesItems := messagesMap["items"].([]interface{})
		messageItemMap := messagesItems[0].(map[string]interface{})
		data, _ := json.Marshal(body)
		ownerData := GetUserMap(owner)
		recipientUser := chat.UserObjs[0]
		log.Println("Yaaaaaaaa: ", recipientUser)

		expectedData := map[string]interface{}{
			"status":  "success",
			"message": "Messages fetched",
			"data": map[string]interface{}{
				"chat": map[string]interface{}{
					"id":          chat.ID,
					"name":        chat.Name,
					"owner":       ownerData,
					"ctype":       chat.Ctype,
					"description": chat.Description,
					"image":       nil,
					"latest_message": map[string]interface{}{
						"sender": ownerData,
						"text":   message.Text,
						"file":   nil,
					},
					"created_at": chatMap["created_at"],
					"updated_at": chatMap["updated_at"],
				},
				"messages": map[string]interface{}{
					"per_page":     400,
					"current_page": 1,
					"last_page":    1,
					"items": []map[string]interface{}{
						{
							"id":         message.ID,
							"chat_id":    chat.ID,
							"sender":     ownerData,
							"text":       message.Text,
							"file":       nil,
							"created_at": messageItemMap["created_at"],
							"updated_at": messageItemMap["updated_at"],
						},
					},
				},
				"users": []map[string]interface{}{
					GetUserMap(recipientUser),
				},
			},
		}
		expectedDataJson, _ := json.Marshal(expectedData)
		assert.JSONEq(t, string(expectedDataJson), string(data))
	})
}

func updateGroupChat(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	chat := CreateGroupChat(db)
	user := chat.UserObjs[0]
	token := AccessToken(db)
	t.Run("Update Group Chat", func(t *testing.T) {
		url := fmt.Sprintf("%s/%s", baseUrl, uuid.New())
		name := "Updated Group chat name"
		desc := "Updated group chat description"
		chatData := schemas.GroupChatInputSchema{Name: &name, Description: &desc}

		// Test for valid response for invalid chat id
		res := ProcessTestBody(t, app, url, "PATCH", chatData, token)
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)
		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_NON_EXISTENT, body["code"])
		assert.Equal(t, "User owns no group chat with that ID", body["message"])

		// Test for valid response for valid entry
		url = fmt.Sprintf("%s/%s", baseUrl, chat.ID)
		res = ProcessTestBody(t, app, url, "PATCH", chatData, token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		data, _ := json.Marshal(body)
		dataMap := body["data"].(map[string]interface{})

		expectedData := map[string]interface{}{
			"status":  "success",
			"message": "Chat updated",
			"data": map[string]interface{}{
				"id":          chat.ID,
				"owner": GetUserMap(chat.OwnerObj),
				"name":        chatData.Name,
				"description": chatData.Description,
				"ctype": choices.CGROUP,
				"users": []map[string]interface{}{
					GetUserMap(user),
				},
				"image": nil,
				"latest_message": nil,
				"created_at": dataMap["created_at"],
				"updated_at": dataMap["updated_at"],
			},
		}
		expectedDataJson, _ := json.Marshal(expectedData)
		assert.JSONEq(t, string(expectedDataJson), string(data))

		// You can test for other error responses yourself

	})
}

func deleteGroupChat(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	chat := CreateGroupChat(db)
	token := AccessToken(db)
	t.Run("Delete Group Chat", func(t *testing.T) {
		url := fmt.Sprintf("%s/%s", baseUrl, uuid.New())
		// Test for valid response for invalid chat id
		req := httptest.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		res, _ := app.Test(req)

		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)
		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_NON_EXISTENT, body["code"])
		assert.Equal(t, "User owns no group chat with that ID", body["message"])

		// Test for valid response for valid entry
		url = fmt.Sprintf("%s/%s", baseUrl, chat.ID)
		req = httptest.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		res, _ = app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		data, _ := json.Marshal(body)
		expectedData := map[string]interface{}{
			"status":  "success",
			"message": "Group Chat Deleted",
		}
		expectedDataJson, _ := json.Marshal(expectedData)
		assert.JSONEq(t, string(expectedDataJson), string(data))
		// You can test for other error responses yourself

	})
}

func updateMessage(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	message := CreateMessage(db)
	sender := message.SenderObj
	token := AccessToken(db)
	t.Run("Update A Message", func(t *testing.T) {
		url := fmt.Sprintf("%s/messages/%s", baseUrl, uuid.New())
		text := "Jesus is Lord"
		messageData := schemas.MessageUpdateSchema{Text: &text}

		// Test for valid response for invalid message id
		res := ProcessTestBody(t, app, url, "PUT", messageData, token)
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)
		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_NON_EXISTENT, body["code"])
		assert.Equal(t, "User has no message with that ID", body["message"])

		// Test for valid response for valid entry
		url = fmt.Sprintf("%s/messages/%s", baseUrl, message.ID)
		res = ProcessTestBody(t, app, url, "PUT", messageData, token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		data, _ := json.Marshal(body)
		dataMap := body["data"].(map[string]interface{})
		expectedData := map[string]interface{}{
			"status":  "success",
			"message": "Message updated",
			"data": map[string]interface{}{
				"id":         message.ID,
				"chat_id":    message.ChatID,
				"sender":     GetUserMap(sender),
				"text":       messageData.Text,
				"file": nil,
				"created_at": dataMap["created_at"],
				"updated_at": dataMap["updated_at"],
			},
		}
		expectedDataJson, _ := json.Marshal(expectedData)
		assert.JSONEq(t, string(expectedDataJson), string(data))

		// You can test for other error responses yourself

	})
}

func deleteMessage(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	message := CreateMessage(db)
	token := AccessToken(db)
	t.Run("Delete A Message", func(t *testing.T) {
		url := fmt.Sprintf("%s/messages/%s", baseUrl, uuid.New())
		// Test for valid response for invalid message id
		req := httptest.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		res, _ := app.Test(req)

		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)
		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_NON_EXISTENT, body["code"])
		assert.Equal(t, "User has no message with that ID", body["message"])

		// Test for valid response for valid entry
		url = fmt.Sprintf("%s/messages/%s", baseUrl, message.ID)
		req = httptest.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		res, _ = app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		data, _ := json.Marshal(body)
		expectedData := map[string]interface{}{
			"status":  "success",
			"message": "Message Deleted",
		}
		expectedDataJson, _ := json.Marshal(expectedData)
		assert.JSONEq(t, string(expectedDataJson), string(data))
	})
}

func createGroupChat(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	user := CreateAnotherTestVerifiedUser(db)
	owner := CreateTestVerifiedUser(db)
	token := AccessToken(db)
	t.Run("Create A Group Chat", func(t *testing.T) {
		url := fmt.Sprintf("%s/groups/group", baseUrl)
		desc := "JESUS is KING"
		chatData := schemas.GroupChatCreateSchema{Name: "New Group Chat", Description: &desc, UsernamesToAdd: []string{"invalid_username"}}

		// Verify the requests fails with invalid username
		res := ProcessTestBody(t, app, url, "POST", chatData, token)
		// Assert Status code
		assert.Equal(t, 422, res.StatusCode)
		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, utils.ERR_INVALID_ENTRY, body["code"])
		assert.Equal(t, "Invalid Entry", body["message"])
		expectedDataErr := map[string]interface{}{"usernames_to_add": "Enter at least one valid username"}
		assert.Equal(t, expectedDataErr, body["data"])

		// Test for valid response for valid entry
		chatData.UsernamesToAdd = []string{user.Username}
		url = fmt.Sprintf("%s/groups/group", baseUrl)
		res = ProcessTestBody(t, app, url, "POST", chatData, token)
		// Assert Status code
		assert.Equal(t, 201, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		data, _ := json.Marshal(body)
		dataBody := body["data"].(map[string]interface{})

		expectedData := map[string]interface{}{
			"status":  "success",
			"message": "Chat created",
			"data": map[string]interface{}{
				"id":          dataBody["id"],
				"owner": GetUserMap(owner),
				"name":        chatData.Name,
				"description": chatData.Description,
				"ctype": choices.CGROUP,
				"users": []map[string]interface{}{
					GetUserMap(user),
				},
				"image": nil,
				"latest_message": nil,
				"created_at": dataBody["created_at"],
				"updated_at": dataBody["updated_at"],
			},
		}
		expectedDataJson, _ := json.Marshal(expectedData)
		assert.JSONEq(t, string(expectedDataJson), string(data))
	})
}

func TestChat(t *testing.T) {
	os.Setenv("ENVIRONMENT", "TESTING")
	app := fiber.New()
	db := Setup(t, app)
	BASEURL := "/api/v6/chats"

	// Run Chat Endpoint Tests
	getChats(t, app, db, BASEURL)
	sendMessage(t, app, db, BASEURL)
	getChatMessages(t, app, db, BASEURL)
	updateGroupChat(t, app, db, BASEURL)
	deleteGroupChat(t, app, db, BASEURL)
	updateMessage(t, app, db, BASEURL)
	deleteMessage(t, app, db, BASEURL)
	createGroupChat(t, app, db, BASEURL)

	// Drop Tables and Close Connectiom
	database.DropTables(db)
	CloseTestDatabase(db)
}
