package store

// Options type
type Options struct {
	ProjectID      string
	ServiceAccount []byte
}

// Option type
type Option func(*Options)

// ProjectID sets project id to options
func ProjectID(projectID string) Option {
	return func(args *Options) {
		args.ProjectID = projectID
	}
}

// ServiceAccount sets service account to options
func ServiceAccount(serviceAccount []byte) Option {
	return func(args *Options) {
		args.ServiceAccount = serviceAccount
	}
}
