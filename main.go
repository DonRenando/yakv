package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	mux "github.com/gorilla/mux"
)

// Globally-available key-value store.
var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

// Logger format strings.
var ftlWriteFormat = "%d\t%d\t%q\t%q\n"
var ftlReadFormat = "%d\t%d\t%q\t%q"

// Initializing logger.
var logger TransactionLogger

// This error is raised when a key is not found in the store.
var ErrorNoSuchKey = errors.New("Key doesn't exist!")

// The interface for a transaction logger.
type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error
	ReadEvents() (<-chan Event, <-chan error)
	Log()
}

// Struct for the file-based transaction logger.
type FileTransactionLogger struct {
	events chan<- Event // Write-only channel for sending events.
	errors <-chan error // Read-only channel for receiving errors.
	lastID uint64       // Last used event ID.
	file   *os.File     // Path for the transaction log.
}

// Event holds the basic information for an event.
type Event struct {
	ID        uint64    // ID assigned to the event.
	EventType EventType // The type of event assigned to the event.
	Key       string    // The key assigned to the event.
	Value     string    // The value assigned to the event.
}

type EventType byte

// Assigns a constant value for each event.
const (
	_                     = iota
	EventDelete EventType = iota
	EventPut    EventType = iota
)

// Struct for defining DELETE request body structure.
type DeleteBody struct {
	Key string
}

// Struct for defining GET request body structure.
type GetBody struct {
	Key string
}

// Struct for defining PUT request body structure.
type PutBody struct {
	Key   string
	Value string
}

// Config struct for connections.
var config struct {
	port int
	host string
}

// Put takes a key and a value as arguments, and sets the value to the given key.
func Put(key string, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()

	return nil
}

