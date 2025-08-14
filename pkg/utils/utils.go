/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package utils

import (
	"time"
)

// CurrentTimeInMillis returns current time as a long int pointer NARF!
func CurrentTimeInMillis() *int64 {
	now := time.Now().UnixMilli()
	return &now
}

// BoolPointer true or false as a boolean pointer POINT!
func BoolPointer(b bool) *bool {
	return &b
}

// StringPointer return a pointer of a string as a convenience function
func StringPointer(s string) *string {
	return &s
}
