package domain

import "fmt"

type UnauthorizedError struct {
	Message string
}

func (e UnauthorizedError) Error() string {
	return fmt.Sprintf("Unauthorized: %s", e.Message)
}

type UnauthenticatedError struct {
	Message string
}

func (e UnauthenticatedError) Error() string {
	return fmt.Sprintf("Unauthenticated: %s", e.Message)
}

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("Not Found: %s", e.Message)
}

type BadRequestError struct {
	Message string
}

func (e BadRequestError) Error() string {
	return fmt.Sprintf("Bad Request: %s", e.Message)
}

type MissingValuesError struct {
	Message string
}

func (e MissingValuesError) Error() string {
	return fmt.Sprintf("Missing Values: %s", e.Message)
}

type ConflictError struct {
	Message string
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("Conflict: %s", e.Message)
}

type InternalServerError struct {
	Message string
}

func (e InternalServerError) Error() string {
	return fmt.Sprintf("Internal Server Error: %s", e.Message)
}

type InvalidInputError struct {
	Message string
}

func (e InvalidInputError) Error() string {
	return fmt.Sprintf("Invalid Input: %s", e.Message)
}
