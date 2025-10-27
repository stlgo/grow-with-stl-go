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
	seedsKey               = "seeds"
	getInventoryRequestKey = "getInventory"
	getDetailRequestKey    = "getDetail"
	purchaseRequestKey     = "purchase"
	addInventoryRequestKey = "addInventory"
)

type purchaseRequest struct {
	Category *string      `json:"category,omitempty"`
	ID       *string      `json:"id,omitempty"`
	Quantity *interface{} `json:"quantity,omitempty"`
}

// Init is called by the launch in the pkg/commands/root.go
func Init() error {
	webservice.AppendToWebsocketFunctionMap(seedsKey, handleWebsocketRequest)
	webservice.AppendToRESTFunctionMap(seedsKey, handleRESTRequest)

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
		case getInventoryRequestKey:
			response.Data, err = getInventory()
		case getDetailRequestKey:
			response.Data, err = getDetail(request.Component, request.SubComponent)
		case purchaseRequestKey:
			response.Data, err = purchase(request.SessionID, request.Data)
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
	restURI := strings.TrimPrefix(r.RequestURI, fmt.Sprintf("/REST/v1.0.0/%s/", seedsKey))
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
	case getInventoryRequestKey:
		if strings.EqualFold(r.Method, http.MethodGet) {
			data, err := getInventory()
			if err != nil {
				log.Error(err)
				http.Error(w, configs.InternalServerError, http.StatusInternalServerError)
				return
			}
			writeRESTResponse(data, w, http.StatusOK)
		}
	case getDetailRequestKey:
		if len(uriParts) >= 2 {
			data, err := getDetail(&uriParts[1], &uriParts[2])
			if err != nil {
				log.Error(err)
				http.Error(w, configs.InternalServerError, http.StatusInternalServerError)
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
	case purchaseRequestKey:
		purchaseREST(w, r)
	case addInventoryRequestKey:
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
		http.Error(w, configs.InternalServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpStatus)
	_, err = w.Write(b)
	if err != nil {
		log.Error(err)
		http.Error(w, configs.InternalServerError, http.StatusInternalServerError)
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

		data, err := purchase(nil, requestBody)
		if err != nil {
			log.Error(err)
			http.Error(w, configs.InternalServerError, http.StatusInternalServerError)
			return
		}

		writeRESTResponse(data, w, http.StatusCreated)
		return
	}

	log.Error("bad request")
	http.Error(w, configs.BadRequestError, http.StatusBadRequest)
}

func purchase(sessionID *string, data interface{}) (*InventoryItem, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var request purchaseRequest
	unmarshalErr := json.Unmarshal(bytes, &request)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	if request.Category != nil && request.ID != nil && request.Quantity != nil {
		inventoryMutex.Lock()
		category, ok := inventoryCache[*request.Category]
		inventoryMutex.Unlock()
		if ok && category != nil {
			if item, ok := category.Items[*request.ID]; ok {
				quantity, err := determineQuantity(*request.Quantity)
				if err != nil {
					return nil, err
				}
				if quantity != nil && *item.Packets >= *quantity {
					remainder := *item.Packets - *quantity
					item.Packets = &remainder
					if err := updateItemInDB(item); err != nil {
						return nil, err
					}
					go notifyAll(request.Category, request.ID, sessionID, item)
					return item, nil
				}
				return nil, errors.New("cannot purchase more packets than what the inventory contains")
			}
		}
	}
	return nil, fmt.Errorf("item not found")
}

func notifyAll(categoryName, seedID, sessionID *string, item *InventoryItem) {
	route := seedsKey
	requestType := "purchase"
	response := &configs.WsMessage{
		Route:        &route,
		Type:         &requestType,
		Component:    categoryName,
		SubComponent: seedID,
		Data:         item,
	}
	webservice.NotifyAll(sessionID, response)
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
