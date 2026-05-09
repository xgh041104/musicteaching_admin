package v1

var (
	ErrSuccess             = newError(200, "success")
	ErrBadRequest          = newError(400, "bad request")
	ErrInvalidBookID       = newError(401, "invalid book ID")
	ErrInvalidCourseID     = newError(402, "invalid course ID")
	ErrUnauthorized        = newError(403, "unauthorized")
	ErrNotFound            = newError(404, "resource not found")
	ErrEmailAlreadyUse     = newError(409, "email already in use")
	ErrInternalServerError = newError(500, "internal server error")
	ErrDeleteCourseFailed  = newError(501, "failed to delete course")
	ErrQueryBooksFailed    = newError(502, "failed to retrieve books")
	ErrQueryCoursesFailed  = newError(503, "failed to retrieve courses")
)
