package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// #key generated from account https://developer.here.com/blog/reverse-geocoding-a-location-using-golang
//key 1
// var apikey = "eiijxY1_PtsTt52GKtFJZeNZK9ZzQAxV6NtwC8YlkTA"
// key 2
var apikey = "L_bktvTy1Z2aL2pD5KzUmbgJB4cxMF-DyfzGBHl1Nmw"

var latitude = 28.7975
var longitude = 76.1322
var database = "covid"
var collection = "states"

type State struct {
	Name        string
	CountryCode string
	Deaths      int
	Province    string
	CityCode    string
	Active      int
	Country     string
	City        string
	Lon         string
	Date        string
	Lat         string
	Recovered   int
	Confirmed   int
	ID          string
}

type Results struct {
	Items []Item
}

// nested within sbserver response
type Item struct {
	Title   string
	Address struct {
		CountryCode string `json:"countryCode"`
		CountryName string `json:"countryName"`
		StateCode   string `json:"stateCode"`
		State       string `json:"state"`
		County      string `json:"county"`
		City        string `json:"city"`
		PostalCode  string `json:"postalCode"`
	}
}

// query method returns a cursor and error.
func query(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {

	// select database and collection.
	collection := client.Database(dataBase).Collection(col)

	// collection has an method Find,
	// that returns a mongo.cursor
	// based on query and field.
	result, err = collection.Find(ctx, query,
		options.Find().SetProjection(field))
	return
}

// / query method returns a cursor and error.
func latestDocFind(client *mongo.Client, ctx context.Context,
	dataBase, col string, province string, query, field interface{}) (result *mongo.Cursor, err error) {

	// select database and collection.
	collection := client.Database(dataBase).Collection(col)

	filter := bson.D{
		{"province", province}}

	//  option remove id field from all documents
	options := options.Find()

	options.SetSort(bson.D{{"date", -1}})
	options.SetLimit(1)

	// collection has an method Find,
	// that returns a mongo.cursor
	// based on query and field.
	result, err = collection.Find(ctx, filter,
		options)

	return
}

// / This is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// This is a user defined method that returns mongo.Client,
// context.Context, context.CancelFunc and error.
// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and
// resource associtated with it.

func connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

func UpdateOne(client *mongo.Client, ctx context.Context, dataBase,
	col string, filter, update interface{}) (result *mongo.UpdateResult, err error) {

	// select the databse and the collection
	collection := client.Database(dataBase).Collection(col)

	// A single document that match with the
	// filter will get updated.
	// update contains the filed which should get updated.
	result, err = collection.UpdateOne(ctx, filter, update)
	return
}

// This is a user defined method that accepts
// mongo.Client and context.Context
// This method used to ping the mongoDB, return error if any.
func ping(client *mongo.Client, ctx context.Context) error {

	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occored, then
	// the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

func search2(province string, client *mongo.Client, ctx context.Context) {
	// create a filter an option of type interface,
	// that stores bjson objects.
	// var filter, option interface{}

	// db.getCollection("states").find({}).sort({"date": -1}).limit(1);
	filter := bson.D{
		{"province", province}}

	//  option remove id field from all documents
	options := options.Find()

	// call the query method with client, context,
	// database name, collection  name, filter and option
	// This method returns momngo.cursor and error if any.
	// cursor, err := query(client, ctx, database,
	// 	collection, filter, option)

	cursor, err := latestDocFind(client, ctx, database,
		collection, province, filter, options)

	// handle the errors.
	if err != nil {
		panic(err)
	}

	var results []bson.D

	// to get bson object  from cursor,
	// returns error if any.
	if err := cursor.All(ctx, &results); err != nil {

		// handle the error
		panic(err)
	}

	// printing the result of query.
	fmt.Println("Query Reult")
	for _, doc := range results {
		fmt.Println(doc)
	}
}

func getSanitizedDates() (string, string, string) {
	now := time.Now()
	after := now.AddDate(0, 0, -1)

	day := after.Day()
	month := after.Month()
	year := after.Year()

	yearStr := strconv.Itoa(year)
	dayStr := strconv.Itoa(day)
	if day < 10 {
		dayStr = "0" + dayStr
	}
	monthInt := int(month) // normally written as 'i := int(m)'
	monthStr := strconv.Itoa(monthInt)
	if monthInt < 10 {
		monthStr = "0" + monthStr
	}
	return dayStr, monthStr, yearStr
}

func getCovidData() []byte {
	dayStr, monthStr, yearStr := getSanitizedDates()
	url := "https://api.covid19api.com/live/country/india/status/confirmed/date/" + (yearStr) + "-" + (monthStr) + "-" + (dayStr) + "T00:00:00Z"
	fmt.Println("Fetching Data from Public Api : ")
	fmt.Println(url)

	response, err := http.Get(url)

	// response, err := http.Get("https://api.covid19api.com/live/country/india/status/confirmed/date/2021-08-31T13:13:30Z")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(responseData))
	return responseData
}

func makeUrlForAddrFetch() string {
	fmt.Println("Enter latitude: ")

	// var then variable name then variable type
	var lat string

	// Taking input from user
	fmt.Scanln(&lat)
	fmt.Println("Enter  longitude: ")
	var lon string
	fmt.Scanln(&lon)

	// Print function is used to
	// display output in the same line
	fmt.Print("latitude  and longitude are: ")

	// Addition of two string
	fmt.Print(lat + " " + lon)
	url := "https://revgeocode.search.hereapi.com/v1/revgeocode?apiKey=" + apikey + "&at=" + fmt.Sprint(lat) + "," + fmt.Sprint(lon)
	return url
}

func getCompleteAddr() []byte {

	fmt.Println("Press Y to input lat and long manually: ")

	// var then variable name then variable type
	var manualInput string

	// Taking input from user
	var url string
	fmt.Scanln(&manualInput)
	if manualInput == "Y" || manualInput == "y" {
		url = makeUrlForAddrFetch()
	} else {
		url = "https://revgeocode.search.hereapi.com/v1/revgeocode?apiKey=" + apikey + "&at=" + fmt.Sprint(latitude) + "," + fmt.Sprint(longitude)
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Reverse GeoEncoding Result: ")

	fmt.Println(string(body))
	return body
}

func getStateOnCoordinateBasis() string {
	body := getCompleteAddr()

	//output of above api
	// {"items":[{"title":"Canara Bank-Bhiwani","id":"here:pds:place:356ttn5u-3f8720fb36f04a1d80fa8ae098fcacd6","resultType":"place","address":{"label":"Canara Bank-Bhiwani, SH-17, Krishna Colony, Bhiwani 127021, India","countryCode":"IND","countryName":"India","stateCode":"HR","state":"Haryana","county":"Bhiwani","city":"Bhiwani","district":"Krishna Colony","street":"SH-17","postalCode":"127021"},"position":{"lat":28.79749,"lng":76.13226},"access":[{"lat":28.79754,"lng":76.13236}],"distance":6,"categories":[{"id":"700-7000-0107","name":"Bank","primary":true},{"id":"700-7010-0108","name":"ATM"}]}]}

	apiRes := &Results{}
	err1 := json.Unmarshal([]byte(body), apiRes)
	if err1 != nil {
		log.Fatal(err1)
	}

	// fmt.Printf("%v\n", apiRes)
	// fmt.Printf("%v\n", apiRes.Items[0].Address.State)
	if apiRes == nil || apiRes.Items == nil || len(apiRes.Items) == 0 {
		return ""
	}
	return apiRes.Items[0].Address.State
}

func storeCovidDataInMongo(responseData []byte, client *mongo.Client, ctx context.Context) {
	var states []State
	json.Unmarshal([]byte(responseData), &states)

	for i := 0; i < len(states); i++ {
		fmt.Println("States %+v", states[i])
	}

	// Ping mongoDB with Ping method
	ping(client, ctx)

	collection := client.Database(database).Collection(collection)
	for i := 0; i < len(states); i++ {
		insertResult, err := collection.InsertOne(context.TODO(), states[i])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a Single Document: ", insertResult.InsertedID)
	}
}

func main() {

	responseData := getCovidData()
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	// Release resource when the main
	// function is returned.
	defer close(client, ctx, cancel)

	storeCovidDataInMongo(responseData, client, ctx)

	province := getStateOnCoordinateBasis()
	if province == "" {
		fmt.Println("No state Found on basis of input coordinates")
		return

	}

	fmt.Println("State got for finding details on the basis of latitude and longitude:")
	fmt.Println(string(province))
	search2(province, client, ctx)

}
