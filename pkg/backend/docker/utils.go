package docker

func shortName(long string) string {
	return string([]rune(long)[:16])
}
