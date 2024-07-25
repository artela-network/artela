package keys

import (
	"fmt"
	"io"

	"sigs.k8s.io/yaml"

	cryptokeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"

	"github.com/artela-network/artela/ethereum/crypto/codec"
)

// available output formats.
const (
	OutputFormatText = "text"
	OutputFormatJSON = "json"
)

type bechKeyOutFn func(k *cryptokeyring.Record) (cryptokeyring.KeyOutput, error)

func printKeyringRecord(w io.Writer, k *cryptokeyring.Record, bechKeyOut bechKeyOutFn, output string) error {
	ko, err := bechKeyOut(k)
	if err != nil {
		return err
	}

	switch output {
	case OutputFormatText:
		if err := printTextRecords(w, []cryptokeyring.KeyOutput{ko}); err != nil {
			return err
		}

	case OutputFormatJSON:
		out, err := codec.KeysCdc.MarshalJSON(ko)
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintln(w, string(out)); err != nil {
			return err
		}
	}

	return nil
}

func printTextRecords(w io.Writer, kos []cryptokeyring.KeyOutput) error {
	out, err := yaml.Marshal(&kos)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, string(out)); err != nil {
		return err
	}

	return nil
}
