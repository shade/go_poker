package db

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type FileDB struct {
	file *os.File
}

func _CSVToDBRecord(r []string) *Record {
	return &Record{
		Name:         r[0],
		Username:     r[1],
		PasswordHash: r[2],
	}
}

func NewFileDB(path string) *FileDB {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}

	return &FileDB{
		file,
	}
}

func (f *FileDB) Get(key DBKey) (*Record, error) {
	f.file.Seek(0, io.SeekStart)
	reader := csv.NewReader(f.file)
	records, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	for _, recordCSV := range records {
		record := _CSVToDBRecord(recordCSV)
		if record.PrimaryKey() == key {
			return record, nil
		}
	}

	return nil, nil
}

func (f *FileDB) Insert(r *Record) error {
	f.file.Seek(0, io.SeekEnd)
	writer := csv.NewWriter(f.file)

	err := writer.Write(r.ToCSVRecord())
	writer.Flush()
	return err
}

func (f *FileDB) Delete(key DBKey) {
	// TODO: mayb not be necessary just yet.
	return
}
