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
	File string
}

func (e *ImportFileNotFoundError) Error() string {
	return fmt.Sprintf("file %s was not found or had errors", e.File)
}

type PackageMultipleDefError struct {
	File string
}

func (e *PackageMultipleDefError) Error() string {
	return fmt.Sprintf("multiple package definitions in %s", e.File)
}

type NotTypeError struct {
	Name      string
	File      string
	Line, Col int
}

func (e *NotTypeError) Error() string {
	return fmt.Sprintf("%s is not a type", e.Name)
}

type NotMessageTypeError struct {
	Name      string
	File      string
	Line, Col int
}

func (e *NotMessageTypeError) Error() string {
	return fmt.Sprintf("%s is not a message type", e.Name)
}

type TypeResolvedNotDefinedError struct {
	Name, ResolvedName string
	File               string
	Line, Col          int
}

func (e *TypeResolvedNotDefinedError) Error() string {
	return fmt.Sprintf("%s is resolved to %s, which is not defined. The innermost scope is searched first in name resolution. Consider using a leading '.' (i.e., %s) to start from the outermost scope", e.Name, e.ResolvedName, "."+e.Name)
}

type TypeNotDefinedError struct {
	Name      string
	File      string
	Line, Col int
}

func (e *TypeNotDefinedError) Error() string {
	return fmt.Sprintf("%s is not defined", e.Name)
}

type TypeNotImportedError struct {
	Name             string
	DefFile, RefFile string
	Line, Col        int
}

func (e *TypeNotImportedError) Error() string {
	return fmt.Sprintf("%s seems to be defined in %s, which is not imported by %s. To use it here, please add the necessary import.", e.Name, e.DefFile, e.RefFile)
}

type TypeRedefinedError struct {
	Name        string
	Files       []string
	Lines, Cols []int
}

func (e *TypeRedefinedError) Error() string {
	return fmt.Sprintf("%s is redefined", e.Name)
}

type TypeUnusedWarning struct {
	Name      string
	File      string
	Line, Col int
}

func (w *TypeUnusedWarning) Warning() string {
	return fmt.Sprintf("%s is defined but not used", w.Name)
}

func (w *TypeUnusedWarning) Error() string {
	return w.Warning()
}

type ImportAlreadyImportedWarning struct {
	ImportingFile, ImportedFile string
	Line, Col                   int
}

func (w *ImportAlreadyImportedWarning) Warning() string {
	return fmt.Sprintf("%s is already imported", w.ImportedFile)
}

func (w *ImportAlreadyImportedWarning) Error() string {
	return w.Warning()
}

type WeakImportNoEffectWarning struct{}

func (w *WeakImportNoEffectWarning) Warning() string {
	return "weak imports have no effect in Protein"
}

func (w *WeakImportNoEffectWarning) Error() string {
	return w.Warning()
}
