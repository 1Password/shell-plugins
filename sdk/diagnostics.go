package sdk

type Diagnostics struct {
	Errors []Error
}

type Error struct {
	Message string
}
