package util

import (
	"contact_app_mux_gorm_main/components/apperror"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func UnmarshalJSON(request *http.Request, out interface{}) error {
	if request.Body == nil {
		fmt.Println("==============================err request.Body == nil==========================")
		return apperror.NewInvalidJSONError("Empty json body")
	}
	// fmt.Println("==============================err (request.Body)==========================")
	// fmt.Println(request.Body)
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println("==============================err ioutil.ReadAll==========================")
		return apperror.NewDatabaseError("Error while reading JSON body")
	}

	if len(body) == 0 {
		fmt.Println("==============================err len(body) == 0==========================")
		return apperror.NewInvalidJSONError("0 lenght json body")
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		fmt.Println("==============================err json.Unmarshal==========================")
		fmt.Println(body)
		fmt.Println(out)
		return apperror.NewInvalidJSONError("Error while unmarshaling the json body")
	}
	return nil
}
