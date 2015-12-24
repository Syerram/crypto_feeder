package main

import (
	"fmt"
	"os"
)

type FileWriter struct {
	BasePath     string
	FilePointers map[string]*os.File
}

func (fWriter *FileWriter) GetCanonicalName() string {
	return "file_writer"
}

// Writer implementation
func (fWriter *FileWriter) Write(tradeRow *TradeRow) (err error) {
	var io_writer, io_err = fWriter.AddOrGetWriter(tradeRow.Triplet)
	if io_err != nil {
		LOGGER.Error.Println("Error occured in getting file writer: ", io_err)
		return io_err
	}
	var _, w_err = io_writer.Write([]byte(fmt.Sprintf("%s\n", tradeRow.ToJson())))
	if w_err != nil {
		LOGGER.Error.Println("Error occured in writing data to the file: ", w_err)
	}
	return w_err
}

// Utility methods for FileWriter
func (fWriter *FileWriter) AddOrGetWriter(key string) (*os.File, error) {
	if writer, ok := fWriter.FilePointers[key]; ok {
		return writer, nil
	}
	writer, err := os.OpenFile(fmt.Sprintf("%s/%s.json", fWriter.BasePath, key), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		LOGGER.Trace.Println("Failed to open error given file:", err)
		return nil, err
	}
	fWriter.FilePointers[key] = writer
	return writer, nil
}

func (fWriter *FileWriter) Close() error {
	for _, fPt := range fWriter.FilePointers {
		fPt.Close()
	}
	return nil
}

func init() {
	//TODO: Externalize paths
	WRITERS.RegisterWriter(&FileWriter{BasePath: "/tmp/crypto_trades", FilePointers: make(map[string]*os.File)})
}
