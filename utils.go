package main

func Contains[E comparable](list []E, item E) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func Remove[E comparable](list []E, item E) []E {
	for i, v := range list {
		if v == item {
			// Remove the element by appending the parts before and after it
			return append(list[:i], list[i+1:]...)
		}
	}
	// Return the original slice if the element is not found
	return list
}
