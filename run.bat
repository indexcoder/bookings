go build -o bookings.exe .\cmd\web\.
bookings.exe -production=false -cache=true -dbHost=localhost -dbName=bookings -dbUser=postgres -dbPass=12345 -dbPort=5432 -dbSsl=disable