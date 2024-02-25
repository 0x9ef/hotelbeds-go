// Copyright (c) 2024 0x9ef. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.
package hotelbeds

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// AuditData represents useful debug information about
// timestamp when request happened, time to process, server ID, etc...
type AuditData struct {
	ProcessTime  ProcessTime  `json:"processTime"`
	Timestamp    Timestamp    `json:"timestamp"`
	RequestHosts Hosts        `json:"requestHost"`
	Environments Environments `json:"environment"`
	ServerID     string       `json:"serverId"`
	Release      string       `json:"release,omitempty"`
	Token        string       `json:"token"`
	Internal     string       `json:"internal"`
}

type ProcessTime time.Duration

func (t *ProcessTime) UnmarshalJSON(data []byte) error {
	n, err := strconv.Atoi(trimUnescapeQuotes(data))
	if err != nil {
		return fmt.Errorf("failed to parse ProcessTime: %w", err)
	}
	*t = ProcessTime(n)
	return nil
}

type Hosts []string

func (rh *Hosts) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	str := strings.ReplaceAll(trimUnescapeQuotes(data), " ", "")
	*rh = strings.Split(str, ",")
	return nil
}

type Environments []string

func (rh *Environments) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	str := trimUnescapeQuotes(data)
	if str[0] == '[' {
		str = str[1 : len(str)-1]
	}
	str = strings.ReplaceAll(str, " ", "")
	*rh = strings.Split(str, ",")
	return nil
}

type CommaSliceString []string

func (s CommaSliceString) MarshalJSON() ([]byte, error) {
	return []byte("\"" + strings.Join(s, ",") + "\""), nil
}

func (s *CommaSliceString) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	*s = strings.Split(string(data), ",")
	return nil
}

type CommaSliceInt []int

func (s CommaSliceInt) MarshalJSON() ([]byte, error) {
	var sb strings.Builder
	for _, elem := range s {
		sb.WriteString(strconv.Itoa(elem) + ",")
	}
	str := sb.String()[:sb.Len()-1]
	return []byte("\"" + str + "\""), nil
}

func (s *CommaSliceInt) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	slice := make([]int, bytes.Count(data, []byte(",")))
	for i, elem := range strings.Split(trimUnescapeQuotes(data), ",") {
		n, err := strconv.Atoi(elem)
		if err != nil {
			return err
		}
		slice[i] = n
	}
	*s = slice
	return nil
}
