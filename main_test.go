package main

import (
	"errors"
	"os"
	"testing"
)

// Helper function for checking whether a file exists or not.
func checkFileExists(filename string) bool {
	// Gets file info and error related to the filename.
	fileInfo, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	// We don't want to return true for directories
	return !fileInfo.IsDir()
}

// Helper function for checking if the current ID is equal to the logger's last ID.
func checkLastID(t *testing.T, tl TransactionLogger, id uint64) {
	if lastID := tl.LastID(); lastID != id {
		t.Errorf("Got unexpected sequence. Expected %d, got %d", id, lastID)
	} else {
		t.Logf("Got expected sequence. Expected %d, got %d", id, lastID)
	}
}

// Function for testing Get operation.
func TestGet(t *testing.T) {
	// Sample data
	const key = "yakv"
	const value = "hello, yakv!"

	// An interface for storing value after the Get operation.
	var val interface{}

	// This boolean helps in detecting in whether a key already exists or not.
	var err error

	// Restore to original state after test.
	defer delete(store.m, key)

	// Perform the Put operation.
	val, err = Get(key)
	if !errors.Is(err, ErrorNoSuchKey) {
		t.Error("Unexpected error:", err)
	}
	if err == nil {
		// We expected an error as we're calling GET on an invalid key.
		t.Error("Expected an error.")
	}

	// Assign key-value pair manually.
	store.m[key] = value

	// Repeat the Get operation.
	val, err = Get(key)
	if err != nil {
		t.Error("Unexpected error:", err)
	}

	// Check if the key holds a valid value.
	if val != value {
		t.Error("val and value don't match. Diverted from expected result.")
	}
}

// Function for testing Put operation.
func TestPut(t *testing.T) {
	// Sample data
	const key = "yakv"
	const value = "hello, yakv!"

	// An interface for storing value after the Put operation.
	var val interface{}

	// This boolean helps in detecting in whether a key already exists or not.
	var alreadyExists bool

	// Restore to original state after test.
	defer delete(store.m, key)

	// Check whether the key already exists.
	_, alreadyExists = store.m[key]
	if alreadyExists {
		t.Error("Key already exists!")
	}

	// Perform the Put operation.
	err := Put(key, value)
	if err != nil {
		t.Error(err)
	}

	// If it is not found, then it failed to create a key.
	val, alreadyExists = store.m[key]
	if !alreadyExists {
		t.Error("Failed to create the key-value pair.")
	}

	// Check if the key holds a valid value.
	if val != value {
		t.Error("val and value don't match. Diverted from expected result.")
	}
}

// Function for testing Delete operation.
func TestDelete(t *testing.T) {
	// Sample data
	const key = "yakv"
	const value = "hello, yakv!"

	// This boolean helps in detecting in whether a key already exists or not.
	var alreadyExists bool

	// Restore to original state after test.
	defer delete(store.m, key)

	// Assign key-value pair manually.
	store.m[key] = value

	// Check whether the key already exists.
	_, alreadyExists = store.m[key]
	if !alreadyExists {
		t.Error("Key doesn't exist! Cannot delete non-existent key.")
	}

	Delete(key)

	// If it is found, then it failed to delete the key.
	_, alreadyExists = store.m[key]
	if alreadyExists {
		t.Error("Failed to delete the key-value pair.")
	}
}

// Function for testing the creation of a transaction log.
func TestInitLogger(t *testing.T) {
	// Temporary log filename.
	const filename = "temp-transaction.log"

	// Restore to original state after test.
	defer os.Remove(filename)

	// Create a new file logger.
	ftl, err := NewFileTransactionLogger(filename)
	if err != nil {
		t.Errorf("Error: %w", err)
	}

	// Check whether a logger was returned from the function.
	if ftl == nil {
		t.Error("No logger was returned from NewFileTransactionLogger()")
	}

	// Check whether the file exists or not.
	if !checkFileExists(filename) {
		t.Errorf("File \"%s\" doesn't exist.", filename)
	}
}

// Function for testing if IDs are sequentially logged.
func TestIDs(t *testing.T) {
	// Temporary log filename.
	const filename = "temp-IDs.log"

	// Restore to original state after test.
	defer os.Remove(filename)

	// Create a new file logger.
	transactionLogger, err := NewFileTransactionLogger(filename)
	if err != nil {
		t.Error(err)
	}

	// Start logging.
	transactionLogger.Log()

	// Close file once tests are done.
	defer transactionLogger.Close()

	// Check expected last ID.
	checkLastID(t, transactionLogger, 0)

	// Send events to the logger.
	transactionLogger.WritePut("yakv", "yak1")
	transactionLogger.WritePut("yakv2", "yak2")
	transactionLogger.Wait()

	// Check expected last ID after WritePut.
	checkLastID(t, transactionLogger, 2)

	// Send events to the logger.
	transactionLogger.WritePut("yakv3", "yak3")
	transactionLogger.WritePut("yakv4", "yak4")
	transactionLogger.Wait()

	// Check expected last ID after WritePut.
	checkLastID(t, transactionLogger, 4)
}

// Function for testing WritePut.
func TestWritePut(t *testing.T) {
	// Temporary log filename.
	const filename = "temp-writeput.log"

	// Restore to original state after test.
	defer os.Remove(filename)

	// Create a new file logger.
	transactionLogger, _ := NewFileTransactionLogger(filename)

	// Start logging.
	transactionLogger.Log()

	// Close file once tests are done.
	defer transactionLogger.Close()

	// Send events to the logger.
	transactionLogger.WritePut("yakv1", "yak1")
	transactionLogger.WritePut("yakv2", "yak2")
	transactionLogger.Wait()

	// Create a new file logger.
	transactionLogger2, _ := NewFileTransactionLogger(filename)

	// Read events and errors.
	inEvents, inErrors := transactionLogger2.ReadEvents()

	// Close file once tests are done.
	defer transactionLogger2.Close()

	// Log events to the console.
	for e := range inEvents {
		t.Log(e)
	}

	// Print error to the console.
	err := <-inErrors
	if err != nil {
		t.Error(err)
	}

	// Check if both loggers have the same last ID.
	if transactionLogger.LastID() != transactionLogger2.LastID() {
		t.Errorf("IDs are not matching: %d != %d", transactionLogger.LastID(), transactionLogger2.LastID())
	}
}
