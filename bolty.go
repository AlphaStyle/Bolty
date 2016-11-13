// Package bolty is a simple and minimalistic
// wrapper for BoltDB
package bolty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
)

// DB session struct
type DB struct {
	*bolt.DB
}

// New Creates a new BoltDB session
func New(filepath string) (*bolt.DB, error) {
	// New BoltDB session
	db, err := bolt.Open(filepath, 0600, nil)
	// Return the new Session and error if any
	return db, err
}

// Bucket creates a new bucket in BoltDB store
func (db *DB) Bucket(bucket string) error {
	// Create the bucket
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bucket))
		return err
	})

	// return error if any or allready exist
	return err
}

// Set will update or add new key and/or value to BoltDB store
func (db *DB) Set(bucket string, key string, value interface{}) error {
	// Convert value from interface to bytes
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// update the BoltDB store
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err = b.Put([]byte(key), []byte(val))
		return err
	})

	// Return error if any
	return err
}

// Get will get the value from BoltDB store
func (db *DB) Get(bucket string, key string) []byte {
	// Get the key and value
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get([]byte(key))

		// return error and value
		return fmt.Errorf("%s", v)
	})

	// return the value or error
	return []byte(err.Error())
}

// Delete will delete key and value from BoltDB store
func (db *DB) Delete(bucket string, key string) error {
	// Delete key if exist
	err := db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucket)).Delete([]byte(key))
	})

	// return error if any
	return err
}

// Seek is search function for BoltDB Store
func (db *DB) Seek(bucket string, key string) error {
	err := db.View(func(tx *bolt.Tx) error {

		c := tx.Bucket([]byte(bucket)).Cursor()

		prefix := []byte(key)
		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		// TODO return results

		return nil
	})

	// return error if any
	return err
}

// RequestBody struct
type RequestBody struct {
	Action string `json:"action"`
	Name   string `json:"name"`
	Pass   string `json:"pass"`
	Mail   string `json:"mail"`
}

// UserHandler handles all user requests such as login, logout, register
func (db *DB) UserHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Get the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("r.Body error: %s\n", err)
	}

	// Make the Request Body to JSON
	var rBody RequestBody
	json.Unmarshal(body, &rBody)

	// TODO add better validation
	if len(rBody.Name) >= 5 && len(rBody.Pass) >= 5 {
		ok := "success"

		// Check the Request action
		switch rBody.Action {
		case "REGISTER":
			fmt.Printf("%s\n", rBody)
			// TODO add to BoltDB if not exist

			if !strings.Contains(rBody.Mail, "@") {
				ok = "failed"
			}

			status, err := json.Marshal(ok)
			if err != nil {
				log.Println(err)
			}

			w.Write(status)
		case "LOGIN":
			fmt.Printf("%s\n", rBody)
			// TODO check if exist in BoltDB and make cookie

			status, err := json.Marshal(ok)
			if err != nil {
				log.Println(err)
			}

			w.Write(status)

		case "LOGOUT":
			fmt.Printf("%s\n", rBody)
			// TODO check if exist in BoltDN and delete cookie

			status, err := json.Marshal(ok)
			if err != nil {
				log.Println(err)
			}

			w.Write(status)

		default:
			log.Printf("unknowen action: %s \n", rBody.Action)

			ok = "invalid"
			status, err := json.Marshal(ok)
			if err != nil {
				log.Println(err)
			}

			w.Write(status)
		}
	}
}
