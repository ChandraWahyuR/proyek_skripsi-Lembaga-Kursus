package constant

import "errors"

const Unauthorized = "Unauthorized"
const InternalServerError = "Internal Server Error"
const BadInput = "Format data not valid"

var ErrBadRequest = errors.New("bad request")

// JWT
var ErrGenerateJWT = errors.New("failed to generate jwt token")
var ErrValidateJWT = errors.New("failed to validate jwt token")

// Validator
var ErrHashPassword = errors.New("failed to hash password")

var ErrEmptyEmailRegister = errors.New("Email cannot be empty")
var ErrEmptyPasswordRegister = errors.New("Password cannot be empty")
var ErrEmptyAddressRegister = errors.New("Address cannot be empty")
var ErrEmptyNameRegister = errors.New("Name cannot be empty")
var ErrPasswordNotMatch = errors.New("Password not match")
var ErrInvalidEmail = errors.New("Email is not valid")
var ErrInvalidUsername = errors.New("Username formating not valid")
var ErrInvalidPhone = errors.New("Phone formating not valid")
var ErrEmptyLogin = errors.New("Email or Password cannot be empty")
var UserNotFound = errors.New("User not found")
