package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/gerifield/search-encrypted/index"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dataKey         = []byte("some-encryption-key-for-data$#@#") // 32 bytes
	firstNameIdxKey = "my-super-secret-encryption-key1"
	bornIdxKey      = "my-super-secret-encryption-key2"
)

type data struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Born      int    `json:"born"`
}

var indexer *index.Index

func main() {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	indexer = index.New(index.WithCaseInsensitive(true))

	err = insertSomeData(db)
	if err != nil {
		log.Println(err)
		return
	}
}

func insertSomeData(db *sqlx.DB) error {
	schema := `CREATE TABLE some_data (id int, data text, firstname_idx text, firstname_bloomidx text, born_idx text);`
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	if err = insertData(db, 1, data{FirstName: "John", LastName: "Carter", Born: 1982}); err != nil {
		return err
	}

	if err = insertData(db, 2, data{FirstName: "John", LastName: "Doe", Born: 1994}); err != nil {
		return err
	}

	if err = insertData(db, 3, data{FirstName: "John", LastName: "Wick", Born: 1999}); err != nil {
		return err
	}

	if err = insertData(db, 4, data{FirstName: "Johnatan", LastName: "Somebody", Born: 1993}); err != nil {
		return err
	}

	if err = insertData(db, 5, data{FirstName: "Somebody", LastName: "Else", Born: 2001}); err != nil {
		return err
	}

	log.Println("---------------------------------------------------------------------------------------------------------------------------------------")
	log.Println("Fetch some data via index:")

	if err = printData(db, 1); err != nil {
		return err
	}
	if err = printData(db, 2); err != nil {
		return err
	}

	log.Println("---------------------------------------------------------------------------------------------------------------------------------------")
	log.Println("Do some blind index search:")

	if err = searchFirstName(db, "John"); err != nil {
		return err
	}

	if err = searchFirstName(db, "Somebody"); err != nil {
		return err
	}

	if err = searchFirstName(db, "Nobody"); err != nil {
		return err
	}

	log.Println("---------------------------------------------------------------------------------------------------------------------------------------")
	log.Println("Do some bloom filter search:")

	if err = searchFistNameBloom(db, "John"); err != nil {
		return err
	}

	log.Println("---------------------------------------------------------------------------------------------------------------------------------------")
	log.Println("Do some bloom filter search with secondary filtering:")

	if err = searchFistNameBloomLastName(db, "John", "Doe"); err != nil {
		return err
	}

	if err = searchFistNameBloomLastName(db, "John", "Carpenter"); err != nil {
		return err
	}

	log.Println("---------------------------------------------------------------------------------------------------------------------------------------")
	log.Println("Do some search and filtering for case insensitive input:")

	if err = searchFirstName(db, "SoMeBoDy"); err != nil {
		return err
	}

	if err = searchFistNameBloomLastName(db, "JoHn", "DoE"); err != nil {
		return err
	}

	log.Println("---------------------------------------------------------------------------------------------------------------------------------------")
	log.Println("Do some bucket-like filtering for a between query:")
	if err = searchBornBetween(db, 1995, 2005); err != nil {
		return err
	}

	return nil
}

// searchBornBetween uses the generated buckets to look for values between some fields
func searchBornBetween(db *sqlx.DB, start, end int) error {
	// Generate the bucket keys from the start till the end
	// We know we use 10 for each step

	first := (start / 10) * 10
	last := (end/10)*10 + 10 // +10 to inlcude the last bucket too

	log.Printf("Looking for values between %d (%d), %d (%d)\n", start, first, end, last)

	var buckets []interface{} // the Select use []interface{}
	var placeHolders []string
	for i := first; i <= last; i += 10 {
		generated := indexer.IntBucket(bornIdxKey, i, 10)
		buckets = append(buckets, generated) // We just need the hmac, but the truncate won't hurt anyway
		placeHolders = append(placeHolders, "?")
		log.Println("Added bucket for", i, generated)
	}

	var results []struct {
		ID      int    `db:"id"`
		EncData []byte `db:"data"`
	}

	err := db.Select(&results, "SELECT id, data FROM some_data WHERE born_idx IN ("+strings.Join(placeHolders, ",")+")", buckets...)
	if err != nil {
		return err
	}

	var data data
	for _, res := range results {
		dec, _ := decryptDBData(res.EncData)
		if err = json.Unmarshal(dec, &data); err != nil {
			log.Println("Unmarshal error", err)
			continue
		}

		if start <= data.Born && data.Born <= end {
			log.Println("Result MATCHED, ID:", res.ID, "Data:", string(dec))
			continue
		}

		// No exact match, just the prefix was the same
		log.Println("Result NOT matched, ID:", res.ID, "Data:", string(dec))
	}

	return nil
}

