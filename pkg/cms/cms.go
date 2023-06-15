package cms

type CMS interface {
	Get(key string) string
	Set(key string, value string) error
}

type InMemCMS map[string]string

func (cms InMemCMS) Get(key string) string              { return cms[key] }
func (cms InMemCMS) Set(key string, value string) error { cms[key] = value; return nil }
