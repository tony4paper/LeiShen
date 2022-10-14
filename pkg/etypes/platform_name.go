package etypes

import (
	"bytes"
	"encoding/gob"

	"github.com/ethereum/go-ethereum/common"
)

func init() {
	gob.Register(new(PlatformName))
}

func DecodePlatformName(data []byte) (*PlatformName, error) {
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	platformName := new(PlatformName)
	err := dec.Decode(platformName)
	return platformName, err
}

type PlatformName struct {
	Address common.Address
	IsEOA   bool
	NameMap map[string]struct{}
}

func (p *PlatformName) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	enc.Encode(p)
	return buf.Bytes()
}

func (p *PlatformName) AddName(name string) {
	p.NameMap[name] = struct{}{}
}

func (p *PlatformName) Names() []string {
	names := make([]string, 0, len(p.NameMap))
	for name := range p.NameMap {
		names = append(names, name)
	}
	return names
}
