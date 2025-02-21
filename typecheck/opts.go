package typecheck

import (
	"github.com/Clement-Jean/protein/source"
)

type TypeCheckerOpt func(*TypeChecker)

// /!\ this is order sensitive. The paths will be checked in the specific order you provided.
func WithIncludePaths(paths ...string) func(*TypeChecker) {
	return func(tc *TypeChecker) {
		tc.includePaths = paths

		// FIX: implement proper path sanitizer
		//for i := 0; i < len(tc.includePaths); i++ {
		// TODO: check the getwd error
		//	tc.includePaths[i], _ = filepath.Abs(tc.includePaths[i])
		//}
	}
}

type SourceCreator = func(string) (*source.Buffer, error)

func WithSourceCreator(creator SourceCreator) func(*TypeChecker) {
	return func(tc *TypeChecker) {
		tc.srcCreator = creator
	}
}

type FileExistsCheck = func(string) bool

func WithFileCheck(check FileExistsCheck) func(*TypeChecker) {
	return func(tc *TypeChecker) {
		tc.fileCheck = check
	}
}
