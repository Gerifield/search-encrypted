package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dataKey         = []byte("some-encryption-key-for-data$#@#") // 32 bytes
	firstNameIdxKey = "my-super-secret-encryption-key1"
)

type data struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	err = insertSomeData(db)
	if err != nil {
		log.Println(err)
		return
	}
}

func insertSomeData(db *sqlx.DB) error {
	schema := `CREATE TABLE some_data (id int, data text, firstname_idx text, index2 text);`
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	if err = insertData(db, 1, data{FirstName: "John", LastName: "Carter"}); err != nil {
		return err
	}

	if err = insertData(db, 2, data{FirstName: "John", LastName: "Doe"}); err != nil {
		return err
	}

	if err = insertData(db, 3, data{FirstName: "John", LastName: "Wick"}); err != nil {
		return err
	}

	if err = insertData(db, 4, data{FirstName: "Johnatan", LastName: "Somebody"}); err != nil {
		return err
	}

	if err = insertData(db, 5, data{FirstName: "Somebody", LastName: "Else"}); err != nil {
		return err
	}

	log.Println("---------------------------------------------------------------------------------------------------------------------------------------")
	log.Println("Test some data:")

	if err = printData(db, 1); err != nil {
		return err
	}
	if err = printData(db, 2); err != nil {
		return err
	}

	log.Println("---------------------------------------------------------------------------------------------------------------------------------------")
	log.Println("Do some index search:")

	if err = searchFirstName(db, "John"); err != nil {
		return err
	}

	if err = searchFirstName(db, "Somebody"); err != nil {
		return err
	}

	if err = searchFirstName(db, "Nobody"); err != nil {
		return err
	}

	return nil
}

func searchFirstName(db *sqlx.DB, firstName string) error {
	log.Printf("Looking for %s\n", firstName)

	firstNameIndex := generateHMACIndex(firstNameIdxKey, firstName)

	var results []struct {
		ID      int    `db:"id"`
		EncData []byte `db:"data"`
	}

	err := db.Select(&results, "SELECT id, data FROM some_data WHERE firstname_idx = ?", firstNameIndex)
	if err != nil {
		return err
	}

	for _, res := range results {
		dec, _ := decryptDBData(res.EncData)
		log.Println("Result, ID:", res.ID, "Data:", string(dec))
	}

	return nil
}

func insertData(db *sqlx.DB, index int, d data) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	log.Println("ID:", index, "Data:", d)
	encrypted, err := encryptDBData(b)
	if err != nil {
		return err
	}

	firstnameIndex := generateHMACIndex(firstNameIdxKey, d.FirstName)

	log.Println("ID:", index, "Encrypted:", base64.StdEncoding.EncodeToString(encrypted), "FirstName idx:", firstnameIndex)
	_, err = db.Exec("INSERT INTO some_data (id, data, firstname_idx) VALUES (?, ?, ?)", index, encrypted, firstnameIndex)
	return err
}

func printData(db *sqlx.DB, index int) error {
	var d []byte

	err := db.Get(&d, "SELECT data FROM some_data WHERE id = ?", index)
	if err != nil {
		return err
	}

	log.Println("ID:", index, "Encrypted:", base64.StdEncoding.EncodeToString(d))
	dec, err := decryptDBData(d)
	if err != nil {
		return err
	}

	var data data
	err = json.Unmarshal(dec, &data)
	if err != nil {
		return err
	}

	log.Println("ID:", index, "Data:", data)
	return nil
}

func encryptDBData(b []byte) ([]byte, error) {
	c, err := aes.NewCipher(dataKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, b, nil), nil
}

func decryptDBData(b []byte) ([]byte, error) {
	c, err := aes.NewCipher(dataKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(b) < nonceSize {
		return nil, err
	}

	nonce, ciphertext := b[:nonceSize], b[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext, err
}

func generateHMACIndex(key string, field string) string {
	sig := hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(field))

	return hex.EncodeToString(sig.Sum(nil))
}
