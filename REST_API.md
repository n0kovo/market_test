# Tochka Free Market REST API

Rest API handles use same endpoints as regular site but with ?json=true added to query path.

## Captcha

Use captcha_id from both responses to display captcha: 

http://iyf3xi4gbq4tkfzjzlm725awhjbnajc7fp7gb5umcehmbt2mxq2eyfid.onion/captcha/:captcha_id

## Login 

Endpoint: /api/auth/login

### GET Response

Example:
	
	> curl http://iyf3xi4gbq4tkfzjzlm725awhjbnajc7fp7gb5umcehmbt2mxq2eyfid.onion/api/auth/login
	
	{"captcha_id":"jVnRGZ9sQHzaJVPiv9mY"}

Captcha id must be used to obtain captcha image on /captcha/:captcha_id route

### POST Response

Example:
	
	> curl --data "username=test&passphrase=test&captcha_id=CumVFfRdnwEr7VPPbK6s&captcha=6207" http://iyf3xi4gbq4tkfzjzlm725awhjbnajc7fp7gb5umcehmbt2mxq2eyfid.onion/api/auth/login
	
	{
		"api_session": {
			"token": "a6949a9fa5f24de272fbe1e0cab9cdf6",
			"end_date": "2018-03-23T22:33:47.259373345+03:00",
			"is_2fa_session": false,
			"is_2fa_completed": false
		}
	}

**All subsequent requests must include token get parameter**

## Register

### GET Response

Example:
	
	> curl http://iyf3xi4gbq4tkfzjzlm725awhjbnajc7fp7gb5umcehmbt2mxq2eyfid.onion/api/auth/register
	
	{"captcha_id":"jVnRGZ9sQHzaJVPiv9mY"}

Captcha id must be used to obtain captcha image on /captcha/:captcha_id route

### POST Response

Example:
	
	> curl --data "username=test&passphrase_1=test&passphrase_2=test&role=seller&captcha_id=CumVFfRdnwEr7VPPbK6s&captcha=6207" http://iyf3xi4gbq4tkfzjzlm725awhjbnajc7fp7gb5umcehmbt2mxq2eyfid.onion/auth/login?json=true
	
	{"captcha_id":"WxA3CXmhiUYs8CJnvSVe","error":"Invalid captcha"} 

## Items SERP Response

### GET Request

Example 

	> curl http://iyf3xi4gbq4tkfzjzlm725awhjbnajc7fp7gb5umcehmbt2mxq2eyfid.onion/api/serp?token=a6949a9fa5f24de272fbe1e0cab9cdf6

	{
		"page": 1,
		"number_of_pages": 481,
		"sort_by": "popularity",
		"city": "All countries",
		"geo_cities": [
			{
				"geonameid": 0,
				"name": "",
				"country": "",
				"subcountry": ""
			},
			...
			{
				"geonameid": 6174041,
				"name": "Victoria",
				"country": "Canada",
				"subcountry": "British Columbia"
			}
		],
		"shipping_from_list": [
			"",
			"Afghanistan",
			"Aland Islands",
			...
			"United Kingdom",
			"United States",
			"Uruguay",
			"Worldwide"
		],
		"shipping_to_list": [
			"",
			"Afghanistan",
			"Andorra",
			...
			"United Arab Emirates",
			"United Kingdom",
			"United States",
			"Vatican",
			"Western Sahara",
			"Worldwide",
			"Yemen"
		],
		"account": "all",
		"available_items": [
			{
				"item_uuid": "cc2bfb67891d422e44415d94db75b492",
				"vendor_uuid": "f268db63a6054352705808cda6042e14",
				"vendor_username": "joshkingseller",
				"vendor_description": "Hello everyone I'm here to make u rich  ! No escrow on my product",
				"vendor_language": "en",
				"vendor_is_premium": true,
				"vendor_is_premium_plus": false,
				"vendor_is_trusted": false,
				"type": "digital",
				"item_created_at": "2017-08-18T05:50:22.34331+03:00",
				"item_name": "xxxx",
				"item_description": "xxxx",
				"item_category_id": 25,
				"item_parent_category_id": 23,
				"item_parent_parent_category_id": 0,
				"vendor_score": 4.57,
				"vendor_score_count": 14784,
				"item_score": 4.59,
				"item_score_count": 14784,
				"country_shipping_from": "Interwebs",
				"country_shipping_to": "Interwebs",
				"geoname_id": 0,
				"vendor_is_online": false,
				"vendor_last_login_date": "2 days ago",
				"vendor_registration_date": "8 months ago",
				"price_range": [
					"15",
					"50"
				],
				"price": "",
				"vendor_btc_tx_number": "100+",
				"vendor_btc_tx_volume": "0.5-1 BTC",
				"item_btc_tx_number": "50-100",
				"vendor_eth_tx_number": "10-20",
				"item_eth_tx_number": "5-10"
			}
		]
	}