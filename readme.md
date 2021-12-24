# Fair Price

>Знал бы прикуп, жил бы в Сочи (a russian mems)

The application is calculated "Fair Price" using the data from source for the last 3 minutes. 


## User story

Task description located in the file doc/Index Price Test Task.pdf  


## Specification (С4)

Description of conditions and restrictions are implemented in the file doc/c4.ru.md


## Usage

run app as:
> go run cmd/fairprice/main.go

## Settings
 - emulated data source are in the file pub_mock_data.go
 - calculation emulation interval can be increased up to 5 seconds by changing the constant calcInterval в main.go

## Project structure
- cmd/fairprice - the main file
- doc - documentation
- internal/eventsource - Message emulator with data source
- internal/exchangeapi - FairPriceCalculator is the main entity
- internal/pricecalcmodels - "fairprice" calculaion model
- internal/tools - app utils
- vendor - folder witj external dependency

## Features of implementation

1. Calculate FairPrice only according to the course data is not possible.
2. Fairprice is a predicted price for the next interval in this case.
3. The model is based on the method of linear approximation by three points
4. The application is created as a reason for the discussion, and is not intended for work in the production




