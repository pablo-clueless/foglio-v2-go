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
        "/api/v2/ws": {
            "get": {
                "summary": "WebSocket connection",
                "description": "Establish a WebSocket connection for real-time notifications and chat messaging. Supported actions: send_message, typing, stop_typing, mark_messages_read, mark_read, ping. Messages are received as notifications with event_type in data field.",
                "tags": ["WebSocket"],
                "security": [{"Bearer": []}],
                "responses": {
                    "101": {"description": "Switching Protocols - WebSocket connection established"},
                    "400": {"description": "Could not upgrade connection"},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/api/v2/ws/stats": {
            "get": {
                "summary": "WebSocket stats",
                "description": "Get current WebSocket connection statistics",
                "tags": ["WebSocket"],
                "responses": {
                    "200": {
                        "description": "Connection statistics",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "connected_clients": {"type": "integer", "description": "Total number of connected clients"},
                                "connected_users": {"type": "integer", "description": "Number of unique users connected"}
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
                        "description": "Authentication successful (or 2FA required)",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "token": {
                                    "type": "string",
                                    "description": "JWT token (empty if 2FA required)"
                                },
                                "user": {
                                    "type": "object",
                                    "description": "User object (empty if 2FA required)"
                                },
                                "requires_two_factor": {
                                    "type": "boolean",
                                    "description": "True if 2FA verification is needed"
                                },
                                "user_id": {
                                    "type": "string",
                                    "description": "User ID for 2FA verification (only if 2FA required)"
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
        "/api/v2/auth/2fa/setup": {
            "post": {
                "summary": "Start 2FA setup",
                "description": "Generate a TOTP secret and QR code URL for setting up 2FA",
                "tags": ["Two-Factor Authentication"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "2FA setup initiated",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "secret": {
                                    "type": "string",
                                    "description": "TOTP secret key"
                                },
                                "qr_code_url": {
                                    "type": "string",
                                    "description": "URL for QR code (otpauth://)"
                                },
                                "backup_codes": {
                                    "type": "array",
                                    "items": {"type": "string"},
                                    "description": "One-time backup codes"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "2FA already enabled"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/auth/2fa/verify-setup": {
            "post": {
                "summary": "Verify and enable 2FA",
                "description": "Verify the TOTP code from the authenticator app and enable 2FA",
                "tags": ["Two-Factor Authentication"],
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
                            "required": ["code"],
                            "properties": {
                                "code": {
                                    "type": "string",
                                    "example": "123456",
                                    "description": "6-digit TOTP code from authenticator app"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "2FA enabled successfully"
                    },
                    "400": {
                        "description": "Invalid verification code"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/auth/2fa/verify": {
            "post": {
                "summary": "Verify 2FA during login",
                "description": "Verify the TOTP code or backup code during login to get the auth token",
                "tags": ["Two-Factor Authentication"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["user_id", "code"],
                            "properties": {
                                "user_id": {
                                    "type": "string",
                                    "description": "User ID returned from signin"
                                },
                                "code": {
                                    "type": "string",
                                    "example": "123456",
                                    "description": "6-digit TOTP code or backup code"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "2FA verification successful",
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
                        "description": "Invalid verification code"
                    }
                }
            }
        },
        "/api/v2/auth/2fa/disable": {
            "post": {
                "summary": "Disable 2FA",
                "description": "Disable two-factor authentication (requires password confirmation)",
                "tags": ["Two-Factor Authentication"],
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
                            "required": ["password"],
                            "properties": {
                                "password": {
                                    "type": "string",
                                    "description": "Current password for confirmation"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "2FA disabled successfully"
                    },
                    "400": {
                        "description": "Invalid password or 2FA not enabled"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/auth/2fa/backup-codes": {
            "post": {
                "summary": "Regenerate backup codes",
                "description": "Generate new backup codes (invalidates old ones, requires password)",
                "tags": ["Two-Factor Authentication"],
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
                            "required": ["password"],
                            "properties": {
                                "password": {
                                    "type": "string",
                                    "description": "Current password for confirmation"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "New backup codes generated",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "backup_codes": {
                                    "type": "array",
                                    "items": {"type": "string"}
                                },
                                "message": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid password or 2FA not enabled"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/auth/2fa/status": {
            "get": {
                "summary": "Get 2FA status",
                "description": "Get the current 2FA status for the authenticated user",
                "tags": ["Two-Factor Authentication"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "2FA status retrieved",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "enabled": {
                                    "type": "boolean"
                                },
                                "backup_codes_left": {
                                    "type": "integer"
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
                    },
                    {
                        "name": "status",
                        "in": "query",
                        "type": "string",
                        "enum": ["PENDING", "REVIEWED", "ACCEPTED", "REJECTED", "HIRED"],
                        "description": "Filter by application status"
                    },
                    {
                        "name": "submission_date",
                        "in": "query",
                        "type": "string",
                        "format": "date",
                        "description": "Filter by submission date (YYYY-MM-DD)"
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
        "/api/v2/jobs/applications/{id}/review": {
            "post": {
                "summary": "Mark application as reviewed",
                "description": "Mark a job application as reviewed (recruiter only)",
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
                                    "description": "Optional notes"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Application marked as reviewed"
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
        "/api/v2/jobs/applications/{id}/hire": {
            "post": {
                "summary": "Hire applicant",
                "description": "Mark applicant as hired (recruiter only)",
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
                                    "description": "Optional message to applicant"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Applicant hired"
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
        },
        "/api/v2/payments/methods": {
            "get": {
                "summary": "Get payment methods",
                "description": "Get all saved payment methods for the authenticated user",
                "tags": ["Payment Methods"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "List of payment methods"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            },
            "post": {
                "summary": "Add payment method",
                "description": "Initialize adding a new payment method (card)",
                "tags": ["Payment Methods"],
                "security": [{"Bearer": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "callback_url": {
                                    "type": "string",
                                    "description": "URL to redirect after card validation"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Redirect URL for card validation"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/payments/methods/{authCode}": {
            "delete": {
                "summary": "Remove payment method",
                "description": "Remove a saved payment method",
                "tags": ["Payment Methods"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {
                        "name": "authCode",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Authorization code of the payment method"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Payment method removed"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "400": {
                        "description": "Bad request"
                    }
                }
            }
        },
        "/api/v2/payments/invoices": {
            "get": {
                "summary": "Get invoices",
                "description": "Get all invoices for the authenticated user",
                "tags": ["Invoices"],
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
                        "description": "List of invoices"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/payments/invoices/{id}": {
            "get": {
                "summary": "Get invoice",
                "description": "Get invoice by ID",
                "tags": ["Invoices"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Invoice UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Invoice details"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Invoice not found"
                    }
                }
            }
        },
        "/api/v2/domain": {
            "get": {
                "summary": "Get domain configuration",
                "description": "Get the authenticated user's domain configuration including subdomain and custom domain",
                "tags": ["Domain"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Domain configuration retrieved"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/domain/check/{subdomain}": {
            "get": {
                "summary": "Check subdomain availability",
                "description": "Check if a subdomain is available for use",
                "tags": ["Domain"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "subdomain",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Subdomain to check"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Availability status"
                    }
                }
            }
        },
        "/api/v2/domain/subdomain": {
            "post": {
                "summary": "Claim subdomain",
                "description": "Claim a subdomain for your profile (e.g., username.foglio.app)",
                "tags": ["Domain"],
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
                            "required": ["subdomain"],
                            "properties": {
                                "subdomain": {
                                    "type": "string",
                                    "example": "johndoe",
                                    "description": "Subdomain (3-32 characters, alphanumeric)"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Subdomain claimed successfully"
                    },
                    "400": {
                        "description": "Subdomain taken, invalid, or reserved"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            },
            "put": {
                "summary": "Update subdomain",
                "description": "Update your subdomain",
                "tags": ["Domain"],
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
                            "required": ["subdomain"],
                            "properties": {
                                "subdomain": {
                                    "type": "string",
                                    "example": "newsubdomain"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Subdomain updated successfully"
                    },
                    "400": {
                        "description": "Subdomain taken or invalid"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/domain/custom": {
            "post": {
                "summary": "Set custom domain",
                "description": "Set a custom domain for your profile (paid users only)",
                "tags": ["Domain"],
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
                            "required": ["custom_domain"],
                            "properties": {
                                "custom_domain": {
                                    "type": "string",
                                    "example": "portfolio.example.com",
                                    "description": "Your custom domain"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Custom domain configured, DNS records provided for verification"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "403": {
                        "description": "Custom domains require a paid subscription"
                    }
                }
            },
            "delete": {
                "summary": "Remove custom domain",
                "description": "Remove your custom domain",
                "tags": ["Domain"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Custom domain removed"
                    },
                    "400": {
                        "description": "No custom domain configured"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/domain/custom/verify": {
            "post": {
                "summary": "Verify custom domain",
                "description": "Verify DNS records for custom domain",
                "tags": ["Domain"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Verification status returned"
                    },
                    "400": {
                        "description": "No custom domain configured"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v2/portfolio": {
            "get": {
                "summary": "Get my portfolio",
                "description": "Get the authenticated user's portfolio",
                "tags": ["Portfolio"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Portfolio retrieved successfully"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Portfolio not found"
                    }
                }
            },
            "post": {
                "summary": "Create portfolio",
                "description": "Create a new portfolio for the authenticated user",
                "tags": ["Portfolio"],
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
                            "required": ["title", "slug"],
                            "properties": {
                                "title": {
                                    "type": "string",
                                    "example": "John Doe's Portfolio"
                                },
                                "slug": {
                                    "type": "string",
                                    "example": "john-doe",
                                    "description": "URL-friendly identifier (3-50 chars, lowercase, alphanumeric with hyphens)"
                                },
                                "tagline": {
                                    "type": "string",
                                    "example": "Full-stack Developer"
                                },
                                "bio": {
                                    "type": "string",
                                    "example": "Passionate developer with 5 years of experience..."
                                },
                                "template": {
                                    "type": "string",
                                    "example": "default",
                                    "description": "Portfolio template"
                                },
                                "theme": {
                                    "type": "object",
                                    "properties": {
                                        "primary_color": {"type": "string", "example": "#3B82F6"},
                                        "secondary_color": {"type": "string", "example": "#1E40AF"},
                                        "accent_color": {"type": "string", "example": "#F59E0B"},
                                        "text_color": {"type": "string", "example": "#1F2937"},
                                        "background_color": {"type": "string", "example": "#FFFFFF"},
                                        "font_family": {"type": "string", "example": "Inter"},
                                        "font_size": {"type": "string", "example": "16px"}
                                    }
                                },
                                "seo": {
                                    "type": "object",
                                    "properties": {
                                        "meta_title": {"type": "string"},
                                        "meta_description": {"type": "string"},
                                        "meta_keywords": {"type": "string"},
                                        "og_image": {"type": "string"}
                                    }
                                },
                                "settings": {
                                    "type": "object",
                                    "properties": {
                                        "show_projects": {"type": "boolean", "default": true},
                                        "show_experiences": {"type": "boolean", "default": true},
                                        "show_education": {"type": "boolean", "default": true},
                                        "show_skills": {"type": "boolean", "default": true},
                                        "show_certifications": {"type": "boolean", "default": true},
                                        "show_contact": {"type": "boolean", "default": true},
                                        "show_social_links": {"type": "boolean", "default": true},
                                        "enable_analytics": {"type": "boolean", "default": false},
                                        "enable_comments": {"type": "boolean", "default": false}
                                    }
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Portfolio created successfully"
                    },
                    "400": {
                        "description": "Portfolio already exists or invalid data"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            },
            "put": {
                "summary": "Update portfolio",
                "description": "Update the authenticated user's portfolio",
                "tags": ["Portfolio"],
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
                            "properties": {
                                "title": {"type": "string"},
                                "slug": {"type": "string"},
                                "tagline": {"type": "string"},
                                "bio": {"type": "string"},
                                "cover_image": {"type": "string"},
                                "logo": {"type": "string"},
                                "template": {"type": "string"},
                                "theme": {"type": "object"},
                                "custom_css": {"type": "string"},
                                "status": {"type": "string", "enum": ["draft", "published", "archived"]},
                                "is_public": {"type": "boolean"},
                                "seo": {"type": "object"},
                                "settings": {"type": "object"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Portfolio updated successfully"
                    },
                    "400": {
                        "description": "Invalid data"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Portfolio not found"
                    }
                }
            },
            "delete": {
                "summary": "Delete portfolio",
                "description": "Delete the authenticated user's portfolio",
                "tags": ["Portfolio"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Portfolio deleted successfully"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Portfolio not found"
                    }
                }
            }
        },
        "/api/v2/portfolio/publish": {
            "post": {
                "summary": "Publish portfolio",
                "description": "Publish the portfolio to make it publicly accessible",
                "tags": ["Portfolio"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Portfolio published successfully"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Portfolio not found"
                    }
                }
            }
        },
        "/api/v2/portfolio/unpublish": {
            "post": {
                "summary": "Unpublish portfolio",
                "description": "Unpublish the portfolio (set to draft)",
                "tags": ["Portfolio"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Portfolio unpublished successfully"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Portfolio not found"
                    }
                }
            }
        },
        "/api/v2/portfolio/sections": {
            "post": {
                "summary": "Create portfolio section",
                "description": "Create a new section in the portfolio",
                "tags": ["Portfolio Sections"],
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
                            "required": ["title", "type"],
                            "properties": {
                                "title": {
                                    "type": "string",
                                    "example": "About Me"
                                },
                                "type": {
                                    "type": "string",
                                    "enum": ["hero", "about", "projects", "experience", "skills", "contact", "custom"],
                                    "example": "about"
                                },
                                "content": {
                                    "type": "string",
                                    "description": "Section content (HTML or markdown)"
                                },
                                "settings": {
                                    "type": "string",
                                    "description": "JSON string for section-specific settings"
                                },
                                "sort_order": {
                                    "type": "integer",
                                    "example": 1
                                },
                                "is_visible": {
                                    "type": "boolean",
                                    "default": true
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Section created successfully"
                    },
                    "400": {
                        "description": "Invalid data"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Portfolio not found"
                    }
                }
            }
        },
        "/api/v2/portfolio/sections/{sectionId}": {
            "put": {
                "summary": "Update portfolio section",
                "description": "Update a section in the portfolio",
                "tags": ["Portfolio Sections"],
                "security": [{"Bearer": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "sectionId",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Section UUID"
                    },
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "title": {"type": "string"},
                                "type": {"type": "string"},
                                "content": {"type": "string"},
                                "settings": {"type": "string"},
                                "sort_order": {"type": "integer"},
                                "is_visible": {"type": "boolean"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Section updated successfully"
                    },
                    "400": {
                        "description": "Invalid data"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Section not found"
                    }
                }
            },
            "delete": {
                "summary": "Delete portfolio section",
                "description": "Delete a section from the portfolio",
                "tags": ["Portfolio Sections"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {
                        "name": "sectionId",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Section UUID"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Section deleted successfully"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Section not found"
                    }
                }
            }
        },
        "/api/v2/portfolio/sections/reorder": {
            "post": {
                "summary": "Reorder portfolio sections",
                "description": "Reorder sections by providing an array of section IDs in desired order",
                "tags": ["Portfolio Sections"],
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
                            "required": ["section_ids"],
                            "properties": {
                                "section_ids": {
                                    "type": "array",
                                    "items": {"type": "string"},
                                    "description": "Array of section UUIDs in desired order"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Sections reordered successfully"
                    },
                    "400": {
                        "description": "Invalid data"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Portfolio not found"
                    }
                }
            }
        },
        "/api/v2/portfolios/{slug}": {
            "get": {
                "summary": "Get public portfolio",
                "description": "Get a published portfolio by its slug (public access)",
                "tags": ["Portfolio"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "slug",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": "Portfolio slug"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Portfolio retrieved successfully"
                    },
                    "404": {
                        "description": "Portfolio not found or not published"
                    }
                }
            }
        },
        "/api/v2/analytics/track/page-view": {
            "post": {
                "summary": "Track page view",
                "description": "Track a page view event",
                "tags": ["Analytics - Tracking"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["path", "session_id"],
                            "properties": {
                                "path": {"type": "string", "example": "/jobs"},
                                "session_id": {"type": "string", "example": "sess_abc123"},
                                "referrer": {"type": "string"},
                                "duration": {"type": "integer", "description": "Time spent in seconds"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "Page view tracked"}
                }
            }
        },
        "/api/v2/analytics/track/job-view": {
            "post": {
                "summary": "Track job view",
                "description": "Track when a job listing is viewed",
                "tags": ["Analytics - Tracking"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["job_id", "session_id"],
                            "properties": {
                                "job_id": {"type": "string", "format": "uuid"},
                                "session_id": {"type": "string"},
                                "referrer": {"type": "string"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "Job view tracked"}
                }
            }
        },
        "/api/v2/analytics/track/profile-view": {
            "post": {
                "summary": "Track profile view",
                "description": "Track when a user profile is viewed",
                "tags": ["Analytics - Tracking"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["profile_user_id", "session_id"],
                            "properties": {
                                "profile_user_id": {"type": "string", "format": "uuid"},
                                "session_id": {"type": "string"},
                                "referrer": {"type": "string"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "Profile view tracked"}
                }
            }
        },
        "/api/v2/analytics/track/portfolio-view": {
            "post": {
                "summary": "Track portfolio view",
                "description": "Track when a portfolio is viewed",
                "tags": ["Analytics - Tracking"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["portfolio_id", "session_id"],
                            "properties": {
                                "portfolio_id": {"type": "string", "format": "uuid"},
                                "session_id": {"type": "string"},
                                "referrer": {"type": "string"},
                                "duration": {"type": "integer"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "Portfolio view tracked"}
                }
            }
        },
        "/api/v2/analytics/track/event": {
            "post": {
                "summary": "Track custom event",
                "description": "Track a custom analytics event",
                "tags": ["Analytics - Tracking"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["event_type", "session_id"],
                            "properties": {
                                "event_type": {"type": "string", "example": "button_click"},
                                "entity_id": {"type": "string", "format": "uuid"},
                                "entity_type": {"type": "string", "example": "job"},
                                "session_id": {"type": "string"},
                                "properties": {"type": "object"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "Event tracked"}
                }
            }
        },
        "/api/v2/analytics/admin/dashboard": {
            "get": {
                "summary": "Admin analytics dashboard",
                "description": "Get comprehensive platform analytics (admin only)",
                "tags": ["Analytics - Admin"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "start_date", "in": "query", "type": "string", "format": "date", "description": "Start date (YYYY-MM-DD)"},
                    {"name": "end_date", "in": "query", "type": "string", "format": "date", "description": "End date (YYYY-MM-DD)"},
                    {"name": "group_by", "in": "query", "type": "string", "enum": ["day", "week", "month"], "description": "Group trend data by"}
                ],
                "responses": {
                    "200": {"description": "Admin dashboard analytics"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Admin access required"}
                }
            }
        },
        "/api/v2/analytics/admin/overview": {
            "get": {
                "summary": "Platform overview",
                "description": "Get platform overview statistics (admin only)",
                "tags": ["Analytics - Admin"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {"description": "Platform overview stats"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Admin access required"}
                }
            }
        },
        "/api/v2/analytics/admin/users": {
            "get": {
                "summary": "User analytics",
                "description": "Get user statistics (admin only)",
                "tags": ["Analytics - Admin"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "start_date", "in": "query", "type": "string", "format": "date"},
                    {"name": "end_date", "in": "query", "type": "string", "format": "date"}
                ],
                "responses": {
                    "200": {"description": "User analytics"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Admin access required"}
                }
            }
        },
        "/api/v2/analytics/admin/jobs": {
            "get": {
                "summary": "Job analytics",
                "description": "Get job statistics (admin only)",
                "tags": ["Analytics - Admin"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "start_date", "in": "query", "type": "string", "format": "date"},
                    {"name": "end_date", "in": "query", "type": "string", "format": "date"}
                ],
                "responses": {
                    "200": {"description": "Job analytics"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Admin access required"}
                }
            }
        },
        "/api/v2/analytics/admin/applications": {
            "get": {
                "summary": "Application analytics",
                "description": "Get application statistics (admin only)",
                "tags": ["Analytics - Admin"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "start_date", "in": "query", "type": "string", "format": "date"},
                    {"name": "end_date", "in": "query", "type": "string", "format": "date"}
                ],
                "responses": {
                    "200": {"description": "Application analytics"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Admin access required"}
                }
            }
        },
        "/api/v2/analytics/admin/revenue": {
            "get": {
                "summary": "Revenue analytics",
                "description": "Get revenue statistics (admin only)",
                "tags": ["Analytics - Admin"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "start_date", "in": "query", "type": "string", "format": "date"},
                    {"name": "end_date", "in": "query", "type": "string", "format": "date"}
                ],
                "responses": {
                    "200": {"description": "Revenue analytics"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Admin access required"}
                }
            }
        },
        "/api/v2/analytics/recruiter/dashboard": {
            "get": {
                "summary": "Recruiter analytics dashboard",
                "description": "Get analytics for recruiter's jobs and applications (recruiter only)",
                "tags": ["Analytics - Recruiter"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "start_date", "in": "query", "type": "string", "format": "date"},
                    {"name": "end_date", "in": "query", "type": "string", "format": "date"},
                    {"name": "group_by", "in": "query", "type": "string", "enum": ["day", "week", "month"]}
                ],
                "responses": {
                    "200": {"description": "Recruiter dashboard analytics"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Recruiter access required"}
                }
            }
        },
        "/api/v2/analytics/recruiter/jobs": {
            "get": {
                "summary": "Recruiter job performance",
                "description": "Get performance metrics for recruiter's jobs",
                "tags": ["Analytics - Recruiter"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "start_date", "in": "query", "type": "string", "format": "date"},
                    {"name": "end_date", "in": "query", "type": "string", "format": "date"}
                ],
                "responses": {
                    "200": {"description": "Job performance metrics"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Recruiter access required"}
                }
            }
        },
        "/api/v2/analytics/recruiter/applications": {
            "get": {
                "summary": "Recruiter application stats",
                "description": "Get application statistics for recruiter's jobs",
                "tags": ["Analytics - Recruiter"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {"description": "Application statistics"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Recruiter access required"}
                }
            }
        },
        "/api/v2/analytics/talent/dashboard": {
            "get": {
                "summary": "Talent analytics dashboard",
                "description": "Get analytics for talent's profile, portfolio, and applications",
                "tags": ["Analytics - Talent"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "start_date", "in": "query", "type": "string", "format": "date"},
                    {"name": "end_date", "in": "query", "type": "string", "format": "date"},
                    {"name": "group_by", "in": "query", "type": "string", "enum": ["day", "week", "month"]}
                ],
                "responses": {
                    "200": {"description": "Talent dashboard analytics"},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/api/v2/analytics/talent/profile-views": {
            "get": {
                "summary": "Profile views analytics",
                "description": "Get profile views statistics for the authenticated user",
                "tags": ["Analytics - Talent"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "start_date", "in": "query", "type": "string", "format": "date"},
                    {"name": "end_date", "in": "query", "type": "string", "format": "date"}
                ],
                "responses": {
                    "200": {"description": "Profile views analytics"},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/api/v2/analytics/talent/portfolio": {
            "get": {
                "summary": "Portfolio analytics",
                "description": "Get portfolio views and engagement statistics",
                "tags": ["Analytics - Talent"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "start_date", "in": "query", "type": "string", "format": "date"},
                    {"name": "end_date", "in": "query", "type": "string", "format": "date"}
                ],
                "responses": {
                    "200": {"description": "Portfolio analytics"},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/api/v2/analytics/talent/applications": {
            "get": {
                "summary": "Talent application stats",
                "description": "Get statistics for user's job applications",
                "tags": ["Analytics - Talent"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {"description": "Application statistics"},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/api/v2/analytics/talent/viewer-insights": {
            "get": {
                "summary": "Viewer insights",
                "description": "Get insights about who viewed your profile (recruiters vs talents)",
                "tags": ["Analytics - Talent"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "start_date", "in": "query", "type": "string", "format": "date"},
                    {"name": "end_date", "in": "query", "type": "string", "format": "date"}
                ],
                "responses": {
                    "200": {"description": "Viewer insights"},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/api/v2/notification-settings": {
            "get": {
                "summary": "Get notification settings",
                "description": "Get current user's notification preferences",
                "tags": ["Notification Settings"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Notification settings retrieved",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "id": {"type": "string", "format": "uuid"},
                                "user_id": {"type": "string", "format": "uuid"},
                                "email": {
                                    "type": "object",
                                    "properties": {
                                        "app_updates": {"type": "boolean"},
                                        "new_messages": {"type": "boolean"},
                                        "job_recommendations": {"type": "boolean"},
                                        "newsletter": {"type": "boolean"},
                                        "marketing_emails": {"type": "boolean"}
                                    }
                                },
                                "push": {
                                    "type": "object",
                                    "properties": {
                                        "app_updates": {"type": "boolean"},
                                        "new_messages": {"type": "boolean"},
                                        "reminders": {"type": "boolean"}
                                    }
                                },
                                "in_app": {
                                    "type": "object",
                                    "properties": {
                                        "activity_updates": {"type": "boolean"},
                                        "mentions": {"type": "boolean"},
                                        "announcements": {"type": "boolean"}
                                    }
                                }
                            }
                        }
                    },
                    "401": {"description": "Unauthorized"}
                }
            },
            "put": {
                "summary": "Update all notification settings",
                "description": "Update all notification preferences at once",
                "tags": ["Notification Settings"],
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
                            "properties": {
                                "email": {
                                    "type": "object",
                                    "properties": {
                                        "app_updates": {"type": "boolean"},
                                        "new_messages": {"type": "boolean"},
                                        "job_recommendations": {"type": "boolean"},
                                        "newsletter": {"type": "boolean"},
                                        "marketing_emails": {"type": "boolean"}
                                    }
                                },
                                "push": {
                                    "type": "object",
                                    "properties": {
                                        "app_updates": {"type": "boolean"},
                                        "new_messages": {"type": "boolean"},
                                        "reminders": {"type": "boolean"}
                                    }
                                },
                                "in_app": {
                                    "type": "object",
                                    "properties": {
                                        "activity_updates": {"type": "boolean"},
                                        "mentions": {"type": "boolean"},
                                        "announcements": {"type": "boolean"}
                                    }
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "Settings updated successfully"},
                    "400": {"description": "Invalid data"},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/api/v2/notification-settings/email": {
            "put": {
                "summary": "Update email notification settings",
                "description": "Update email notification preferences only",
                "tags": ["Notification Settings"],
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
                            "properties": {
                                "app_updates": {"type": "boolean"},
                                "new_messages": {"type": "boolean"},
                                "job_recommendations": {"type": "boolean"},
                                "newsletter": {"type": "boolean"},
                                "marketing_emails": {"type": "boolean"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "Email settings updated"},
                    "400": {"description": "Invalid data"},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/api/v2/notification-settings/push": {
            "put": {
                "summary": "Update push notification settings",
                "description": "Update push notification preferences only",
                "tags": ["Notification Settings"],
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
                            "properties": {
                                "app_updates": {"type": "boolean"},
                                "new_messages": {"type": "boolean"},
                                "reminders": {"type": "boolean"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "Push settings updated"},
                    "400": {"description": "Invalid data"},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/api/v2/notification-settings/in-app": {
            "put": {
                "summary": "Update in-app notification settings",
                "description": "Update in-app notification preferences only",
                "tags": ["Notification Settings"],
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
                            "properties": {
                                "activity_updates": {"type": "boolean"},
                                "mentions": {"type": "boolean"},
                                "announcements": {"type": "boolean"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "In-app settings updated"},
                    "400": {"description": "Invalid data"},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/api/v2/announcements": {
            "get": {
                "summary": "Get announcements",
                "description": "Get announcements for the current user based on their role",
                "tags": ["Announcements"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "limit", "in": "query", "type": "integer", "default": 10}
                ],
                "responses": {
                    "200": {
                        "description": "Announcements retrieved",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {
                                    "type": "array",
                                    "items": {
                                        "type": "object",
                                        "properties": {
                                            "id": {"type": "string", "format": "uuid"},
                                            "title": {"type": "string"},
                                            "content": {"type": "string"},
                                            "type": {"type": "string", "enum": ["ANNOUNCEMENT", "APP_UPDATE", "SYSTEM_ALERT", "MAINTENANCE"]},
                                            "priority": {"type": "string", "enum": ["LOW", "NORMAL", "HIGH", "CRITICAL"]},
                                            "show_as_banner": {"type": "boolean"},
                                            "action_url": {"type": "string"},
                                            "action_text": {"type": "string"},
                                            "is_read": {"type": "boolean"},
                                            "is_dismissed": {"type": "boolean"},
                                            "published_at": {"type": "string", "format": "date-time"}
                                        }
                                    }
                                },
                                "total_items": {"type": "integer"},
                                "total_pages": {"type": "integer"},
                                "page": {"type": "integer"},
                                "limit": {"type": "integer"}
                            }
                        }
                    },
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/api/v2/announcements/banners": {
            "get": {
                "summary": "Get active banners",
                "description": "Get active announcement banners for display",
                "tags": ["Announcements"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Active banners retrieved",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "properties": {
                                    "id": {"type": "string", "format": "uuid"},
                                    "title": {"type": "string"},
                                    "content": {"type": "string"},
                                    "type": {"type": "string"},
                                    "priority": {"type": "string"},
                                    "banner_color": {"type": "string"},
                                    "action_url": {"type": "string"},
                                    "action_text": {"type": "string"}
                                }
                            }
                        }
                    },
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/api/v2/announcements/{id}/read": {
            "put": {
                "summary": "Mark announcement as read",
                "description": "Mark an announcement as read for the current user",
                "tags": ["Announcements"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "description": "Announcement UUID"}
                ],
                "responses": {
                    "200": {"description": "Announcement marked as read"},
                    "401": {"description": "Unauthorized"},
                    "404": {"description": "Announcement not found"}
                }
            }
        },
        "/api/v2/announcements/{id}/dismiss": {
            "put": {
                "summary": "Dismiss announcement",
                "description": "Dismiss an announcement for the current user",
                "tags": ["Announcements"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "description": "Announcement UUID"}
                ],
                "responses": {
                    "200": {"description": "Announcement dismissed"},
                    "401": {"description": "Unauthorized"},
                    "404": {"description": "Announcement not found"}
                }
            }
        },
        "/api/v2/admin/announcements": {
            "get": {
                "summary": "List all announcements (Admin)",
                "description": "Get all announcements with filtering (admin only)",
                "tags": ["Announcements - Admin"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "limit", "in": "query", "type": "integer", "default": 10},
                    {"name": "type", "in": "query", "type": "string", "enum": ["ANNOUNCEMENT", "APP_UPDATE", "SYSTEM_ALERT", "MAINTENANCE"]},
                    {"name": "is_published", "in": "query", "type": "boolean"}
                ],
                "responses": {
                    "200": {"description": "Announcements retrieved"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Admin access required"}
                }
            },
            "post": {
                "summary": "Create announcement (Admin)",
                "description": "Create a new announcement (admin only)",
                "tags": ["Announcements - Admin"],
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
                            "required": ["title", "content", "type"],
                            "properties": {
                                "title": {"type": "string", "example": "New Feature Released"},
                                "content": {"type": "string", "example": "We have released a new messaging feature!"},
                                "type": {"type": "string", "enum": ["ANNOUNCEMENT", "APP_UPDATE", "SYSTEM_ALERT", "MAINTENANCE"]},
                                "target_audience": {"type": "string", "enum": ["ALL_USERS", "ADMINS_ONLY", "RECRUITERS_ONLY", "PREMIUM_ONLY"], "default": "ALL_USERS"},
                                "priority": {"type": "string", "enum": ["LOW", "NORMAL", "HIGH", "CRITICAL"], "default": "NORMAL"},
                                "show_as_banner": {"type": "boolean", "default": false},
                                "banner_color": {"type": "string", "example": "#3B82F6"},
                                "action_url": {"type": "string", "example": "https://example.com/feature"},
                                "action_text": {"type": "string", "example": "Learn More"},
                                "scheduled_at": {"type": "string", "format": "date-time"},
                                "expires_at": {"type": "string", "format": "date-time"}
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {"description": "Announcement created"},
                    "400": {"description": "Invalid data"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Admin access required"}
                }
            }
        },
        "/api/v2/admin/announcements/{id}": {
            "get": {
                "summary": "Get announcement (Admin)",
                "description": "Get a single announcement by ID (admin only)",
                "tags": ["Announcements - Admin"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "description": "Announcement UUID"}
                ],
                "responses": {
                    "200": {"description": "Announcement retrieved"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Admin access required"},
                    "404": {"description": "Announcement not found"}
                }
            },
            "put": {
                "summary": "Update announcement (Admin)",
                "description": "Update an announcement (admin only)",
                "tags": ["Announcements - Admin"],
                "security": [{"Bearer": []}],
                "consumes": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "description": "Announcement UUID"},
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "title": {"type": "string"},
                                "content": {"type": "string"},
                                "type": {"type": "string", "enum": ["ANNOUNCEMENT", "APP_UPDATE", "SYSTEM_ALERT", "MAINTENANCE"]},
                                "target_audience": {"type": "string", "enum": ["ALL_USERS", "ADMINS_ONLY", "RECRUITERS_ONLY", "PREMIUM_ONLY"]},
                                "priority": {"type": "string", "enum": ["LOW", "NORMAL", "HIGH", "CRITICAL"]},
                                "show_as_banner": {"type": "boolean"},
                                "banner_color": {"type": "string"},
                                "action_url": {"type": "string"},
                                "action_text": {"type": "string"},
                                "scheduled_at": {"type": "string", "format": "date-time"},
                                "expires_at": {"type": "string", "format": "date-time"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "Announcement updated"},
                    "400": {"description": "Invalid data"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Admin access required"},
                    "404": {"description": "Announcement not found"}
                }
            },
            "delete": {
                "summary": "Delete announcement (Admin)",
                "description": "Delete an announcement (admin only)",
                "tags": ["Announcements - Admin"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "description": "Announcement UUID"}
                ],
                "responses": {
                    "200": {"description": "Announcement deleted"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Admin access required"},
                    "404": {"description": "Announcement not found"}
                }
            }
        },
        "/api/v2/admin/announcements/{id}/publish": {
            "put": {
                "summary": "Publish announcement (Admin)",
                "description": "Publish an announcement (admin only)",
                "tags": ["Announcements - Admin"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "description": "Announcement UUID"}
                ],
                "responses": {
                    "200": {"description": "Announcement published"},
                    "400": {"description": "Announcement already published"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Admin access required"},
                    "404": {"description": "Announcement not found"}
                }
            }
        },
        "/api/v2/admin/announcements/{id}/unpublish": {
            "put": {
                "summary": "Unpublish announcement (Admin)",
                "description": "Unpublish an announcement (admin only)",
                "tags": ["Announcements - Admin"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "description": "Announcement UUID"}
                ],
                "responses": {
                    "200": {"description": "Announcement unpublished"},
                    "400": {"description": "Announcement not published"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Admin access required"},
                    "404": {"description": "Announcement not found"}
                }
            }
        },
        "/api/v2/chat/conversations": {
            "get": {
                "summary": "Get conversations",
                "description": "Get all conversations for the current user",
                "tags": ["Chat"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "limit", "in": "query", "type": "integer", "default": 20}
                ],
                "responses": {
                    "200": {
                        "description": "Conversations retrieved",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {
                                    "type": "array",
                                    "items": {
                                        "type": "object",
                                        "properties": {
                                            "id": {"type": "string", "format": "uuid"},
                                            "other_user": {
                                                "type": "object",
                                                "properties": {
                                                    "id": {"type": "string", "format": "uuid"},
                                                    "name": {"type": "string"},
                                                    "username": {"type": "string"},
                                                    "image": {"type": "string"}
                                                }
                                            },
                                            "last_message": {
                                                "type": "object",
                                                "properties": {
                                                    "id": {"type": "string", "format": "uuid"},
                                                    "content": {"type": "string"},
                                                    "sender_id": {"type": "string", "format": "uuid"},
                                                    "created_at": {"type": "string", "format": "date-time"}
                                                }
                                            },
                                            "unread_count": {"type": "integer"},
                                            "created_at": {"type": "string", "format": "date-time"},
                                            "updated_at": {"type": "string", "format": "date-time"}
                                        }
                                    }
                                },
                                "total_items": {"type": "integer"},
                                "total_pages": {"type": "integer"},
                                "page": {"type": "integer"},
                                "limit": {"type": "integer"}
                            }
                        }
                    },
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/api/v2/chat/conversations/{id}": {
            "get": {
                "summary": "Get conversation",
                "description": "Get a single conversation by ID",
                "tags": ["Chat"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "description": "Conversation UUID"}
                ],
                "responses": {
                    "200": {"description": "Conversation retrieved"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Not a participant"},
                    "404": {"description": "Conversation not found"}
                }
            },
            "delete": {
                "summary": "Delete conversation",
                "description": "Delete a conversation",
                "tags": ["Chat"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "description": "Conversation UUID"}
                ],
                "responses": {
                    "200": {"description": "Conversation deleted"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Not a participant"},
                    "404": {"description": "Conversation not found"}
                }
            }
        },
        "/api/v2/chat/conversations/user/{userId}": {
            "get": {
                "summary": "Get or create conversation",
                "description": "Get an existing conversation with a user or create a new one",
                "tags": ["Chat"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "userId", "in": "path", "required": true, "type": "string", "description": "Other user's UUID"}
                ],
                "responses": {
                    "200": {"description": "Conversation retrieved or created"},
                    "400": {"description": "Cannot message self"},
                    "401": {"description": "Unauthorized"},
                    "404": {"description": "User not found"}
                }
            }
        },
        "/api/v2/chat/messages": {
            "post": {
                "summary": "Send message",
                "description": "Send a direct message to another user. Can include text content, media attachments, or both. Messages can also be sent via WebSocket using the action 'send_message'.",
                "tags": ["Chat"],
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
                            "required": ["recipient_id"],
                            "properties": {
                                "recipient_id": {"type": "string", "format": "uuid", "description": "Recipient's user UUID"},
                                "content": {"type": "string", "description": "Message content (optional if media is provided)", "example": "Hello, how are you?"},
                                "media": {
                                    "type": "array",
                                    "description": "Media attachments (max 10)",
                                    "maxItems": 10,
                                    "items": {
                                        "type": "object",
                                        "required": ["type", "url"],
                                        "properties": {
                                            "type": {"type": "string", "enum": ["IMAGE", "VIDEO", "AUDIO", "DOCUMENT", "FILE"], "description": "Media type"},
                                            "url": {"type": "string", "format": "url", "description": "URL of the uploaded media"},
                                            "file_name": {"type": "string", "description": "Original file name"},
                                            "file_size": {"type": "integer", "description": "File size in bytes"},
                                            "mime_type": {"type": "string", "description": "MIME type", "example": "image/jpeg"},
                                            "width": {"type": "integer", "description": "Width in pixels (for images/videos)"},
                                            "height": {"type": "integer", "description": "Height in pixels (for images/videos)"},
                                            "duration": {"type": "integer", "description": "Duration in seconds (for audio/video)"},
                                            "thumbnail": {"type": "string", "description": "Thumbnail URL (for videos)"}
                                        }
                                    }
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Message sent",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "id": {"type": "string", "format": "uuid"},
                                "conversation_id": {"type": "string", "format": "uuid"},
                                "sender_id": {"type": "string", "format": "uuid"},
                                "recipient_id": {"type": "string", "format": "uuid"},
                                "content": {"type": "string"},
                                "media": {
                                    "type": "array",
                                    "items": {
                                        "type": "object",
                                        "properties": {
                                            "id": {"type": "string"},
                                            "type": {"type": "string"},
                                            "url": {"type": "string"},
                                            "file_name": {"type": "string"},
                                            "file_size": {"type": "integer"},
                                            "mime_type": {"type": "string"},
                                            "width": {"type": "integer"},
                                            "height": {"type": "integer"},
                                            "duration": {"type": "integer"},
                                            "thumbnail": {"type": "string"}
                                        }
                                    }
                                },
                                "status": {"type": "string", "enum": ["SENT", "DELIVERED", "READ"]},
                                "created_at": {"type": "string", "format": "date-time"}
                            }
                        }
                    },
                    "400": {"description": "Cannot message self, message must have content or media"},
                    "401": {"description": "Unauthorized"},
                    "404": {"description": "Recipient not found"}
                }
            }
        },
        "/api/v2/chat/conversations/{id}/messages": {
            "get": {
                "summary": "Get messages",
                "description": "Get messages in a conversation",
                "tags": ["Chat"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "description": "Conversation UUID"},
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "limit", "in": "query", "type": "integer", "default": 50}
                ],
                "responses": {
                    "200": {
                        "description": "Messages retrieved",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {
                                    "type": "array",
                                    "items": {
                                        "type": "object",
                                        "properties": {
                                            "id": {"type": "string", "format": "uuid"},
                                            "conversation_id": {"type": "string", "format": "uuid"},
                                            "sender_id": {"type": "string", "format": "uuid"},
                                            "recipient_id": {"type": "string", "format": "uuid"},
                                            "content": {"type": "string"},
                                            "media": {
                                                "type": "array",
                                                "items": {
                                                    "type": "object",
                                                    "properties": {
                                                        "id": {"type": "string"},
                                                        "type": {"type": "string", "enum": ["IMAGE", "VIDEO", "AUDIO", "DOCUMENT", "FILE"]},
                                                        "url": {"type": "string"},
                                                        "file_name": {"type": "string"},
                                                        "file_size": {"type": "integer"},
                                                        "mime_type": {"type": "string"},
                                                        "width": {"type": "integer"},
                                                        "height": {"type": "integer"},
                                                        "duration": {"type": "integer"},
                                                        "thumbnail": {"type": "string"}
                                                    }
                                                }
                                            },
                                            "status": {"type": "string", "enum": ["SENT", "DELIVERED", "READ"]},
                                            "read_at": {"type": "string", "format": "date-time"},
                                            "created_at": {"type": "string", "format": "date-time"},
                                            "sender": {
                                                "type": "object",
                                                "properties": {
                                                    "id": {"type": "string", "format": "uuid"},
                                                    "name": {"type": "string"},
                                                    "username": {"type": "string"},
                                                    "image": {"type": "string"}
                                                }
                                            }
                                        }
                                    }
                                },
                                "total_items": {"type": "integer"},
                                "total_pages": {"type": "integer"},
                                "page": {"type": "integer"},
                                "limit": {"type": "integer"}
                            }
                        }
                    },
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Not a participant"},
                    "404": {"description": "Conversation not found"}
                }
            }
        },
        "/api/v2/chat/conversations/{id}/read": {
            "put": {
                "summary": "Mark messages as read",
                "description": "Mark all messages in a conversation as read",
                "tags": ["Chat"],
                "security": [{"Bearer": []}],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "description": "Conversation UUID"}
                ],
                "responses": {
                    "200": {"description": "Messages marked as read"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Not a participant"},
                    "404": {"description": "Conversation not found"}
                }
            }
        },
        "/api/v2/chat/unread": {
            "get": {
                "summary": "Get unread count",
                "description": "Get total unread message count for the current user",
                "tags": ["Chat"],
                "security": [{"Bearer": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Unread count retrieved",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "unread_count": {"type": "integer"}
                            }
                        }
                    },
                    "401": {"description": "Unauthorized"}
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