// Get takes a key as an argument, and gets the value assigned to the key.
func Get(key string) (string, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

// Delete takes a key as an argument, and deletes it from the store.
func Delete(key string) error {
	delete(store.m, key)

	return nil
}

// Malformed request struct
type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

// DecodeJSONBody parses the JSON response and returns an appropriate request.
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value := r.Header.Get("Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	// Limit size of incoming request body
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

// Handler function for DELETE
func DeleteHandler(rw http.ResponseWriter, r *http.Request) {
	var body DeleteBody

	// Use custom JSON decoder
	decodeErr := DecodeJSONBody(rw, r, &body)
	defer r.Body.Close()

	if decodeErr != nil {
		var mr *malformedRequest

		// Match errors with malformed requests
		if errors.As(decodeErr, &mr) {
			http.Error(rw, mr.msg, mr.status)
		} else {
			log.Println(decodeErr.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Get key from DeleteBody struct
	key := body.Key

	// Calls Delete for deleting a key-value pair
	err := Delete(key)

	fmt.Println("deleting key:", key)
	if errors.Is(err, ErrorNoSuchKey) {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Any other error that can't be handled
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the DELETE event to the log.
	logger.WriteDelete(key)
}

// Handler function for GET.
func GetHandler(rw http.ResponseWriter, r *http.Request) {
	var body GetBody

	// Use custom JSON decoder
	decodeErr := DecodeJSONBody(rw, r, &body)
	defer r.Body.Close()

	if decodeErr != nil {
		var mr *malformedRequest

		// Match errors with malformed requests
		if errors.As(decodeErr, &mr) {
			http.Error(rw, mr.msg, mr.status)
		} else {
			log.Println(decodeErr.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Get key from GetBody struct
	key := body.Key

	// Calls Get to get the value assigned to the key
	value, err := Get(key)

	fmt.Printf("value found for key \"%s\", value: %s\n", key, string(value))
	if errors.Is(err, ErrorNoSuchKey) {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Any other error that can't be handled
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// ResponseWriter takes byte as argument
	rw.Write([]byte(value))
}

// Handler function for PUT
func PutHandler(rw http.ResponseWriter, r *http.Request) {
	var body PutBody

	// Use custom JSON decoder
	decodeErr := DecodeJSONBody(rw, r, &body)
	defer r.Body.Close()

	if decodeErr != nil {
		var mr *malformedRequest

		// Match errors with malformed requests
		if errors.As(decodeErr, &mr) {
			http.Error(rw, mr.msg, mr.status)
		} else {
			log.Println(decodeErr.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Get key and value from PutBody struct
	key := body.Key
	value := body.Value

	if decodeErr != nil {
		http.Error(rw, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	// Call the Put function to add a key-value pair.
	err := Put(key, strings.Replace(string(value), "\n", "", -1))

	fmt.Printf("added value: \"%s\" to key \"%s\"\n", string(value), key)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the PUT event to the log.
	logger.WritePut(key, string(value))
	rw.WriteHeader(http.StatusCreated)
}

// Sends events of type EventPut to the file-based transaction logger's events channel.
func (ftl *FileTransactionLogger) WritePut(key, value string) {
	ftl.events <- Event{EventType: EventPut, Key: key, Value: value}
}

// Sends events of type EventDelete to the file-based transaction logger's events channel.
func (ftl *FileTransactionLogger) WriteDelete(key string) {
	ftl.events <- Event{EventType: EventDelete, Key: key}
}

// Sends errors to the file-based transaction logger's errors channel
func (ftl *FileTransactionLogger) Err() <-chan error {
	return ftl.errors
}

// Logs transactions to the transaction log.
func (ftl *FileTransactionLogger) Log() {
	// Buffered channel for events.
	events := make(chan Event, 16)
	ftl.events = events

	// Non-blocking buffer for sending errors.
	errors := make(chan error, 1)
	ftl.errors = errors

	// Goroutine retrieves events from the events channel.
	go func() {
		for e := range events {
			ftl.lastID++

			// Log the transaction in the log file.
			_, err := fmt.Fprintf(ftl.file, ftlWriteFormat, ftl.lastID, e.EventType, e.Key, strings.TrimSpace(e.Value))

			if err != nil {
				// Send the error to errors channel.
				errors <- err
				return
			}
		}
	}()
}

// Reads all transactions from the transaction log.
func (ftl *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(ftl.file) // Scanner for transaction log
	outEvent := make(chan Event)          // Unbuffered channel for events.
	outError := make(chan error, 1)       // Buffered channel for errors.

	// Goroutine for parsing transactions.
	go func() {
		var e Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			// Scans the transaction from the log.
			if _, err := fmt.Sscanf(line, ftlReadFormat, &e.ID, &e.EventType, &e.Key, &e.Value); err != nil {
				outError <- fmt.Errorf("Failed while parsing input. %w", err)
				return
			}

			// Checks for seqeuence. Abnormal sequences are not suitable for replaying transactions.
			if ftl.lastID >= e.ID {
				outError <- fmt.Errorf("Transaction IDs out of sequence.")
				return
			}

			// Last used ID is updated to current value.
			ftl.lastID = e.ID

			// Sends the event to the outEvent channel.
			outEvent <- e
		}

		// Send any error to the outError channel.
		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("Failed reading transaction log. %w", err)
			return
		}
	}()

	return outEvent, outError
}

// Creates a new file-based transaction logger.
func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)

	if err != nil {
		return nil, fmt.Errorf("Failure: Failed to read transaction log file. %w", err)
	}

	return &FileTransactionLogger{file: file}, nil
}

// Initialize the transaction log and mutates the state of the key-value store by replaying previously stored transactions.
func InitLog(filename string) error {
	var err error

	// Filename for logs is "transaction.log" by default.
	logger, err = NewFileTransactionLogger(filename)
	if err != nil {
		return fmt.Errorf("Failed to create logger! %w", err)
	}

	// Reads all events and errors.
	fmt.Println("yakv is reading previous transactions from the log.... 🔎")
	events, errors := logger.ReadEvents()
	e, ok := Event{}, true

	// Checks each transaction and performs it (i.e. replaying).
	fmt.Println("yakv is replaying all previous transactions.... ⏯")
	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = Delete(e.Key)
			case EventPut:
				err = Put(e.Key, e.Value)
			}
		}
	}

	// Actively call Log() to log transactions to the transaction log.
	logger.Log()
	return err
}

func main() {
	// Filename for the transaction log extracted from the -filename flag.
	var logFilename string

	// Flag values for TLS-based connection.
	var secure string
	var certFilename string
	var keyFilename string

	// default address is 127.0.0.1:8080
	flag.IntVar(&config.port, "port", 8080, "Port Number.")
	flag.StringVar(&config.host, "host", "127.0.0.1", "Host Address.")

	// default connections are not secured using TLS
	flag.StringVar(&secure, "secure", "insecure", "TLS-secured connection.")
	flag.StringVar(&certFilename, "cert", "cert.pem", "Filename for certificate.")
	flag.StringVar(&keyFilename, "key", "key.pem", "Filename for private key.")

	// default transaction log filename is "transaction.log"
	flag.StringVar(&logFilename, "filename", "transaction.log", "Filename for the transaction log.")

	flag.Parse()

	addr := fmt.Sprintf("%s:%d", config.host, config.port)
	fmt.Printf("yakv is starting on address: %s 🥳\n", addr)
	fmt.Println("yakv is up and running! 🚀🥳")

	r := mux.NewRouter()
	fmt.Println("yakv is initializing the transaction log! 🔨")

	InitLog(logFilename)

	// yakv URLs are set to v0.
	r.HandleFunc("/yakv/v0/put", PutHandler).Methods("PUT")
	r.HandleFunc("/yakv/v0/get", GetHandler).Methods("GET")
	r.HandleFunc("/yakv/v0/delete", DeleteHandler).Methods("DELETE")

	// Handle secure flag and serve.
	if strings.TrimSpace(secure) == "tls" {
		fmt.Println("yakv is running in secure mode.... 🔒")
		log.Fatal(http.ListenAndServeTLS(addr, certFilename, keyFilename, r))
	} else {
		fmt.Println("yakv is running in insecure mode.... 🔓❎")
		log.Fatal(http.ListenAndServe(addr, r))
	}
}
