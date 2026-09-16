# Auction House

## HTTP API

All monetary values are integer cents. For example, `1250` means `12.50`.
Timestamps are encoded as RFC 3339 strings in UTC.

### Create an auction

`POST /auctions`

The JSON body contains `productName`, `startingPrice`, `minIncrement`, and `buyoutPrice`.

### Get an auction

`GET /auctions/{id}`

### Place a bid

`POST /auctions/{id}/bids`

The JSON body contains `amount`, expressed in cents.

Malformed input returns `400`, a missing auction returns `404`, and a bid that conflicts with the current auction state returns `409`.
