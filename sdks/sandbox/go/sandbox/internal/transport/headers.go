package transport

func MergeHeaders(base, endpoint map[string]string) map[string]string {
	merged := make(map[string]string, len(base)+len(endpoint))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range endpoint {
		merged[k] = v
	}
	return merged
}
