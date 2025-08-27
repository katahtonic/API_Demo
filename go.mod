module locationVerify

go 1.23.5

require API_Demo/handlers v0.0.0

require (
	API_Demo/models v0.0.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)

replace (
	API_Demo/handlers => ./handlers
	API_Demo/models => ./models
)
