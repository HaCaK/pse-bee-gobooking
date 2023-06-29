
# goBooking demo

## API demo via Insomnia

### Bookings

1. Retrieve Bookings => Empty list

2. Get, Update, Delete Booking by Id => Not found

3. Create Booking => Property not found => Create Property first

### Properties

1. Retrieve Properties => Empty list

2. Delete Property by Id 1 => Not found (same for Get and Update by Id)

3. Create Property => has id 1 and status "FREE"
```json
{
	"description": "Family-friendly vacation home in Davenport with water park.",
	"ownerName": "Mickey Mouse",
	"address": "Davenport, Florida",
	"name": "Mansion with pool"
}
```

4. Get Properties => returns property

5. Get Property By Id 1 => returns matching property

6. Update Name of Property 1
```json
{
    "description": "Family-friendly vacation home in Davenport with water park.",
	"ownerName": "Mickey Mouse",
	"address": "Davenport, Florida",
	"name": "Wonderful Mansion near Disneyland"
}
```

7. Delete Property By Id 1 => deletes property

8. Re-create Property
```json
{
	"description": "Family-friendly vacation home in Davenport with water park.",
	"ownerName": "Mickey Mouse",
	"address": "Davenport, Florida",
	"name": "Mansion with pool"
}
```
=> now has id 2 and status "FREE"


### Bookings

1. Create Booking for Property 2
```json
{
	"comment": "We would love to book your amazing property.",
	"customerName": "Dagobert Duck and family",
	"propertyId": 2
}
```
=> has id 1 and status "CONFIRMED"

### Properties

1. Get Property By Id 2 => has status "BOOKED"

2. Create another Booking for Property 2
```json
{
	"comment": "We cannot wait to try out this great place!",
	"customerName": "Goofy and Co",
	"propertyId": 2
}
```
=> declined, because property already "BOOKED"

3. Delete Booking By Id 1
=> Property 2 has status "FREE"

4. Try again to create another Booking for Property 2
```json
{
	"comment": "We cannot wait to try out this great place!",
	"customerName": "Goofy and Co",
	"propertyId": 2
}
```
=> accepted

### Properties

1. Delete Property By Id 2 => not possible, because already booked


## Code

- General code structure (3 Microservices, Booking and Property with known structure, Proxy with gRPC Gateway in gen.go and MUX with Gin in main.go)
- Multi-stage Dockerfile
- Booking tests with mock PropertyInternalServer