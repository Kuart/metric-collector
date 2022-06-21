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
	noFileErr     = errors.New("there is no file in the selected path")
	openFileErr   = errors.New("error opening file")
	noFileInfoErr = errors.New("can't get file info")
	emptyFileErr  = errors.New("the file in the selected path is empty")
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
	defer func() {
		err = util.ErrorWrap("can't load to storage", err)
		if err != nil {
			log.Print(err)
		}
	}()

	if !isRestore {
		return nil
	}

	file, err := os.OpenFile(fs.path, os.O_RDONLY, 0)
	defer file.Close()

	if err != nil {
		return openFileErr
	}

	fi, err := file.Stat()
	if err != nil {
		return noFileInfoErr
	}

	size := fi.Size()
	if size == 0 {
		return emptyFileErr
	}

	var fileStorage storage.FileStorage

	if err := json.NewDecoder(file).Decode(&fileStorage); err != nil {
		return err
	}

	fs.storage.UpdateFromFile(fileStorage)

	return nil
}

func (fs Storage) InitSaver() {
	ticker := time.NewTicker(fs.interval)

	for {
		<-ticker.C
		fs.SaveToFile()
	}
}

func (fs *Storage) SaveToFile() (err error) {
	defer func() {
		err = util.ErrorWrap("can't save from storage", err)
		if err != nil {
			log.Print(err)
		}
	}()

	flags := os.O_WRONLY | os.O_CREATE

	file, err := os.OpenFile(fs.path, flags, 0644)
	defer file.Close()

	if err != nil {
		return openFileErr
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
	return nil
}
