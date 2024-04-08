/*******************************************************************************
 * Copyright 2024 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package models

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/oklog/ulid/v2"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"io/ioutil"
	"math/rand"
	"time"
)

const alphanumericCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

type SampleData struct {
	Description string    `json:"description,omitempty"`
	Id          ulid.ULID `json:"id,omitempty"`
	Seed        string    `json:"seed,omitempty"`
	Signature   string    `json:"signature,omitempty"`
	Timestamp   string    `json:"timestamp,omitempty"`
}

func NewSampleData(cfg config.KeyInfo) (SampleData, error) {
	key, err := ioutil.ReadFile(cfg.Path)
	if err != nil {
		return SampleData{}, err
	}

	x := SampleData{
		Description: factoryRandomFixedLengthString(128, alphanumericCharset),
		Id:          newULID(),
		Seed:        factoryRandomFixedLengthString(64, alphanumericCharset),
		Timestamp:   time.Now().UTC().Format(time.RFC3339Nano),
	}

	keyDecoded := make([]byte, hex.DecodedLen(len(key)))
	hex.Decode(keyDecoded, key)
	signed := ed25519.Sign(keyDecoded, []byte(x.Seed))
	x.Signature = fmt.Sprintf("%x", signed)
	return x, nil
}

func factoryRandomFixedLengthString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func newULID() ulid.ULID {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	entropy := rand.New(source)

	id, _ := ulid.New(ulid.Timestamp(time.Now()), entropy)

	return id
}

type MongoRecord struct {
	Description  string    `json:"description,omitempty"`
	Id           string    `json:"id,omitempty"`
	Seed         string    `json:"seed,omitempty"`
	Signature    string    `json:"signature,omitempty"`
	Timestamp    string    `json:"timestamp,omitempty"`
	TimestampISO time.Time `json:"timestampiso,omitempty"`
	Confidence   float64   `json:"confidence"`
}

func MongoFromSampleData(data SampleData) MongoRecord {
	parsed, _ := time.Parse(time.RFC3339Nano, data.Timestamp)
	return MongoRecord{
		Description:  data.Description,
		Id:           data.Id.String(),
		Seed:         data.Seed,
		Signature:    data.Signature,
		Timestamp:    data.Timestamp,
		TimestampISO: parsed,
	}
}
