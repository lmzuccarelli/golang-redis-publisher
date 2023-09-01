package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/lmzuccarelli/golang-redis-publisher/pkg/connectors"
	"github.com/lmzuccarelli/golang-redis-publisher/pkg/schema"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

// SendPayloadHandler - api function handler that sends events to redis pub/sub bus
func SendPayloadHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var cp *schema.GenericSchema

	addHeaders(w, r)

	// read the jwt token data in the body
	// we don't use authorization header as the token can get quite large due to form data
	// ensure we don't have nil - it will cause a null pointer exception
	if r.Body == nil {
		r.Body = io.NopCloser(bytes.NewBufferString(""))
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		msg := "Body data (JWT) error : access forbidden %v"
		b := responseErrorFormat(http.StatusForbidden, w, msg, err)
		fmt.Fprintf(w, "%s", string(b))
		return
	}

	// unmarshal result from mw backend
	errs := json.Unmarshal(body, &cp)
	if errs != nil {
		msg := "SendPayloadHandler could not unmarshal input data to schema %v"
		con.Error(msg, errs)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, errs)
		fmt.Fprintf(w, "%s", string(b))
		return
	}

	// check the jwt token
	//creds, err := verifyJwtToken(bp.JwtToken)
	//if err != nil {
	//	msg := "SendPayloadHandler verifyToken  %v"
	//	con.Error(msg, err)
	//	b := responseErrorFormat(http.StatusForbidden, w, msg, err)
	//	fmt.Fprintf(w, string(b))
	//	return
	//}

	//token := os.Getenv("TOKEN")
	//url := os.Getenv("URL")
	// check for the correct parameters. Do not check text for password resets
	//msg := "SendPayloadHandler is not valid %s"
	//con.Error(msg, "")
	//b := responseErrorFormat(http.StatusInternalServerError, w, msg, "")
	//fmt.Fprintf(w, "%s", string(b))
	//return

	// do some funky transforms ;)
	//generic := &schema.GenericSchema{Request: cp}
	payload := `{ "number":"{{ .Number }}", "email":"{{ .Email }}" }`

	con.Trace("SendPayloadHandler new schema %v", cp.Request)
	tmpl := template.New("publish")
	//parse some content and generate a template
	tmp, _ := tmpl.Parse(payload)
	var tpl bytes.Buffer
	err = tmp.Execute(&tpl, cp.Request)
	if err != nil {
		con.Error("SendPayloadHandler parse template %v", errs)
		b := responseErrorFormat(http.StatusInternalServerError, w, " %v", errs)
		fmt.Fprintf(w, "%s", string(b))
		return
	}

	// now make the call to get all data
	con.Trace("SendPayloadHandler payload %s", tpl.String())

	err = con.Publish(context.Background(), os.Getenv("TOPIC"), tpl.String())
	if err != nil {
		con.Error("SendPayloadHandler publish request %v", errs)
		b := responseErrorFormat(http.StatusInternalServerError, w, " %v", errs)
		fmt.Fprintf(w, "%s", string(b))
		return
	}

	msg := "SendPayloadHandler published successfully"
	con.Debug(msg+" %v", string(body))
	response := &schema.Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: msg}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, "%s", string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"version\" : \""+os.Getenv("VERSION")+"\" , \"name\": \""+os.Getenv("NAME")+"\" }")
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// responsFormat - utility function
func responseErrorFormat(code int, w http.ResponseWriter, msg string, val ...interface{}) []byte {
	var b []byte
	response := &schema.Response{Name: os.Getenv("NAME"), StatusCode: strconv.Itoa(code), Status: "ERROR", Message: fmt.Sprintf(msg, val...)}
	w.WriteHeader(code)
	b, _ = json.MarshalIndent(response, "", "	")
	return b
}

// verifyJwtToken - private function
/*
func verifyJwtToken(tokenStr string) (*schema.Credentials, error) {
	var creds *schema.Credentials

	if tokenStr == "" {
		return creds, errors.New("jwt token is invalid/empty")
	}
	// local function
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRETKEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["user"] == nil || claims["customerNumber"] == nil {
			return creds, errors.New("JWT invalid user/customerNumber empty")
		}
		user := claims["user"].(string)
		cn := claims["customerNumber"].(string)
		creds = &schema.Credentials{User: user, Password: "", CustomerNumber: cn}
		return creds, nil
	}
	return creds, errors.New("jwt token is invalid")
}

// encodeBase64 - suffix the data with : and base64 encode the data. This is a requirement by the MW plugin
//func encodeBase64(data string) string {
//	data = data + ":"
//	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
//	return sEnc
//}
*/
