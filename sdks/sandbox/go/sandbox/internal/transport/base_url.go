package transport

func BaseURL(protocol, domain string) string {
	return protocol + "://" + domain
}
