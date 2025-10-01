package docs

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
)

const template = `{
	"schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "description": "Enter your bearer token in the format: Bearer {token}"
        }
    },
    "paths": {
        "/api/v2/health": {
            "get": {
                "summary": "Health check",
                "description": "Check service health",
                "tags": ["Health"],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "status": {
                                    "type": "string",
                                    "example": "healthy"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/v2/auth/signup": {
            "post": {
                "summary": "User signup",
                "description": "Create a new user account",
                "tags": ["Authentication"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["name", "email", "password"],
                            "properties": {
                                "name": {
                                    "type": "string",
                                    "example": "John Doe"
                                },
                                "email": {
                                    "type": "string",
                                    "example": "john@example.com"
                                },
                                "password": {
                                    "type": "string",
                                    "example": "password123"
                                },
                                "isRecruiter": {
                                    "type": "boolean",
                                    "example": false
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User created successfully"
                    },
                    "400": {
                        "description": "Bad request"
                    }
                }
            }
        },
        "/api/v2/auth/signin": {
            "post": {
                "summary": "User signin",
                "description": "Authenticate user and return token",
                "tags": ["Authentication"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["email", "password"],
                            "properties": {
                                "email": {
                                    "type": "string",
                                    "example": "john@example.com"
                                },
                                "password": {
                                    "type": "string",
                                    "example": "password123"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Authentication successful",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "token": {
                                    "type": "string"
                                },
                                "user": {
                                    "type": "object"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/auth/verification": {
            "post": {
                "summary": "Account verification",
                "description": "Verify user account with OTP",
                "tags": ["Authentication"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["email", "otp"],
                            "properties": {
                                "email": {
                                    "type": "string",
                                    "example": "john@example.com"
                                },
                                "otp": {
                                    "type": "string",
                                    "example": "123456"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Verification successful"
                    },
                    "400": {
                        "description": "Invalid OTP"
                    }
                }
            }
        },
        "/api/v2/auth/forgot-password": {
            "post": {
                "summary": "Forgot password",
                "description": "Request password reset OTP",
                "tags": ["Authentication"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["email"],
                            "properties": {
                                "email": {
                                    "type": "string",
                                    "example": "john@example.com"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Reset OTP sent"
                    }
                }
            }
        },
        "/api/v2/auth/reset-password": {
            "post": {
                "summary": "Reset password",
                "description": "Reset user password with OTP",
                "tags": ["Authentication"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["email", "otp", "newPassword"],
                            "properties": {
                                "email": {
                                    "type": "string",
                                    "example": "john@example.com"
                                },
                                "otp": {
                                    "type": "string",
                                    "example": "123456"
                                },
                                "newPassword": {
                                    "type": "string",
                                    "example": "newPassword123"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Password reset successful"
                    }
                }
            }
        },
        "/api/v2/auth/update-password": {
            "post": {
                "summary": "Update password",
                "description": "Change user password",
                "tags": ["Authentication"],
                "security": [{"Bearer": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["currentPassword", "newPassword"],
                            "properties": {
                                "currentPassword": {
                                    "type": "string",
                                    "example": "oldPassword123"
                                },
                                "newPassword": {
                                    "type": "string",
                                    "example": "newPassword123"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Password updated"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/users/": {
            "get": {
                "summary": "List users",
                "description": "Get all users",
                "tags": ["Users"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "List of users"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/users/{id}": {
            "get": {
                "summary": "Get user",
                "description": "Get user by ID",
                "tags": ["Users"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "User UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User details"
                    },
                    "404": {
                        "description": "User not found"
                    }
                }
            },
            "put": {
                "summary": "Update user",
                "description": "Update user by ID",
                "tags": ["Users"],
                "security": [{"Bearer": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "User UUID"
                    },
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "name": {"type": "string"},
                                "headline": {"type": "string"},
                                "location": {"type": "string"},
                                "summary": {"type": "string"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User updated"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            },
            "delete": {
                "summary": "Delete user",
                "description": "Delete user by ID",
                "tags": ["Users"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "User UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User deleted"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/users/{id}/avatar": {
            "put": {
                "summary": "Update avatar",
                "description": "Update user avatar",
                "tags": ["Users"],
                "security": [{"Bearer": []}],
                "consumes": ["multipart/form-data"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "User UUID"
                    },
                    {
                        "name": "avatar",
                        "in": "formData",
                        "required": true,
                        "type": "file",
                        "description": "Avatar image file"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Avatar updated"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/jobs/": {
            "post": {
                "summary": "Create job",
                "description": "Create a new job posting",
                "tags": ["Jobs"],
                "security": [{"Bearer": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["title", "company", "location", "description", "employmentType"],
                            "properties": {
                                "title": {"type": "string", "example": "Senior Software Engineer"},
                                "company": {"type": "string", "example": "Tech Corp"},
                                "location": {"type": "string", "example": "San Francisco, CA"},
                                "description": {"type": "string"},
                                "requirements": {"type": "array", "items": {"type": "string"}},
                                "employmentType": {"type": "string", "example": "Full-time"},
                                "isRemote": {"type": "boolean", "example": true}
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Job created"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            },
            "get": {
                "summary": "List jobs",
                "description": "Get all jobs with optional filters",
                "tags": ["Jobs"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "location",
                        "in": "query",
                        "type": "string",
                        "description": "Filter by location"
                    },
                    {
                        "name": "isRemote",
                        "in": "query",
                        "type": "boolean",
                        "description": "Filter by remote jobs"
                    },
                    {
                        "name": "page",
                        "in": "query",
                        "type": "integer",
                        "description": "Page number"
                    },
                    {
                        "name": "size",
                        "in": "query",
                        "type": "integer",
                        "description": "Items per page"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of jobs"
                    }
                }
            }
        },
        "/api/v2/jobs/me": {
            "get": {
                "summary": "List my jobs",
                "description": "Get all jobs created by XXX authenticated user",
                "tags": ["Jobs"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "List of user's jobs"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/jobs/{id}": {
            "get": {
                "summary": "Get job",
                "description": "Get job by ID",
                "tags": ["Jobs"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Job UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Job details"
                    },
                    "404": {
                        "description": "Job not found"
                    }
                }
            },
            "put": {
                "summary": "Update job",
                "description": "Update job by ID",
                "tags": ["Jobs"],
                "security": [{"Bearer": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Job UUID"
                    },
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "title": {"type": "string"},
                                "description": {"type": "string"},
                                "isRemote": {"type": "boolean"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Job updated"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            },
            "delete": {
                "summary": "Delete job",
                "description": "Delete job by ID",
                "tags": ["Jobs"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Job UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Job deleted"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/jobs/{id}/apply": {
            "post": {
                "summary": "Apply to job",
                "description": "Submit job application",
                "tags": ["Jobs"],
                "security": [{"Bearer": []}],
                "consumes": ["multipart/form-data"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Job UUID"
                    },
                    {
                        "name": "resume",
                        "in": "formData",
                        "required": true,
                        "type": "file",
                        "description": "Resume file"
                    },
                    {
                        "name": "coverLetter",
                        "in": "formData",
                        "type": "string",
                        "description": "Cover letter text"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Application submitted"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/notifications": {
            "get": {
                "summary": "Get notifications",
                "description": "Get all notification",
                "tags": ["Notifications"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "page",
                        "in": "query",
                        "type": "integer",
                        "description": "Page number"
                    },
                    {
                        "name": "size",
                        "in": "query",
                        "type": "integer",
                        "description": "Items per page"
                    }
                ]
            }
        },
        "/api/v2/notifications/{id}": {
            "get": {
                "summary": "Get notification",
                "description": "Get notification by ID",
                "tags": ["Notifications"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Notification UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Notification details"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Notification not found"
                    }
                }
            },
            "put": {
                "summary": "Update notification",
                "description": "Mark notification as read",
                "tags": ["Notifications"],
                "security": [{"Bearer": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Notification UUID"
                    },
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "isRead": {
                                    "type": "boolean",
                                    "example": true
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Notification updated"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Notification not found"
                    }
                }
            },
            "delete": {
                "summary": "Delete notification",
                "description": "Delete notification by ID",
                "tags": ["Notifications"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Notification UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Notification deleted"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Notification not found"
                    }
                }
            }
        }
    }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "Foglio API",
	Description:      "Professional networking and job platform API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  template,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func InitializeSwagger() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

func SetupSwagger(router *gin.Engine) {
	InitializeSwagger()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/docs", func(ctx *gin.Context) {
		ctx.Redirect(302, "/swagger/index.html")
	})
}
