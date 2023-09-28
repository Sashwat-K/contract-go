// Copyright 2023 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ioeither

import (
	"crypto/rsa"
	"testing"

	RA "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	FIOE "github.com/IBM/fp-go/ioeither/file"
	"github.com/stretchr/testify/assert"
)

func TestPubOrCertToRsaKey(t *testing.T) {
	// convert a batch of test files
	res := F.Pipe2(
		RA.From(
			"../../build/mavenResolver/se-encrypt-basic.pub",
			"../../build/mavenResolver/se-encrypt-basic.crt",
			"../../build/mavenResolver/hpse-pipeline-dev-gen2-enclaved.crt",
		),
		RA.Map(F.Flow2(
			FIOE.ReadFile,
			IOE.ChainEitherK(pubOrCertToRsaKey),
		)),
		IOE.SequenceArray[error, *rsa.PublicKey],
	)

	// validate that we have a valid result
	assert.True(t, E.IsRight(res()))
}

// TestCryptoSignature checks if the signature works when created and verified by the crypto APIs
func TestCryptoSignature(t *testing.T) {
	SignatureTest(
		CryptoPrivateKey,
		CryptoPublicKey,
		CryptoRandomPassword(3333),
		CryptoSignDigest,
		CryptoVerifyDigest,
	)(t)
}

// TestCryptoSignature checks if the signature works when created by OpenSSL and verified by the crypto APIs
func TestOpenSSLCryptoSignature(t *testing.T) {
	SignatureTest(
		OpenSSLPrivateKey,
		CryptoPublicKey,
		OpenSSLRandomPassword(3333),
		OpenSSLSignDigest,
		CryptoVerifyDigest,
	)(t)
}

func TestCryptoPrivateKey(t *testing.T) {
	// generate a random key
	privKeyE := CryptoPrivateKey()
	// extract public key via openSSL and crypto and compare
	fpOpenSSL := F.Pipe1(
		privKeyE,
		E.Chain(OpenSSLPrivKeyFingerprint),
	)
	fpCrypto := F.Pipe1(
		privKeyE,
		E.Chain(CryptoPrivKeyFingerprint),
	)

	assert.Equal(t, fpOpenSSL, fpCrypto)
}

func byteLen(data []byte) int {
	return len(data)
}

func TestCryptRandomPassword(t *testing.T) {

	n := keylen
	pwd := CryptoRandomPassword(n)

	lenE := F.Pipe1(
		pwd(),
		E.Map[error](byteLen),
	)

	assert.Equal(t, E.Of[error](n), lenE)
}

// func TestCryptPrivKey(t *testing.T) {
// 	privKeyE := readFileE("../build/key.priv")
// 	// fingerprint from openSSL
// 	fpOpenSSL := F.Pipe1(
// 		privKeyE,
// 		E.Chain(OpenSSLPrivKeyFingerprint),
// 	)
// 	// fingerprint directly from crypto
// 	fpCrypto := F.Pipe1(
// 		privKeyE,
// 		E.Chain(CryptoPrivKeyFingerprint),
// 	)
// 	// make sure they match
// 	assert.Equal(t, fpOpenSSL, fpCrypto)
// }
