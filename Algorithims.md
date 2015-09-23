# The challenge 
	
	generate a list of the top 10 items sold for the month up to the current time. You should write three functions in pseudo-code:

	reset -- Resets all the counts to zero and is called automatically at the start of every month.
	inc_id(item id) -- Updates the number of times the given item id has been purchased by 1. This is called the instant a purchase completes.
	get_top_10() -- Returns the top 10 items and number of times they’ve been purchased since the last call to reset.

In order to solve this I will use a hash table called `my_map` which uses `item_id` as the key. I choose this data structure because of its O(1) look up tiem, which will be usefull, especially on a busy ecommerce site where users are constantly buying things. This will also take O(n) space.

This however makes the reset and get_top_10 methods more computationally expensive.

In the case of the reset method, we need to loop over all of the keys in the map and delete the keys one by one. This is an O(n) opperation.

The top 10 method also has an O(n) time to run through all the keys in the map, but also has the extra task of creating and computing a top 10 list. To do this, we'll keep track of the smallest item in the top 10 list and replace any item that is bigger than the smallest item in the list we'll pop that item out, and replace it with the new smallest item. This allows us to save a loop over the top ten list, without having to do extra computational steps like sorting.

```python

	func inc_id(item_id) { 
		my_map[item_id]+=1
	}
	
	func reset () {
		for key in my_map:
			del my_map[key]
	}

	func get_top_10 () {
		top_10 = []
		min = None
		for key in my_map: 
			num = my_map[key]
			if min is None:
				min = num
				top10.append(min)
				continue
			if num < min and len(top_10) < 10:
				min = num
				top10.append(min)
				continue
			if num < min and len(top_10) == 10:
				top10.pop(min)
				min = num
				top10.append(min)
				continue
		return top_10	
	}
```
