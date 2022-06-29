package file

import (
	"encoding/json"
	"errors"
	"fmt"
	config "github.com/Kuart/metric-collector/config/server"
	"github.com/Kuart/metric-collector/internal/storage/inmemory"
	"github.com/Kuart/metric-collector/internal/util"
	"log"
	"os"
)

var (
	errNoFile     = errors.New("there is no file in the selected path")
	errOpenFile   = "error opening file"
	errNoFileInfo = errors.New("can't get file info")
	errEmptyFile  = errors.New("the file in the selected path is empty")
)

type Storage struct {
	path       string
	openedFile *os.File
}

func New(cfg config.Config) Storage {
	return Storage{
		path: cfg.StoreFile,
	}
}

func (fs *Storage) GetFileData() (ifs inmemory.FileStorage, err error) {
	file, err := os.OpenFile(fs.path, os.O_RDONLY, 0)

	if err != nil {
		return inmemory.FileStorage{}, fmt.Errorf("%s, Action:%s, File:%s", errOpenFile, "GetFileData", fs.path)
	}

	fi, err := file.Stat()
	if err != nil {
		return inmemory.FileStorage{}, errNoFile
	}

	size := fi.Size()
	if size == 0 {
		return inmemory.FileStorage{}, errEmptyFile
	}

	var fileStorage inmemory.FileStorage

	if err := json.NewDecoder(file).Decode(&fileStorage); err != nil {
		return inmemory.FileStorage{}, err
	}

	defer file.Close()
	defer func() {
		err = util.ErrorWrap("can't load to storage", err)
	}()

	return fileStorage, nil
}

func (fs *Storage) Save(metrics map[string]interface{}) (err error) {
	if fs.openedFile == nil {
		flags := os.O_WRONLY | os.O_CREATE
		fs.openedFile, err = os.OpenFile(fs.path, flags, 0644)

		if err != nil {
			return fmt.Errorf("%s, Action:%s, File:%s", errOpenFile, "GetFileData", fs.path)
		}
	}

	encoder := json.NewEncoder(fs.openedFile)

	if err := encoder.Encode(metrics); err != nil {
		return err
	}

	log.Printf("metrics saved to file %s", fs.path)

	defer func() {
		err = util.ErrorWrap("can't save from storage", err)
		if err != nil {
			log.Print(err)
		}
	}()

	return nil
}

func (fs *Storage) CloseFile() {
	if fs.openedFile != nil {
		fs.openedFile.Close()
	}
}
