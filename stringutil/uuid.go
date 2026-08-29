package stringutil

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/jbenet/go-base58"
	"go.gary.cool/uuid"
)

type Uuid struct {
	uuid.UUID
}

func UuidFromBytes(b []byte) (*Uuid, error) {
	if uuid, err := uuid.FromBytes(b); err == nil {
		return &Uuid{uuid}, nil
	} else {
		return nil, err
	}
}

func ParseUUID(in string) (*Uuid, error) {
	if uuid, err := uuid.Parse(in); err == nil {
		return &Uuid{
			UUID: uuid,
		}, nil
	} else {
		return nil, err
	}
}

func MustUUID(in string) *Uuid {
	if uuid, err := uuid.Parse(in); err == nil {
		return &Uuid{
			UUID: uuid,
		}
	} else {
		panic(err)
	}
}

func UUID() *Uuid {
	return &Uuid{
		UUID: uuid.New(),
	}
}

func (self *Uuid) Bytes() []byte {
	return []byte(self.UUID[:])
}

func (self *Uuid) Hex() string {
	return hex.EncodeToString(self.Bytes())
}

func (self *Uuid) Base64() string {
	return base64.StdEncoding.EncodeToString(self.Bytes())
}

func (self *Uuid) Base58() string {
	return base58.EncodeAlphabet(self.Bytes(), base58.BTCAlphabet)
}
