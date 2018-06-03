package main_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// func init() {
// 	log.SetFlags(log.LstdFlags)
// }

// APITestSuite type
type APITestSuite struct {
	suite.Suite
	Server   *httptest.Server
	Response *http.Response
	Request  *http.Request
	err      error
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Success bool `json:"success"`
	}{true})
}

// SetupSuite func
func (suite *APITestSuite) SetupSuite() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheckHandler)
	suite.Server = httptest.NewServer(mux)
}

// TearDownSuite func
func (suite *APITestSuite) TearDownSuite() {
	suite.Server.Close()
}

// SetupTest func
func (suite *APITestSuite) SetupTest() {
	suite.Response = nil
	suite.Request = nil
}

// TearDownTest func
func (suite *APITestSuite) TearDownTest() {
	suite.Response = nil
	suite.Request = nil
}

// TestAPISuite func
func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

// TestHealthAPI func
func (suite *APITestSuite) TestHealthAPI() {
	// Make POST request with JSON body
	// userJson := `{"username": "dennis", "balance": 200}`
	// reader = strings.NewReader(userJson) //Convert string to reader
	// request, err := http.NewRequest("POST", usersUrl, reader)

	// Make GET request
	suite.Request, suite.err = http.NewRequest("GET", suite.Server.URL+"/health", nil)
	suite.Response, suite.err = http.DefaultClient.Do(suite.Request)
	if suite.err != nil {
		suite.Error(suite.err) //Something is wrong while sending request
	}
	defer suite.Response.Body.Close()
	body, _ := ioutil.ReadAll(suite.Response.Body)
	suite.JSONEq(`{"success":true}`, string(body))
	suite.Equal(http.StatusOK, suite.Response.StatusCode)
}

