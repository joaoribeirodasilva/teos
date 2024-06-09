package requests

type UserUsersCreate struct {
	FirstName     string `form:"firstName"`
	Surename      string `form:"surename"`
	Email         string `form:"email"`
	Password      string `form:"password"`
	CheckPassword string `form:"checkPassword"`
	Terms         bool   `form:"terms"`
}
