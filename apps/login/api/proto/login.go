package proto

type Login struct {
	Ver string `version:"1.0.0"`
	XXX any
	// custom method begin
	Login func(string, string) (bool, error)
	// custom method end
}

func init() {
	reg["Login"] = new(Login)
}