// TestHealthHandler func
func TestHealthHandler(t *testing.T) {
	// TDD
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := executeRequest(healthCheckHandler, req)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.JSONEq(t, `{"success":true}`, res.Body.String())

	// BDD
	Convey("Given target url", t, func() {
		req, err := http.NewRequest("GET", "/health", nil)
		if err != nil {
			t.Fatal(err)
		}
		Convey("When make request", func() {
			res := executeRequest(healthCheckHandler, req)
			Convey("It should status ok", func() {
				So(res.Code, ShouldEqual, http.StatusOK)
				So(res.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func executeRequest(handlerFunc http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	return rr
}

/*
package postgres_test

import (
	"testing"

	_ "github.com/lib/pq"
	// "github.com/payboxth/cloud/postgres"
	// . "github.com/smartystreets/goconvey/convey"
)

// func init() {
// 	log.SetFlags(log.LstdFlags)
// }

// func TestVendingRepo_Ota(t *testing.T) {
// 	tests := map[string]struct {
// 		clientID  int64
// 		vendingID int64
// 		otaType   vending.OtaType
// 		smsErr    error
// 		err       error
// 	}{
// 		"normal update": {0, 2, vending.OTAUPDATE, nil, nil},
// 		"new machine":   {0, 2, vending.OTANEW, nil, nil},
// 	}

// 	db := connectDB()
// 	defer db.Close()
// 	v, err := postgres.NewVendingRepository(db)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	for k, test := range tests {
// 		fmt.Println("test #x : TestVendingRepo_Ota in case : " + k)
// 		// log.Println(test.err)
// 		// expect := &mobile.Config{}
// 		expect := test.err
// 		// log.Println(test.otaType)
// 		_, err := v.Ota(test.clientID, test.vendingID, test.otaType)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return
// 		}
// 		assert.Equal(t, expect, 1, "")
// 	}

// 	// newOta := &vending.Ota{}

// 	// expect := "success"

// 	// db := connectDB()
// 	// defer db.Close()
// 	// x, err := postgres.NewMobileRepository(db)
// 	// if err != nil {
// 	// 	fmt.Println(err.Error())
// 	// 	return
// 	// }

// 	// err = x.StoreConfig(0, newConfig)
// 	// if err != nil {
// 	// 	actual = "fail"
// 	// 	fmt.Println("fail -->", err.Error())
// 	// } else {
// 	// 	actual = "success"
// 	// }

// 	// assert.Equal(t, expect, actual, "err")
// }

// Internal function testing
func TestVersionConvertor(t *testing.T) {
	// fmt.Println("test #x : Test OTA version convertor version <==> number")
	// tests := []struct {
	// 	number  int64
	// 	version string
	// }{
	// 	{100100001, "10.10.1"},
	// 	{100100010, "10.10.10"},
	// 	{100100100, "10.10.100"},
	// 	{100101000, "10.10.1000"},

	// 	{120100001, "12.10.1"},
	// 	{100120010, "10.12.10"},
	// 	{100100120, "10.10.120"},
	// 	{100101002, "10.10.1002"},

	// 	{30010000, "3.1.0"},
	// 	{20000001, "2.0.1"},
	// 	{10150035, "1.15.35"},
	// 	{10000000, "1.0.0"},
	// }

	// var version string
	// var number int64
	// var err error

	// for _, test := range tests {
	// 	version = NumberToVersion(test.number)
	// 	assert.Equal(t, test.version, version)
	// 	number, err = VersionToNumber(test.version)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, test.number, number)
	// }
}

// func NumberToVersion(number int64) string {
// 	major := number / 10000000
// 	minor := int64(math.Mod(float64(number), 10000000)) / 10000
// 	patch := int64(math.Mod(float64(number), 10000))
// 	version := fmt.Sprintf("%d.%d.%d", major, minor, patch)
// 	return version
// }

// func VersionToNumber(version string) (int64, error) {
// 	s := strings.Split(version, ".")
// 	var number int64
// 	major, err := strconv.ParseInt(s[0], 10, 64)
// 	if err != nil {
// 		return 0, err
// 	}
// 	minor, err := strconv.ParseInt(s[1], 10, 64)
// 	if err != nil {
// 		return 0, err
// 	}
// 	patch, err := strconv.ParseInt(s[2], 10, 64)
// 	if err != nil {
// 		return 0, err
// 	}
// 	number = (major * 10000000) + (minor * 10000) + patch
// 	return number, nil
// }

func TestOtaContent(t *testing.T) {
	// db := connectDB()
	// defer db.Close()
	// pdb, err := postgresql.New(db)

	// clientID := int64(2)
	// vendingID := int64(2)
	// otp := vending.Otp{}

	// // content := vending.OtaContent{}
	// // cm, err := content.OtaContent(pdb, clientID, vendingID, otp)
	// // if err != nil {
	// // }
	// // log.Println(cm)

	// //  get Client Meta
	// lccommand := "select id,meta from clients where id =  " + strconv.FormatInt(clientID, 10)
	// fmt.Println(lccommand)
	// //fmt.Println("pid = ", clientID)
	// rs, err := pdb.QueryRow(lccommand)
	// cm := vending.ClientMetaDB{}
	// if err != nil {
	// }
	// var _meta string
	// var _id int64
	// rs.Scan(&_id, &_meta)
	// // log.Println("xxxxx Client Id xxxxx ", _id)
	// // log.Println("xxxxx Client Meta xxxxx ", _meta)
	// //stringMeta := _c.Meta
	// _ = json.Unmarshal([]byte(_meta), &cm)
	// // fmt.Println("client meta object ", cm)
	// //fmt.Println(m.ClientConfig.Logo)

	// // GEN VENDING META
	// lccommand = "select id,meta from vendings where id =  " + strconv.FormatInt(vendingID, 10) + " and client_id = " + strconv.FormatInt(clientID, 10)
	// rs, err = pdb.QueryRow(lccommand)
	// vm := vending.VendingMetaDB{}
	// if err != nil {
	// }
	// var _vMeta string
	// var _vID int64
	// rs.Scan(&_vID, &_vMeta)
	// // log.Println("xxxxx Vending ID xxxxx ", _vID)
	// // log.Println("xxxxx Vending Meta xxxxx ", _vMeta)

	// _ = json.Unmarshal([]byte(_vMeta), &vm)
	// // log.Println(err.Error())
	// // fmt.Println("vending Meta Object ", vm)

	// // Get vending access token

	// // // gen Host
	// host := vending.Host{}
	// // host.VendingID = vendingID
	// // //host.XAccessToken = ""
	// // host.XAccessToken = ""
	// // host.MockXAccessToken = ""

	// res := vending.OtaContent{
	// 	Time_servers:      cm.Time_servers,
	// 	Otp:               otp,
	// 	Host:              host,
	// 	Client:            cm.Client,
	// 	Hardware:          vm.Hardware,
	// 	Cash:              cm.Cash,
	// 	CashAcceptedValue: cm.CashAcceptedValue,
	// 	ChangeRefillValue: cm.ChangeRefillValue,
	// }

	// log.Println(res)
	// // log.Println(vm.Hardware)
	// // log.Println(cm)
}

func TestListUpdateByVending(t *testing.T) {

}
func TestShiftClose(t *testing.T) {

}
func TestShiftDeposit(t *testing.T) {

}
func TestOta(t *testing.T) {

}
func TestOtaDone(t *testing.T) {

}
func TestImageList(t *testing.T) {

}

*/
