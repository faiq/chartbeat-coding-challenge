#Chartbeat Coding Challenge

## The Challenge
Implement a prototype web service that exposes an endpoint and responds with JSON. It only
needs to work for the top 100 pages in a domain, but should work for any domain for which you
have data. The service should also service at least 100 requests per second. The API call to get
the top pages for a domain is guaranteed to return within 100 milliseconds.

## My Approach
In order to tackle this problem, I broke down to several parts:

1. A function that makes requests to the chartbeat API
2. A routine that takes data from the chartbeat API and saves them for later processing
3. A webserver with a route that takes a parameter `host` and determines whether or not to start polling that host from the chartbeat API
4. A routine that schedules polling a given host on a given interval
5. A function that takes the state of the application and returns increasing page paths for that host

And tried to come up with solutions for each peice

## My Soultions for These Subproblems
1. A function that makes requests to the chartbeat API
- To solve this, I wrote a function that takes in a given url for a host and a golang channel. I then make a request to the given host, stream over the results for a given path, and send them over a channel to a routine that saves them for later processing.

2. A routine that takes data from the chartbeat API and saves them for later processing 
- To solve this, I started a routine that listens on a channel that sends `Page` structs. When a struct is sent from routine 1, I call a new function `HandlePage` which looks to see if the path for that `Page` struct exists in the global memory map known as `state`. If it exists, it updates the `Visitors` and `PrevVisitors` feilds for that page struct in `state`. Otherwise, it will add this `Page` object to `state`.

The `state` object looks like this

```json
	{
		"gizmodo":[
			{ 
				"I": "something cool",	
				"Path": "/cool-article-or-something",
				"Visitors": 1,
				"PrevVisitors": 0
			}
		]
```

3. A webserver with a route that takes a parameter `host` and determines whether or not to start polling that host from the chartbeat API
- To do this I just set up a basic webserver and router using negroni and gorilla/mux, and set up one route which extracts a host parameter from the URL. If the host is already being polled, it will return increaising paths for that host	
4. A routine that schedules polling a given host on a given interval
- To solve this problem, I wrote the function `Poll` which returns starts a go routine that will call the function that makes requests to the chartbeat API whenever the poll interval is hit. It returns the channel that the `Page` structs come through

5. A function that takes the state of the application and returns increasing page paths for that host
- This involved me looping over all of the `Page` structs comparing thier `PrevVisitors` and `Visitors` fields and calculating the difference between the two, and returning all of the increasing numbers.

## Bonus
> Discuss why this is bad and implement a better one that is more useful to a potential Chartbeat user.

