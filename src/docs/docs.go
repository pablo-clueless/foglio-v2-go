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
        "/api/v2/test/email": {
            "get": {
                "summary": "Test endpoint",
                "description": "Test email service",
                "tags": ["Test", "Email"],
                "responses": {
                    "200": {
                        "description": "Test endpoint response"
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
        "/api/v2/auth/request-verification": {
            "post": {
                "summary": "Request verification",
                "description": "Request a new verification OTP",
                "tags": ["Authentication"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "email",
                        "in": "query",
                        "required": true,
                        "type": "string",
                        "description": "User email address"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Verification OTP sent"
                    },
                    "400": {
                        "description": "Bad request"
                    }
                }
            }
        },
        "/api/v2/auth/{provider}": {
            "get": {
                "summary": "Get OAuth URL",
                "description": "Get OAuth authorization URL for provider",
                "tags": ["Authentication"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "provider",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "enum": ["google", "github"],
                        "description": "OAuth provider"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OAuth URL generated"
                    }
                }
            }
        },
        "/api/v2/auth/{provider}/callback": {
            "get": {
                "summary": "OAuth callback",
                "description": "Handle OAuth provider callback",
                "tags": ["Authentication"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "provider",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "enum": ["google", "github"],
                        "description": "OAuth provider"
                    },
                    {
                        "name": "code",
                        "in": "query",
                        "required": true,
                        "type": "string",
                        "description": "Authorization code"
                    },
                    {
                        "name": "state",
                        "in": "query",
                        "type": "string",
                        "description": "State parameter"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Authentication successful"
                    },
                    "400": {
                        "description": "Bad request"
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
        "/api/v2/me": {
            "get": {
                "summary": "Get current user",
                "description": "Get authenticated user's profile",
                "tags": ["Users"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "User profile"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/me/jobs": {
            "get": {
                "summary": "Get my jobs",
                "description": "Get jobs created by authenticated user",
                "tags": ["Jobs"],
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
                        "name": "limit",
                        "in": "query",
                        "type": "integer",
                        "description": "Items per page"
                    }
                ],
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
        "/api/v2/user/profile": {
            "get": {
                "summary": "Get user profile",
                "description": "Get authenticated user's profile (alias for /me)",
                "tags": ["Users"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "User profile"
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
        "/api/v2/jobs/applications/user": {
            "get": {
                "summary": "Get my applications",
                "description": "Get all job applications submitted by authenticated user",
                "tags": ["Job Applications"],
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
                        "name": "limit",
                        "in": "query",
                        "type": "integer",
                        "description": "Items per page"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of applications"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/jobs/applications/job/{id}": {
            "get": {
                "summary": "Get applications for job",
                "description": "Get all applications for a specific job (recruiter only)",
                "tags": ["Job Applications"],
                "security": [{"Bearer": []}],
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
                        "name": "page",
                        "in": "query",
                        "type": "integer",
                        "description": "Page number"
                    },
                    {
                        "name": "limit",
                        "in": "query",
                        "type": "integer",
                        "description": "Items per page"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of applications"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "403": {
                        "description": "Forbidden - not a recruiter"
                    }
                }
            }
        },
        "/api/v2/jobs/applications/{applicationId}": {
            "get": {
                "summary": "Get application",
                "description": "Get a specific job application by ID",
                "tags": ["Job Applications"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "applicationId",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Application UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Application details"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Application not found"
                    }
                }
            }
        },
        "/api/v2/jobs/applications/{id}/accept": {
            "post": {
                "summary": "Accept application",
                "description": "Accept a job application (recruiter only)",
                "tags": ["Job Applications"],
                "security": [{"Bearer": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Application UUID"
                    },
                    {
                        "in": "body",
                        "name": "body",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "reason": {
                                    "type": "string",
                                    "description": "Optional reason/message"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Application accepted"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "403": {
                        "description": "Forbidden - not a recruiter"
                    }
                }
            }
        },
        "/api/v2/jobs/applications/{id}/reject": {
            "post": {
                "summary": "Reject application",
                "description": "Reject a job application (recruiter only)",
                "tags": ["Job Applications"],
                "security": [{"Bearer": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Application UUID"
                    },
                    {
                        "in": "body",
                        "name": "body",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "reason": {
                                    "type": "string",
                                    "description": "Optional rejection reason"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Application rejected"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "403": {
                        "description": "Forbidden - not a recruiter"
                    }
                }
            }
        },
        "/api/v2/jobs/{id}/comment": {
            "post": {
                "summary": "Add comment",
                "description": "Add a comment to a job posting",
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
                            "required": ["content"],
                            "properties": {
                                "content": {
                                    "type": "string",
                                    "example": "Great opportunity!"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Comment added"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            },
            "delete": {
                "summary": "Delete comment",
                "description": "Delete a comment from a job posting",
                "tags": ["Jobs"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Comment UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Comment deleted"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Comment not found"
                    }
                }
            }
        },
        "/api/v2/jobs/{id}/reaction/{reaction}": {
            "post": {
                "summary": "Add reaction",
                "description": "Add a reaction to a job posting",
                "tags": ["Jobs"],
                "security": [{"Bearer": []}],
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
                        "name": "reaction",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "enum": ["LIKE", "DISLIKE"],
                        "description": "Reaction type"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Reaction added"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/jobs/{id}/reaction": {
            "delete": {
                "summary": "Remove reaction",
                "description": "Remove reaction from a job posting",
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
                        "description": "Reaction removed"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Reaction not found"
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
        },
        "/api/v2/subscriptions": {
            "get": {
                "summary": "List subscription tiers",
                "description": "Get all available subscription tiers",
                "tags": ["Subscriptions"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "page",
                        "in": "query",
                        "type": "integer",
                        "description": "Page number"
                    },
                    {
                        "name": "limit",
                        "in": "query",
                        "type": "integer",
                        "description": "Items per page"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of subscription tiers"
                    }
                }
            },
            "post": {
                "summary": "Create subscription tier",
                "description": "Create a new subscription tier (admin only)",
                "tags": ["Subscriptions"],
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
                            "required": ["name", "type", "tier", "price", "billingCycleDays"],
                            "properties": {
                                "name": {"type": "string", "example": "Pro Plan"},
                                "description": {"type": "string", "example": "Best for professionals"},
                                "type": {"type": "string", "enum": ["monthly", "yearly", "lifetime"], "example": "monthly"},
                                "tier": {"type": "string", "enum": ["free", "basic", "premium", "business"], "example": "premium"},
                                "price": {"type": "number", "example": 9.99},
                                "currency": {"type": "string", "example": "USD"},
                                "billingCycleDays": {"type": "integer", "example": 30},
                                "trialPeriodDays": {"type": "integer", "example": 14},
                                "features": {"type": "object"},
                                "sortOrder": {"type": "integer", "example": 1}
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Subscription tier created"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/subscriptions/{id}": {
            "get": {
                "summary": "Get subscription tier",
                "description": "Get subscription tier by ID",
                "tags": ["Subscriptions"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Subscription UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Subscription tier details"
                    },
                    "404": {
                        "description": "Subscription not found"
                    }
                }
            },
            "put": {
                "summary": "Update subscription tier",
                "description": "Update subscription tier by ID (admin only)",
                "tags": ["Subscriptions"],
                "security": [{"Bearer": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Subscription UUID"
                    },
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "name": {"type": "string"},
                                "description": {"type": "string"},
                                "type": {"type": "string", "enum": ["monthly", "yearly", "lifetime"]},
                                "tier": {"type": "string", "enum": ["free", "basic", "premium", "business"]},
                                "price": {"type": "number"},
                                "currency": {"type": "string"},
                                "billingCycleDays": {"type": "integer"},
                                "trialPeriodDays": {"type": "integer"},
                                "features": {"type": "object"},
                                "sortOrder": {"type": "integer"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Subscription tier updated"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Subscription not found"
                    }
                }
            },
            "delete": {
                "summary": "Delete subscription tier",
                "description": "Delete subscription tier by ID (admin only)",
                "tags": ["Subscriptions"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Subscription UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Subscription tier deleted"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Subscription not found"
                    }
                }
            }
        },
        "/api/v2/user/subscriptions": {
            "get": {
                "summary": "Get user subscriptions",
                "description": "Get current user's subscription history",
                "tags": ["User Subscriptions"],
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
                        "name": "limit",
                        "in": "query",
                        "type": "integer",
                        "description": "Items per page"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User's subscription list"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/user/subscriptions/{id}": {
            "get": {
                "summary": "Get user subscription",
                "description": "Get user subscription by ID",
                "tags": ["User Subscriptions"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "User Subscription UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User subscription details"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Subscription not found"
                    }
                }
            }
        },
        "/api/v2/user/subscriptions/{tierId}/subscribe": {
            "post": {
                "summary": "Subscribe to tier",
                "description": "Subscribe current user to a subscription tier",
                "tags": ["User Subscriptions"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "tierId",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Subscription Tier UUID"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Subscribed successfully"
                    },
                    "400": {
                        "description": "User already has an active subscription"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/user/subscriptions/{tierId}/upgrade": {
            "put": {
                "summary": "Upgrade subscription",
                "description": "Upgrade current user's subscription to a higher tier",
                "tags": ["User Subscriptions"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "tierId",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "New Subscription Tier UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Subscription upgraded successfully"
                    },
                    "400": {
                        "description": "No active subscription found"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/user/subscriptions/{tierId}/downgrade": {
            "put": {
                "summary": "Downgrade subscription",
                "description": "Downgrade current user's subscription to a lower tier",
                "tags": ["User Subscriptions"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "tierId",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "New Subscription Tier UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Subscription downgraded successfully"
                    },
                    "400": {
                        "description": "No active subscription found"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/user/subscriptions/unsubscribe": {
            "delete": {
                "summary": "Unsubscribe",
                "description": "Cancel current user's active subscription",
                "tags": ["User Subscriptions"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Unsubscribed successfully"
                    },
                    "400": {
                        "description": "No active subscription found"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/payments/initialize": {
            "post": {
                "summary": "Initialize payment",
                "description": "Initialize a Paystack payment transaction for subscription",
                "tags": ["Payments"],
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
                            "required": ["subscription_tier_id"],
                            "properties": {
                                "subscription_tier_id": {
                                    "type": "string",
                                    "description": "UUID of the subscription tier to purchase",
                                    "example": "123e4567-e89b-12d3-a456-426614174000"
                                },
                                "callback_url": {
                                    "type": "string",
                                    "description": "URL to redirect after payment",
                                    "example": "https://yoursite.com/payment/callback"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Payment initialized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "authorization_url": {
                                    "type": "string",
                                    "description": "URL to redirect user for payment"
                                },
                                "access_code": {
                                    "type": "string"
                                },
                                "reference": {
                                    "type": "string",
                                    "description": "Transaction reference for verification"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad request - user already has subscription or invalid tier"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/payments/verify": {
            "get": {
                "summary": "Verify payment",
                "description": "Verify a Paystack payment transaction and activate subscription",
                "tags": ["Payments"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "reference",
                        "in": "query",
                        "required": true,
                        "type": "string",
                        "description": "Payment reference from Paystack"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Payment verified and subscription activated"
                    },
                    "400": {
                        "description": "Payment not successful or already processed"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/payments/cancel": {
            "delete": {
                "summary": "Cancel subscription",
                "description": "Cancel the user's active Paystack subscription",
                "tags": ["Payments"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Subscription cancelled successfully"
                    },
                    "400": {
                        "description": "No active subscription found"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/payments/webhook": {
            "post": {
                "summary": "Paystack webhook",
                "description": "Webhook endpoint for Paystack events (charge.success, subscription.create, etc.)",
                "tags": ["Payments"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "x-paystack-signature",
                        "in": "header",
                        "required": true,
                        "type": "string",
                        "description": "HMAC SHA512 signature for verification"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Webhook processed successfully"
                    },
                    "400": {
                        "description": "Invalid request body"
                    },
                    "401": {
                        "description": "Invalid signature"
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
