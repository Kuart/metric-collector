package file

import (
	"encoding/json"
	"errors"
	"github.com/Kuart/metric-collector/internal/storage/storage"
	"github.com/Kuart/metric-collector/internal/util"
	"log"
	"os"
	"time"
)

var (
	errNoFile     = errors.New("there is no file in the selected path")
	errOpenFile   = errors.New("error opening file")
	errNoFileInfo = errors.New("can't get file info")
	errEmptyFile  = errors.New("the file in the selected path is empty")
)

type Storage struct {
	path     string
	interval time.Duration
	storage  *storage.Storage
}

func New(path string, interval time.Duration, storage *storage.Storage) Storage {
	return Storage{
		path:     path,
		interval: interval,
		storage:  storage,
	}
}

func (fs *Storage) LoadToStorage(isRestore bool) (err error) {
	if !isRestore {
		return nil
	}

	file, err := os.OpenFile(fs.path, os.O_RDONLY, 0)

	if err != nil {
		return errOpenFile
	}

	fi, err := file.Stat()
	if err != nil {
		return errNoFile
	}

	size := fi.Size()
	if size == 0 {
		return errEmptyFile
	}

	var fileStorage storage.FileStorage

	if err := json.NewDecoder(file).Decode(&fileStorage); err != nil {
		return err
	}

	fs.storage.UpdateFromFile(fileStorage)

	defer file.Close()
	defer func() {
		err = util.ErrorWrap("can't load to storage", err)
		if err != nil {
			log.Print(err)
		}
	}()

	return nil
}

func (fs Storage) InitSaver() {
	if fs.path == "" {
		return
	}

	ticker := time.NewTicker(fs.interval)

	for {
		<-ticker.C
		fs.SaveToFile()
	}
}

func (fs *Storage) SaveToFile() (err error) {
	flags := os.O_WRONLY | os.O_CREATE

	file, err := os.OpenFile(fs.path, flags, 0644)

	if err != nil {
		return errOpenFile
	}

	temp := storage.FileStorage{
		Gauge:   fs.storage.GetGauge(),
		Counter: fs.storage.GetCounter(),
	}

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(temp); err != nil {
		return err
	}

	log.Printf("metrics saved to file %s", fs.path)

	defer file.Close()
	defer func() {
		err = util.ErrorWrap("can't save from storage", err)
		if err != nil {
			log.Print(err)
		}
	}()

	return nil
}
