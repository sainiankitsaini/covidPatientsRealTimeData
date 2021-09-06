# RealTime Covid Patients State Specific Data using reverse Geoencoding 

## Specifications
1. Language: Golang
2. Web Framework: Echo
3. Database: Mongo 
4. Used publicly available APIs for fetching Covid-19 data and reverse geocoding.


## General Description

General Description:
Build an API that fetches the number of Covid-19 cases in each state and in India
and persist it in MongoDB.


Using this data, build an API that takes the user's GPS coordinates as input and
returns the total number of Covid-19 cases in the user's state and in
India(assume India specific coordinates only) and the last update time of data.


## Installation
Clone Repo and Install Other modules from command Below

```python
go get "go.mongodb.org/mongo-driver/bson"
go get "go.mongodb.org/mongo-driver/mongo"
go get "go.mongodb.org/mongo-driver/mongo/options"
go get "go.mongodb.org/mongo-driver/mongo/readpref"

```
## Screenshots

![Alt text](/goProject/images/china.png?raw=true "Outside India Coordinates")
![Alt text](/goProject/images/haryana.png?raw=true "Default Coordinates")
![Alt text](/goProject/images/up.png?raw=true "Coordinates Input Manually")


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
