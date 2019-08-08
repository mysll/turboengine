package proto

type Echo struct {
	Ver string `version:"1.0.0"`
	XXX interface{}
	// custom method begin
	Print func(string) error
	Echo  func(string) (string, error)
	// custom method end
}

func init() {
	reg["Echo"] = new(Echo)
}
