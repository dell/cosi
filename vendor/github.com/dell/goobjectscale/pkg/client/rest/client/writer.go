//
//
//  Copyright © 2021 - 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//       http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//
//

package client

import "io"

var _ io.Writer = (*CountWriter)(nil) // interface guard

// CountWriter is an io.Writer that counts bytes written without actually writing anything.
type CountWriter struct {
	N int
}

// Write adds len(p) to the current value of N and returns len(p) without
// writing anything.
func (w *CountWriter) Write(p []byte) (int, error) {
	w.N += len(p)
	return len(p), nil
}
