package utils

func DefaultIfEmpty(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
