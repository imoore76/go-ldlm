// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
This file contains common types and interfaces for the session package
*/
package session

import (
	cl "github.com/imoore76/go-ldlm/server/clientlock"
)

// sessionStorer defines the Store interface
type sessionStorer interface {
	Write(map[string][]cl.Lock) error
	Read() (map[string][]cl.Lock, error)
}
