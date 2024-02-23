// Copyright (c) 2024 0x9ef. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.
package hotelbeds

import (
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
