package contract

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/andybalholm/brotli"
	"github.com/ethereum/go-ethereum/crypto"
)

// ParseByteCode byte code format rules:
// 1. The bytecode format is: [4 Bytes Header][4 Bytes CheckSum][Data...]
// 2. Header format is [3 Byte Reserved][1 Byte Compression Algorithm]
// 3. The first 3 bytes of the header are reserved for future extensions.
// 4. The last byte of the header indicates the compression algorithm used (0x01 for Brotli).
// 5. The 4 Bytes CheckSum is calculated using the `keccak256(data)[:4]` rule, where `data` is the compressed bytecode.
// 6. If the bytecode does not conform to the above rules, it is considered unprocessed WASM bytecode and will be directly validated and executed.
func ParseByteCode(bytecode []byte) ([]byte, error) {
	if len(bytecode) < 8 {
		// no header
		return bytecode, nil
	}

	// Decode Header and CheckSum
	header := bytecode[:4]

	// parsed reserved parts
	for _, reserved := range header[:3] {
		if reserved != 0 {
			// not a valid header, ignore the processing
			return bytecode, nil
		}
	}

	// check compression algo
	compressionAlgo := header[3]
	if compressionAlgo == 0 {
		// not a valid compression algo, ignore the processing
		return bytecode, nil
	}

	// check checksum
	checkSum := bytecode[4:8]
	data := bytecode[8:]
	if err := verifyChecksum(data, checkSum); err != nil {
		return nil, err
	}

	var (
		decompressed []byte
		err          error
	)

	switch compressionAlgo {
	case 0x01:
		decompressed, err = brotliDecompress(data)
	default:
		err = fmt.Errorf("unsupported compression algorithm: %d", compressionAlgo)
	}

	return decompressed, err
}

// verifyChecksum validate whether the checksum is correct
func verifyChecksum(data, expectedCheckSum []byte) error {
	actualCheckSum := crypto.Keccak256(data)[:4]
	if !bytes.Equal(actualCheckSum, expectedCheckSum) {
		return errors.New("invalid checksum")
	}
	return nil
}

// brotliDecompress brotli decompress
func brotliDecompress(data []byte) ([]byte, error) {
	reader := brotli.NewReader(bytes.NewReader(data))
	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data: %w", err)
	}
	return decompressedData, nil
}
