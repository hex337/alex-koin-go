
## Welcome

Alex Koin is being rewritten in golang. If you are interesteded in helping out, please read on!

## Development

1. `make up`
2. `ngrok http -subdomain=yourdomain 3000`
3. Configure Slack app's event endpoint to call https://yourdomain.ngrok.io/events
4. Add slack settings to .env file (copy .env.template)

### Database Setup

``` psql
CREATE DATABASE akc;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

This should be run with your local `psql` command. The `\copy` function will pull from your local directory.
``` SQL
\COPY
	(SELECT ID,
			EMAIL,
			FIRST_NAME,
			LAST_NAME,
			SLACK_ID,
			INSERTED_AT AS CREATED_AT,
			UPDATED_AT
		FROM USERS) TO 'users.csv' WITH CSV HEADER;

\COPY
	(SELECT COINS.ID,
			COINS.HASH,
			COINS.ORIGIN,
			WALLETS.USER_ID AS USER_ID,
			COINS.MINED_BY_ID AS MINED_BY_USER_ID,
			COINS.CREATED_BY_USER_ID,
			COINS.INSERTED_AT AS CREATED_AT,
			COINS.UPDATED_AT
		FROM COINS
		LEFT JOIN WALLETS ON COINS.WALLET_ID = WALLETS.ID) TO 'coins.csv' WITH CSV HEADER

\COPY
  (SELECT 
    ID,
    AMOUNT,
    MEMO,
    FROM_ID AS FROM_USER_ID,
    TO_ID AS TO_USER_ID,
    COIN_ID,
    INSERTED_AT AS CREATED_AT,
    UPDATED_AT
    FROM TRANSACTIONS
    ) TO 'transactions.csv' WITH CSV HEADER
 ```

 ``` SQL
\COPY USERS (
  ID,
  EMAIL,
  FIRST_NAME,
  LAST_NAME,
  SLACK_ID,
  CREATED_AT,
  UPDATED_AT)
FROM 'users.csv' WITH CSV HEADER;

\COPY COINS (
  ID,
  HASH,
  ORIGIN,
  USER_ID,
  MINED_BY_USER_ID,
  CREATED_BY_USER_ID,
  CREATED_AT,
  UPDATED_AT)
FROM 'coins.csv' WITH CSV HEADER

\COPY TRANSACTIONS (
  ID,
  AMOUNT,
  MEMO,
  FROM_USER_ID,
  TO_USER_ID,
  COIN_ID,
  CREATED_AT,
  UPDATED_AT)
FROM 'transactions.csv' WITH CSV HEADER
```

### Proxy outbound requests

Configure a proxy by setting the env var `http_proxy`, eg `export http_proxy=http://127.0.0.1:9999`. I use Charles proxy for this to see the calls to Slack.
