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

package configs

// constant keys used by the app
const (
	// UI Components
	Auth            string = "auth"
	AuthComplete    string = "authcomplete"
	GetPagelet      string = "getPagelet"
	Initialize      string = "initialize"
	Keepalive       string = "keepalive"
	UI              string = "ui"
	WebsocketClient string = "websocketclient"

	// auth subcomponents
	Approved     string = "approved"
	Authenticate string = "authenticate"
	Denied       string = "denied"
	Refresh      string = "refresh"
	Validate     string = "validate"

	// http error message "json" to return on errors as a const
	NotFoundError       = `{"error": "Not Found", "status": 404}`
	NotImplementedError = `{"error": "Not Implemented", "status": 501}`
	BadRequestError     = `{"error": "Bad Request", "status": 400}`
	InternalServerError = `{"error": "Internal Server Error", "status": 500}`
	UnauthorizedError   = `{"error": "Unauthorized", "status": 401}`

	// time format
	HHMMSS = "15:04:05"
)
