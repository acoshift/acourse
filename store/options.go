package store

// Options type
type Options struct {
	ProjectID string
	JSONKey   []byte
}

// Option type
type Option func(*Options)

// ProjectID sets project id to options
func ProjectID(projectID string) Option {
	return func(args *Options) {
		args.ProjectID = projectID
	}
}

// JSONKey sets jsonKey to options
func JSONKey(jsonKey []byte) Option {
	return func(args *Options) {
		args.JSONKey = jsonKey
	}
}
