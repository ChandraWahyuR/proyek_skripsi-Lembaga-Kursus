package constant

import "errors"

const Unauthorized = "Unauthorized"
const InternalServerError = "Internal Server Error"
const BadInput = "Format data not valid"

var ErrBadRequest = errors.New("bad request")
var ErrUnauthorized = errors.New("unauthorized")

var ErrEmptyOtp = errors.New("otp cannot be empty")

// JWT
var ErrGenerateJWT = errors.New("failed to generate jwt token")
var ErrValidateJWT = errors.New("failed to validate jwt token")

// Validator
var ErrHashPassword = errors.New("failed to hash password")

// Empty Register
var ErrEmptyEmailRegister = errors.New("email cannot be empty")
var ErrEmptyNameRegister = errors.New("username cannot be empty")
var ErrEmptyPasswordRegister = errors.New("password cannot be empty")
var ErrPasswordNotMatch = errors.New("password not match")

// Register Format Not Valid
var ErrInvalidEmail = errors.New("email is not valid")
var ErrInvalidUsername = errors.New("username formating not valid")
var ErrInvalidPhone = errors.New("phone formating not valid")

// Login
var ErrEmptyLogin = errors.New("email or Password cannot be empty")
var ErrUserNotFound = errors.New("user not found")
var ErrLenPassword = errors.New("password must be at least 8 characters")
var ErrInvalidPassword = errors.New("password must contain at least 1 number, 1 uppercase letter, one punctuation symbol and 1 lowercase letter")
var ErrEmptyPasswordLogin = errors.New("password cannot be empty")
