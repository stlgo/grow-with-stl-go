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

package seeds

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/webservice"
)

const (
	seeds               = "seeds"
	getInventoryRequest = "getInventory"
	getDetailRequest    = "getDetail"
	purchaseRequest     = "purchase"
	addInventoryRequest = "addInventory"
)

// Init is called by the launch in the pkg/commands/root.go
func Init() error {
	webservice.AppendToWebsocketFunctionMap(seeds, handleWebsocketRequest)
	webservice.AppendToRESTFunctionMap(seeds, handleRESTRequest)

	// make sure the tables exist
	if err := setupTables(); err != nil {
		return err
	}

	// warm up the inventory cache
	if _, err := getInventory(); err != nil {
		return err
	}

	log.Debugf("%d categories in inventory cache", len(inventoryCache))
	return nil
}

func handleWebsocketRequest(request, response *configs.WsMessage) {
	if request.Type != nil && response != nil {
		var err error
		switch *request.Type {
		case getInventoryRequest:
			response.Data, err = getInventory()
		case getDetailRequest:
			response.Data, err = getDetail(request.Component, request.SubComponent)
		case purchaseRequest:
			response.Data, err = purchase(request.Component, request.SubComponent, request.Data)
			if err == nil {
				webservice.NotifyAll(request.SessionID, response)
			}
		default:
			err = fmt.Errorf("type %s not implemented", *request.Type)
		}

		if response.Data == nil {
			e := "no data found"
			log.Error(e)
			response.Error = &e
		}

		if err != nil {
			log.Error(err)
			e := err.Error()
			response.Error = &e
			response.Data = nil
		}
		return
	}
}

func handleRESTRequest(w http.ResponseWriter, r *http.Request) {
	restURI := strings.TrimPrefix(r.RequestURI, fmt.Sprintf("/REST/v1.0.0/%s/", seeds))
	uriParts := strings.Split(restURI, "/")
	w.Header().Set("content-type", "application/json")
	switch r.Method {
	case http.MethodGet:
		getHelper(restURI, uriParts, w, r)
	case http.MethodPost:
		postHelper(restURI, uriParts, w, r)
	default:
		log.Errorf("%s URI with Method %s not possible", restURI, r.Method)
		http.Error(w, configs.NotImplementedError, http.StatusNotImplemented)
	}
}

func getHelper(restURI string, uriParts []string, w http.ResponseWriter, r *http.Request) {
	switch uriParts[0] {
	case getInventoryRequest:
		if strings.EqualFold(r.Method, http.MethodGet) {
			data, err := getInventory()
			if err != nil {
				log.Error(err)
				http.Error(w, configs.IntenralServerError, http.StatusInternalServerError)
				return
			}
			writeRESTResponse(data, w, http.StatusOK)
		}
	case getDetailRequest:
		if len(uriParts) >= 2 {
			data, err := getDetail(&uriParts[1], &uriParts[2])
			if err != nil {
				log.Error(err)
				http.Error(w, configs.IntenralServerError, http.StatusInternalServerError)
				return
			}
			writeRESTResponse(data, w, http.StatusOK)
			return
		}
		http.Error(w, configs.BadRequestError, http.StatusBadRequest)
	default:
		log.Errorf("%s URI with Method %s not possible", restURI, r.Method)
		http.Error(w, configs.NotImplementedError, http.StatusNotImplemented)
	}
}

func postHelper(restURI string, uriParts []string, w http.ResponseWriter, r *http.Request) {
	switch uriParts[0] {
	case purchaseRequest:
		purchaseREST(w, r)
	case addInventoryRequest:
		http.Error(w, configs.NotImplementedError, http.StatusNotImplemented)
	default:
		log.Errorf("%s URI with Method %s not possible", restURI, r.Method)
		http.Error(w, configs.NotImplementedError, http.StatusNotImplemented)
	}
}

func writeRESTResponse(data interface{}, w http.ResponseWriter, httpStatus int) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		http.Error(w, configs.IntenralServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpStatus)
	_, err = w.Write(b)
	if err != nil {
		log.Error(err)
		http.Error(w, configs.IntenralServerError, http.StatusInternalServerError)
		return
	}
}

func purchaseREST(w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold(r.Method, http.MethodPost) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, configs.BadRequestError, http.StatusBadRequest)
			return
		}

		var requestBody map[string]interface{}
		if err := json.Unmarshal(body, &requestBody); err != nil {
			log.Errorf("bad request body: %s", body)
			http.Error(w, configs.BadRequestError, http.StatusBadRequest)
			return
		}

		if category, ok := requestBody["category"].(string); ok {
			if seedID, ok := requestBody["id"].(string); ok {
				data, err := purchase(&category, &seedID, requestBody)
				if err != nil {
					log.Error(err)
					http.Error(w, configs.IntenralServerError, http.StatusInternalServerError)
					return
				}

				writeRESTResponse(data, w, http.StatusCreated)
				return
			}
		}
	}

	log.Error("bad request")
	http.Error(w, configs.BadRequestError, http.StatusBadRequest)
}

func purchase(categoryName, seedID *string, data interface{}) (*InventoryItem, error) {
	if categoryName != nil && seedID != nil && data != nil {
		if rawData, ok := data.(map[string]interface{}); ok {
			if category, ok := inventoryCache[*categoryName]; ok && category != nil {
				if item, ok := category.Items[*seedID]; ok {
					quantity, err := determineQuantity(rawData["quantity"])
					if err != nil {
						return nil, err
					}
					if quantity != nil && *item.Packets > *quantity {
						remainder := *item.Packets - *quantity
						item.Packets = &remainder
						if err := updateItemInDB(item); err != nil {
							return nil, err
						}
						return item, nil
					}
					return nil, errors.New("cannot purchase more packets than what the inventory contains")
				}
			}
		}
	}
	return nil, fmt.Errorf("item not found")
}

func determineQuantity(data interface{}) (*int, error) {
	var quantity int
	switch t := data.(type) {
	case string:
		q, err := strconv.Atoi(t)
		if err != nil {
			return nil, err
		}
		quantity = q
	case float32:
		quantity = int(t)
	case float64:
		quantity = int(t)
	case int:
		quantity = t
	case int32:
		quantity = int(t)
	case int64:
		quantity = int(t)
	default:
		return nil, fmt.Errorf("cannot determine quantity type %T", data)
	}
	return &quantity, nil
}
