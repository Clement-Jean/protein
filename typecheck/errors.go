package typecheck

import (
	"fmt"
	"strings"
)

type Warning interface {
	Warning() string
}

type ImportCycleError struct {
	Files []string
}

func (e *ImportCycleError) Error() string {
	var msg strings.Builder

	for _, file := range e.Files {
		msg.WriteString(file)
		msg.WriteString(" -> ")
	}
	msg.WriteString(e.Files[0])

	return fmt.Sprintf("cycle found: %s", msg.String())
}

type ImportFileNotFoundError struct {
	File         string
	IncludePaths []string
}

func (e *ImportFileNotFoundError) Error() string {
	return fmt.Sprintf("file %s was not found in include paths %v", e.File, e.IncludePaths)
}

type PackageMultipleDefError struct {
	File string
}

func (e *PackageMultipleDefError) Error() string {
	return fmt.Sprintf("multiple package definitions in %s", e.File)
}

type TypeNotDefinedError struct {
	Name string
}

func (e *TypeNotDefinedError) Error() string {
	return fmt.Sprintf("%s is not defined", e.Name)
}

type TypeRedefinedError struct {
	Name string
}

func (e *TypeRedefinedError) Error() string {
	return fmt.Sprintf("%s is redefined", e.Name)
}

type TypeUnusedWarning struct {
	Name string
}

func (w *TypeUnusedWarning) Warning() string {
	return fmt.Sprintf("%s is defined but not used", w.Name)
}

func (w *TypeUnusedWarning) Error() string {
	return w.Warning()
}