// searchFistNameBloom works on a "truncated" index which will results multiple rows which MAY contain the result we are looking for
// After the initial search we should check the fetched results for the given name!
func searchFistNameBloom(db *sqlx.DB, firstName string) error {
	log.Printf("Looking for %s\n", firstName)

	firstNameBloomIndex := bloomIndex(indexer.HMAC(firstNameIdxKey, firstName))

	var results []struct {
		ID      int    `db:"id"`
		EncData []byte `db:"data"`
	}

	err := db.Select(&results, "SELECT id, data FROM some_data WHERE firstname_bloomidx = ?", firstNameBloomIndex)
	if err != nil {
		return err
	}

	var data data
	for _, res := range results {
		dec, _ := decryptDBData(res.EncData)
		if err = json.Unmarshal(dec, &data); err != nil {
			log.Println("Unmarshal error", err)
			continue
		}

		if data.FirstName == firstName {
			log.Println("Result MATCHED, ID:", res.ID, "Data:", string(dec))
			continue
		}

		// No exact match, just the prefix was the same
		log.Println("Result NOT matched, ID:", res.ID, "Data:", string(dec))
	}

	return nil
}

// searchFistNameBloomSecondName similar to the previous one, but since we'll decrypt the data anyway let's filter it also via last name!
func searchFistNameBloomLastName(db *sqlx.DB, firstName, lastName string) error {
	log.Printf("Looking for %s %s\n", firstName, lastName)

	firstNameBloomIndex := bloomIndex(indexer.HMAC(firstNameIdxKey, firstName))

	var results []struct {
		ID      int    `db:"id"`
		EncData []byte `db:"data"`
	}

	err := db.Select(&results, "SELECT id, data FROM some_data WHERE firstname_bloomidx = ?", firstNameBloomIndex)
	if err != nil {
		return err
	}

	var data data
	for _, res := range results {
		dec, _ := decryptDBData(res.EncData)
		if err = json.Unmarshal(dec, &data); err != nil {
			log.Println("Unmarshal error", err)
			continue
		}

		if strings.ToUpper(data.FirstName) == strings.ToUpper(firstName) && strings.ToUpper(data.LastName) == strings.ToUpper(lastName) {
			log.Println("Result MATCHED, ID:", res.ID, "Data:", string(dec))
			continue
		}

		// No exact match, just the prefix was the same
		log.Println("Result NOT matched, ID:", res.ID, "Data:", string(dec))
	}

	return nil
}

func searchFirstName(db *sqlx.DB, firstName string) error {
	// Note: If we need queries like "first name starts with letter J and last name should finish to `body`" then we can also build and index on that!
	log.Printf("Looking for %s\n", firstName)

	firstNameIndex := indexer.HMAC(firstNameIdxKey, firstName)

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

	firstnameIndex := indexer.HMAC(firstNameIdxKey, d.FirstName)
	firstnameBloomidx := bloomIndex(firstnameIndex)
	bornIdx := indexer.IntBucket(bornIdxKey, d.Born, 10)

	log.Println("ID:", index, "Encrypted:", base64.StdEncoding.EncodeToString(encrypted), "FirstName idx:", firstnameIndex, "FirstName bloom:", firstnameBloomidx, "Born index bucket:", bornIdx)
	_, err = db.Exec("INSERT INTO some_data (id, data, firstname_idx, firstname_bloomidx, born_idx) VALUES (?, ?, ?, ?, ?)", index, encrypted, firstnameIndex, firstnameBloomidx, bornIdx)
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

func bloomIndex(hmac string) string {
	// Truncate to a fixed prefix
	return hmac[:32]
}
