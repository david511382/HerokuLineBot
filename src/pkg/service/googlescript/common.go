package googlescript

func New(url string) *GoogleScript {
	return &GoogleScript{
		url: url,
	}
}
