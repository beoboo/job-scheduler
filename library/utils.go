package main

func fatalf(format string, args ...interface{}) {
	logger.Fatalf(format+"\n", args...)
}

func warnf(format string, args ...interface{}) {
	logger.Warnf(format+"\n", args...)
}

func infof(format string, args ...interface{}) {
	logger.Infof(format+"\n", args...)
}

func do(val string, err error) string {
	check(err)

	return val
}

func check(err error) {
	if err != nil {
		fatalf("Unexpected: %s", err)
	}
}
