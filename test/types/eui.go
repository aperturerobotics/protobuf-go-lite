// Copyright © 2021 The Things Industries B.V.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/TheThingsIndustries/protoc-gen-go-json/jsonplugin"
)

type EUI64 [8]byte

func (t EUI64) Marshal() ([]byte, error) {
	return t[:], nil
}

func (t *EUI64) MarshalTo(data []byte) (n int, err error) {
	return copy(data, t[:]), nil
}

func (t *EUI64) Unmarshal(data []byte) error {
	if len(data) != 8 {
		return fmt.Errorf("invalid data length: got %d, want 8", len(data))
	}
	var dto EUI64
	copy(dto[:], data)
	*t = dto
	return nil
}

func (t *EUI64) Size() int { return 8 }

func MarshalHEX(s *jsonplugin.MarshalState, b []byte) {
	if b == nil {
		s.WriteNil()
		return
	}
	s.WriteString(fmt.Sprintf("%X", b))
}

func UnmarshalHEX(s *jsonplugin.UnmarshalState) []byte {
	str := s.ReadString()
	if s.Err() != nil {
		return nil
	}
	b, err := hex.DecodeString(str)
	if err != nil {
		s.SetError(err)
		return nil
	}
	return b
}

func MarshalHEXArray(s *jsonplugin.MarshalState, bs [][]byte) {
	s.WriteArrayStart()
	var wroteElement bool
	for _, b := range bs {
		s.WriteMoreIf(&wroteElement)
		s.WriteString(fmt.Sprintf("%X", b))
	}
	s.WriteArrayEnd()
}

func UnmarshalHEXArray(s *jsonplugin.UnmarshalState) [][]byte {
	var bs [][]byte
	s.ReadArray(func() {
		bs = append(bs, UnmarshalHEX(s))
	})
	if s.Err() != nil {
		return nil
	}
	return bs
}

func MarshalStringHEXMap(s *jsonplugin.MarshalState, bs map[string][]byte) {
	s.WriteObjectStart()
	var wroteElement bool
	for k, b := range bs {
		s.WriteMoreIf(&wroteElement)
		s.WriteObjectField(k)
		s.WriteString(fmt.Sprintf("%X", b))
	}
	s.WriteObjectEnd()
}

func UnmarshalStringHEXMap(s *jsonplugin.UnmarshalState) map[string][]byte {
	bs := make(map[string][]byte)
	s.ReadObject(func(key string) {
		bs[key] = UnmarshalHEX(s)
	})
	if s.Err() != nil {
		return nil
	}
	return bs
}

func (t EUI64) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%X"`, t[:])), nil
}

func (t *EUI64) UnmarshalJSON(data []byte) error {
	hexStr, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	b, err := hex.DecodeString(hexStr)
	var dto EUI64
	copy(dto[:], b)
	*t = dto
	return nil
}

func (t *EUI64) MarshalProtoJSON(s *jsonplugin.MarshalState) {
	if t == nil {
		s.WriteNil()
		return
	}
	s.WriteString(fmt.Sprintf("%X", t[:]))
}

func (t *EUI64) UnmarshalProtoJSON(s *jsonplugin.UnmarshalState) {
	b, err := hex.DecodeString(s.ReadString())
	if err != nil {
		s.SetError(err)
	}
	var dto EUI64
	copy(dto[:], b)
	*t = dto
}
