/*
Provider for locating active Deliveroo orders.

The `order_index` field allows selecting a preference for which order should be located, if multiple are active at once.

It should be the (zero-indexed) index of the order desired, when sorted by earliest expected delivery time.

0 will select the earliest expected delivery (default)
-1 will select the last expected delivery
A positive number higher than the number of active orders is equivalent to -1
A negative number higher than the number of active orders is equivalent to 0

	{ "type": "deliveroo" }

	{ "type": "deliveroo", "order_index": -1 }
*/
package deliveroo
