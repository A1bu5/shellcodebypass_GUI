package sgngui

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"
	"math/rand"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// OPERANDS string array containing logical & arithmetic operands
// for encoding the decoder stub
var OPERANDS = []string{"XOR", "SUB", "ADD", "ROL", "ROR", "NOT"}

// SCHEMA contains the operand and keys to apply single step encoding
type SCHEMA []struct {
	OP  string
	Key []byte
}

// Encoder struct for keeping encoder specs
type Encoder struct {
	architecture     int
	ObfuscationLimit int
	PlainDecoder     bool
	Seed             byte
	EncodingCount    int
	SaveRegisters    bool
}

// NewEncoder for creating new encoder structures
func NewEncoder(arch int) (*Encoder, error) {
	// Create with default settings
	encoder := Encoder{
		ObfuscationLimit: 50,
		PlainDecoder:     false,
		Seed:             GetRandomByte(),
		EncodingCount:    1,
		SaveRegisters:    false,
	}

	switch arch {
	case 32:
		encoder.architecture = 32
	case 64:
		encoder.architecture = 64
	default:
		return nil, errors.New("invalid architecture")
	}

	return &encoder, nil
}

// SetArchitecture sets the encoder architecture
func (encoder *Encoder) SetArchitecture(arch int) error {
	switch arch {
	case 32:
		encoder.architecture = 32
	case 64:
		encoder.architecture = 64
	default:
		return errors.New("invalid architecture")
	}
	return nil
}

// GetArchitecture returns the encoder architecture
func (encoder *Encoder) GetArchitecture() int {
	return encoder.architecture
}

// Encode function is the primary encode method for SGN
// all necessary options and parameters are contained inside the encoder struct
func (encoder *Encoder) Encode(payload []byte) ([]byte, error) {
	var final []byte
	if encoder.SaveRegisters {
		payload = append(payload, SafeRegisterSuffix[encoder.architecture]...)
	}

	// Add garbage instructions before the un-encoded payload
	garbage, err := encoder.GenerateGarbageInstructions()
	if err != nil {
		return nil, err
	}
	payload = append(garbage, payload...)

	// Apply ADFL cipher to payload
	cipheredPayload := CipherADFL(payload, encoder.Seed)
	encodedPayload, err := encoder.AddADFLDecoder(cipheredPayload)
	if err != nil {
		return nil, err
	}

	if encoder.PlainDecoder {
		final = encodedPayload
	} else {
		// Add more garbage instructions before the decoder stub
		garbage, err = encoder.GenerateGarbageInstructions()
		if err != nil {
			return nil, err
		}
		encodedPayload = append(garbage, encodedPayload...)

		// Calculate schema size
		schemaSize := ((len(encodedPayload) - len(cipheredPayload)) / (encoder.architecture / 8)) + 1
		randomSchema := encoder.NewCipherSchema(schemaSize)

		obfuscatedEncodedPayload := encoder.SchemaCipher(encodedPayload, 0, randomSchema)
		final, err = encoder.AddSchemaDecoder(obfuscatedEncodedPayload, randomSchema)
		if err != nil {
			return nil, err
		}
	}

	if encoder.EncodingCount > 1 {
		encoder.EncodingCount--
		encoder.Seed = GetRandomByte()
		final, err = encoder.Encode(final)
		if err != nil {
			return nil, err
		}
	}

	if encoder.SaveRegisters {
		final = append(SafeRegisterPrefix[encoder.architecture], final...)
	}

	return final, nil
}

// CipherADFL (Additive Feedback Loop) performs an additive feedback XOR operation
// similar to LFSR (Linear-feedback shift register) IN REVERSE ORDER !! with the supplied seed
func CipherADFL(data []byte, seed byte) []byte {
	for i := 1; i < len(data)+1; i++ {
		current := data[len(data)-i]
		data[len(data)-i] ^= seed
		seed = byte((int(current) + int(seed)) % 256)
	}
	return data
}

// SchemaCipher encodes a part of the given binary starting from the given index.
// Encoding done without using any loop conditions based on the schema values.
// Function performs logical/arithmetic operations given in the schema array.
// If invalid operand supplied function returns nil
func (encoder *Encoder) SchemaCipher(data []byte, index int, schema SCHEMA) []byte {
	for _, cursor := range schema {
		switch cursor.OP {
		case "XOR":
			binary.BigEndian.PutUint32(data[index:index+4], (binary.BigEndian.Uint32(data[index:index+4]) ^ binary.LittleEndian.Uint32(cursor.Key)))
		case "ADD":
			binary.LittleEndian.PutUint32(data[index:index+4], (binary.LittleEndian.Uint32(data[index:index+4])-binary.BigEndian.Uint32(cursor.Key))%0xFFFFFFFF)
		case "SUB":
			binary.LittleEndian.PutUint32(data[index:index+4], (binary.LittleEndian.Uint32(data[index:index+4])+binary.BigEndian.Uint32(cursor.Key))%0xFFFFFFFF)
		case "ROL":
			binary.LittleEndian.PutUint32(data[index:index+4], bits.RotateLeft32(binary.LittleEndian.Uint32(data[index:index+4]), -int(binary.BigEndian.Uint32(cursor.Key))))
		case "ROR":
			binary.LittleEndian.PutUint32(data[index:index+4], bits.RotateLeft32(binary.LittleEndian.Uint32(data[index:index+4]), int(binary.BigEndian.Uint32(cursor.Key))))
		case "NOT":
			binary.BigEndian.PutUint32(data[index:index+4], (^binary.BigEndian.Uint32(data[index : index+4])))
		}
		index += 4
	}
	return data
}

// RandomOperand generates a random operand string
func RandomOperand() string {
	return OPERANDS[rand.Intn(len(OPERANDS))]
}

// GetRandomByte generates a random single byte
func GetRandomByte() byte {
	return byte(rand.Intn(255))
}

// GetRandomBytes generates a random byte slice with given size
func GetRandomBytes(num int) []byte {
	slice := make([]byte, num)
	for i := range slice {
		slice[i] = GetRandomByte()
	}
	return slice
}

// CoinFlip implements a coin flip which returns true/false
func CoinFlip() bool {
	return rand.Intn(2) == 0
}

// NewCipherSchema generates random schema for
// using in the SchemaCipher function.
// Generated schema contains random operands and keys.
func (encoder *Encoder) NewCipherSchema(num int) SCHEMA {
	schema := make(SCHEMA, num)

	for i, cursor := range schema {
		cursor.OP = RandomOperand()
		if cursor.OP == "NOT" {
			cursor.Key = nil
		} else if cursor.OP == "ROL" || cursor.OP == "ROR" {
			cursor.Key = []byte{0, 0, 0, GetRandomByte()}
		} else {
			cursor.Key = GetRandomBytes(4)
		}
		schema[i] = cursor
	}
	return schema
}

// GetSchemaTable returns the printable encoder schema table
func GetSchemaTable(schema SCHEMA) string {
	data := &strings.Builder{}
	table := tablewriter.NewWriter(data)
	table.SetHeader([]string{"OPERAND", "KEY"})
	for _, cursor := range schema {
		if cursor.Key == nil {
			table.Append([]string{cursor.OP, "0x00000000"})
		} else {
			table.Append([]string{cursor.OP, fmt.Sprintf("0x%x", cursor.Key)})
		}
	}
	table.Render()

	return data.String()
}
