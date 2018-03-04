package error


func Fatal(err error) {
	if err != nil {
		panic(err)
	}
}
