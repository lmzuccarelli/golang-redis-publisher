package schema

// Response schema
type Response struct {
	Name       string           `json:"name"`
	StatusCode string           `json:"statuscode"`
	Status     string           `json:"status"`
	Message    string           `json:"message"`
	Payload    *SchemaInterface `json:"payload,omitempty"`
}

// Token Schema
type TokenDetail struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

// Credentials (from JWT)
type Credentials struct {
	User           string `json:"user"`
	Password       string `json:"password"`
	CustomerNumber string `json:"customerNumber"`
}

// GenericSchema - used in the GenericHandler (complex data object)
type GenericSchema struct {
	//Url     string
	//Token   string
	//Creds   *Credentials
	Request *CustomerPayload
}

type CustomerPayload struct {
	Email     string `json:"email,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	Number    string `json:"number"`
	JwtToken  string `json:"jwttoken,omitempty"`
	Address   string `json:"address,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

// All the go microservices will using this schema
type SchemaInterface struct {
	ID         int64  `json:"_id,omitempty"`
	LastUpdate int64  `json:"lastupdate,omitempty"`
	MetaInfo   string `json:"metainfo,omitempty"`
}
