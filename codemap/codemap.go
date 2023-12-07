package codemap

import (
	"log"

	"github.com/Clement-Jean/protein/token"
)

type CodeMap struct {
	files map[string]FileMap
}

func New() CodeMap {
	return CodeMap{files: make(map[string]FileMap)}
}

func (cm *CodeMap) Remove(fileName string) {
	if cm.files == nil {
		return
	}
	delete(cm.files, fileName)
}

func (cm *CodeMap) Insert(fileName string, content []byte) *FileMap {
	if cm.files == nil {
		return nil
	}
	fm := FileMap{content: content}
	cm.files[fileName] = fm
	return &fm
}

func (cm *CodeMap) Lookup(id token.UniqueID) []byte {
	if cm.files == nil {
		return nil
	}
	for _, file := range cm.files {
		if slice := file.Lookup(id); slice != nil {
			return slice
		}
	}

	log.Panicf("%d wasn't in any of the FileMaps", id)
	return nil
}
