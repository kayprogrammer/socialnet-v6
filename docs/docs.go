// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "security": [
                    {
                        "GuestUserAuth": []
                    }
                ],
                "description": "This endpoint generates new access and refresh tokens for authentication",
                "tags": [
                    "Auth"
                ],
                "summary": "Login a user",
                "parameters": [
                    {
                        "description": "User login",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.LoginSchema"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "This endpoint logs a user out from our application",
                "tags": [
                    "Auth"
                ],
                "summary": "Logout a user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/refresh": {
            "post": {
                "description": "This endpoint refresh tokens by generating new access and refresh tokens for a user",
                "tags": [
                    "Auth"
                ],
                "summary": "Refresh tokens",
                "parameters": [
                    {
                        "description": "Refresh token",
                        "name": "refresh",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.RefreshTokenSchema"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "` + "`" + `This endpoint registers new users into our application.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.RegisterUser"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/schemas.RegisterResponseSchema"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/resend-verification-email": {
            "post": {
                "description": "` + "`" + `This endpoint resends new otp to the user's email.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Resend Verification Email",
                "parameters": [
                    {
                        "description": "Email data",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.EmailRequestSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/send-password-reset-otp": {
            "post": {
                "description": "` + "`" + `This endpoint sends new password reset otp to the user's email.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Send Password Reset Otp",
                "parameters": [
                    {
                        "description": "Email object",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.EmailRequestSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/set-new-password": {
            "post": {
                "description": "` + "`" + `This endpoint verifies the password reset otp.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Set New Password",
                "parameters": [
                    {
                        "description": "Password reset object",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.SetNewPasswordSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/verify-email": {
            "post": {
                "description": "` + "`" + `This endpoint verifies a user's email.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Verify a user's email",
                "parameters": [
                    {
                        "description": "Verify Email object",
                        "name": "verify_email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.VerifyEmailRequestSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/general/site-detail": {
            "get": {
                "description": "This endpoint retrieves few details of the site/application.",
                "tags": [
                    "General"
                ],
                "summary": "Retrieve site details",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.SiteDetailResponseSchema"
                        }
                    }
                }
            }
        },
        "/healthcheck": {
            "get": {
                "description": "This endpoint checks the health of our application.",
                "tags": [
                    "HealthCheck"
                ],
                "summary": "HealthCheck",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.HealthCheckSchema"
                        }
                    }
                }
            }
        },
        "/profiles/cities": {
            "get": {
                "description": "This endpoint retrieves the first 10 cities that matches the query params",
                "tags": [
                    "Profiles"
                ],
                "summary": "Retrieve cities based on query params",
                "parameters": [
                    {
                        "type": "string",
                        "description": "City name",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.CitiesResponseSchema"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.City": {
            "type": "object",
            "properties": {
                "country": {
                    "type": "string",
                    "example": "Nigeria"
                },
                "countryObj": {
                    "$ref": "#/definitions/models.Country"
                },
                "country_id": {
                    "type": "string"
                },
                "id": {
                    "type": "string",
                    "example": "d10dde64-a242-4ed0-bd75-4c759644b3a6"
                },
                "name": {
                    "type": "string",
                    "example": "Lekki"
                },
                "region": {
                    "type": "string",
                    "example": "Lagos"
                },
                "regionObj": {
                    "$ref": "#/definitions/models.Region"
                },
                "region_id": {
                    "type": "string"
                }
            }
        },
        "models.Country": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string",
                    "example": "NG"
                },
                "id": {
                    "type": "string",
                    "example": "d10dde64-a242-4ed0-bd75-4c759644b3a6"
                },
                "name": {
                    "type": "string",
                    "example": "Nigeria"
                }
            }
        },
        "models.Region": {
            "type": "object",
            "properties": {
                "country": {
                    "$ref": "#/definitions/models.Country"
                },
                "country_id": {
                    "type": "string"
                },
                "id": {
                    "type": "string",
                    "example": "d10dde64-a242-4ed0-bd75-4c759644b3a6"
                },
                "name": {
                    "type": "string",
                    "example": "Lagos"
                }
            }
        },
        "models.SiteDetail": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string",
                    "example": "234, Lagos, Nigeria"
                },
                "email": {
                    "type": "string",
                    "example": "kayprogrammer1@gmail.com"
                },
                "fb": {
                    "type": "string",
                    "example": "https://facebook.com"
                },
                "id": {
                    "type": "string",
                    "example": "d10dde64-a242-4ed0-bd75-4c759644b3a6"
                },
                "ig": {
                    "type": "string",
                    "example": "https://instagram.com"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string",
                    "example": "+2348133831036"
                },
                "tw": {
                    "type": "string",
                    "example": "https://twitter.com"
                },
                "wh": {
                    "type": "string",
                    "example": "https://wa.me/2348133831036"
                }
            }
        },
        "routes.HealthCheckSchema": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "string",
                    "example": "pong"
                }
            }
        },
        "schemas.CitiesResponseSchema": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.City"
                    }
                },
                "message": {
                    "type": "string",
                    "example": "Data fetched/created/updated/deleted"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "schemas.EmailRequestSchema": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 5,
                    "example": "johndoe@email.com"
                }
            }
        },
        "schemas.LoginSchema": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "johndoe@email.com"
                },
                "password": {
                    "type": "string",
                    "example": "password"
                }
            }
        },
        "schemas.RefreshTokenSchema": {
            "type": "object",
            "required": [
                "refresh"
            ],
            "properties": {
                "refresh": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"
                }
            }
        },
        "schemas.RegisterResponseSchema": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/schemas.EmailRequestSchema"
                },
                "message": {
                    "type": "string",
                    "example": "Data fetched/created/updated/deleted"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "schemas.RegisterUser": {
            "type": "object",
            "required": [
                "email",
                "first_name",
                "last_name",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 5,
                    "example": "johndoe@email.com"
                },
                "first_name": {
                    "type": "string",
                    "maxLength": 50,
                    "example": "John"
                },
                "last_name": {
                    "type": "string",
                    "maxLength": 50,
                    "example": "Doe"
                },
                "password": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 8,
                    "example": "strongpassword"
                },
                "terms_agreement": {
                    "type": "boolean"
                }
            }
        },
        "schemas.ResponseSchema": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Data fetched/created/updated/deleted"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "schemas.SetNewPasswordSchema": {
            "type": "object",
            "required": [
                "email",
                "otp",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 5,
                    "example": "johndoe@email.com"
                },
                "otp": {
                    "type": "integer",
                    "example": 123456
                },
                "password": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 8,
                    "example": "newstrongpassword"
                }
            }
        },
        "schemas.SiteDetailResponseSchema": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.SiteDetail"
                },
                "message": {
                    "type": "string",
                    "example": "Data fetched/created/updated/deleted"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "schemas.VerifyEmailRequestSchema": {
            "type": "object",
            "required": [
                "email",
                "otp"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 5,
                    "example": "johndoe@email.com"
                },
                "otp": {
                    "type": "integer",
                    "example": 123456
                }
            }
        },
        "utils.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "data": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "Type 'Bearer jwt_string' to correctly set the API Key",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "6.0",
	Host:             "",
	BasePath:         "/api/v6",
	Schemes:          []string{},
	Title:            "SOCIALNET API",
	Description:      "## A Realtime Social Networking API built with FIBER & GORM ORM.\n\n### WEBSOCKETS:\n\n#### Notifications\n\n- URL: `wss://{host}/api/v6/ws/notifications`\n\n- Requires authorization, so pass in the Bearer Authorization header.\n\n- You can only read and not send notification messages into this socket.\n\n\n#### Chats\n\n- URL: `wss://{host}/api/v6/ws/chats/{id}`\n- Requires authorization, so pass in the Bearer Authorization header.\n- Use chat_id as the ID for an existing chat or username if it's the first message in a DM.\n- You cannot read realtime messages from a username that doesn't belong to the authorized user, but you can surely send messages.\n- Only send a message to the socket endpoint after the message has been created or updated, and files have been uploaded.\n- Fields when sending a message through the socket:\n\n  ```json\n  { \"status\": \"CREATED\", \"id\": \"fe4e0235-80fc-4c94-b15e-3da63226f8ab\" }\n  ```\n",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
