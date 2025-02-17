package linker

type LinkerOpt func(*Linker)

// /!\ this is order sensitive. The paths will be checks in the specific order you provided.
func WithIncludePaths(paths ...string) func(*Linker) {
	return func(l *Linker) {
		l.includePaths = paths
	}
}
