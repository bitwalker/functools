// functools provides high level functions for calling and manipulating functions
package functools

import (
	"reflect"
)

// Anything: represents any possible type
type Anything interface{}

// Function: represents any possible function that returns a single value
type Function func(...Anything) Anything

// MultiFunction: represents any possible function that returns two values
type MultiFunction func(...Anything) (Anything, Anything)

/*
   Apply is used to partially apply function arguments to
   produce a function with smaller arity. This can lead to
   clearer, more intuitive code that builds on the benefits
   of DRY.

   Example:
       // Say we have a basic function Add which sums two arguments
       func Add(x, y int) int {
           return x + y
       }

       // And say we have a lot of code that simply adds 1 to another argument
       // We could simply call add(1, num) in all of those places, or, we could
       // define a new function Increment, which does this for us. Increment 
       // could be defined in a number of ways, but the DRYest makes use of Apply:
       var Increment = Apply(Add, 1)

       // So which of these do you feel is the most intuitive and easy to read?
       z := Add(1, 10)
       // or..
       z := Increment(10)
*/
func Apply(f Anything, args ...Anything) Function {
	// In order to work with any function type, we have to box it
	// in Anything, and extract the true value using reflection.
	fn := reflect.ValueOf(f)

	// We return a function which takes any number of additional arguments (0..N),
	// which when called will call the original function with all of the arguments
	// aggregated.
	var applied Function
	applied = func(moreargs ...Anything) Anything {
		// Aggregate the two sets of arguments and extract
		// the original argument values (arguments is []Anything)
		values := AnythingToValues(append(args, moreargs...))
		// Call the function using reflection, and return the value boxed as Anything
		result := fn.Call(values)
		val := result[0].Interface()
		return val
	}

	return applied
}

/*
   ApplyMulti performs the same function as Apply, but does it for
   functions with multiple return values. The behavior is more or
   less identical, so comments are stripped, as the differences should
   be self-explanatory.
*/
func ApplyMulti(f Anything, args ...Anything) MultiFunction {
	fn := reflect.ValueOf(f)

	var applied MultiFunction
	applied = func(moreargs ...Anything) (Anything, Anything) {
		values := AnythingToValues(append(args, moreargs...))
		// The convention with most multiple return functions is to store
		// the value of the operation in the first value, and the error, if
		// any, in the second. I've named the variables accordingly here, but
		// be aware that the values could really be any combination of two types.
		result := fn.Call(values)
		val := result[0].Interface()
		err := result[1].Interface()
		return val, err
	}

	return applied
}

/*
   Compose takes two functions, f1 and f2, and returns a new function
   that when called, applies it's arguments to f2, then applies the
   result as a single argument to f1, and then returns the result.

   Example:
       func Add(a, b int) int {
           return a + b
       }
       func Square(x int) int {
           return x * x
       }

       var SquareSum = Compose(Square, Add)

       SquareSum(3, 3) // => 36
*/
func Compose(f1 Anything, f2 Anything) Function {
	fn1 := reflect.ValueOf(f1)
	fn2 := reflect.ValueOf(f2)

	var composed Function
	composed = func(args ...Anything) Anything {
		values := AnythingToValues(args)
		inside := fn2.Call(values)[0].Interface()
		result := fn1.Call([]reflect.Value{reflect.ValueOf(inside)})[0].Interface()
		return result
	}

	return composed
}

/*
   AnythingToValues is used to return a slice of reflected values
   for a slice of type Anything (which is really just interface{})
*/
func AnythingToValues(items []Anything) []reflect.Value {
	values := make([]reflect.Value, len(items))
	for k, v := range items {
		values[k] = reflect.ValueOf(v)
	}
	return values
}

// LinkedList is simply a pointer to a function which will return the first Node
type LinkedList func() *Node

// Empty denotes the end of the list. It is a Thunk which returns nil.
var Empty *LinkedList

/*
   Every List is composed of Nodes. Each node contains a Head, or the
   current value of this node; and a Tail, which is just another List
*/
type Node struct {
	Head Anything
	Tail *LinkedList
}

/* 
   Creates a LinkedList from a head element and a tail Thunk, this is used
   just like the `cons` operator in Lisp. You can chain Cons to build
   a list, though it is quite verbose:

   Example:
       list := Cons("A", Cons("B", Cons("C", Empty)))
*/
func Cons(head Anything, tail *LinkedList) *LinkedList {
	var list LinkedList
	list = func() *Node {
		return &Node{head, tail}
	}
	return &list
}

/*
   Create a List from the provided arguments (or a slice using the ... syntax)

   Example:
       nums  := [...]int{1, 2, 3}
       list  := List(1, 2, 3) // => [1, 2, 3]
       list2 := List(nums...) // => [1, 2, 3]
*/
func List(elements ...Anything) *LinkedList {
	result := Empty
	if len(elements) > 0 {
		result = Cons(elements[0], List(elements[1:]...))
	}
	return result
}

/*
   Gets the length of the List. Calling this on an infinite list
   will cause an endless loop. Care is required!
*/
func (list *LinkedList) Length() int {
	length := 0
	node := (*list)()
	for node != nil {
		if node.Tail != nil {
			node = (*node.Tail)()
		} else {
			node = nil
		}
		length++
	}
	return length
}

/*
   Converts a slice of any type to a LinkedList

   Example:
       nums := [...]int{1, 2, 3}
       list := ToList(nums)

*/
func ToList(elements Anything) *LinkedList {
	sliceType := reflect.TypeOf(elements)
	result := Empty
	if elements == nil || sliceType.Kind() != reflect.Slice {
		panic("Attempted to call ToList on a value of the wrong type. Must be Slice.")
	} else {
		val := reflect.ValueOf(elements)
		// Build the list in reverse
		for i := val.Len() - 1; i >= 0; i-- {
			result = Cons(val.Index(i).Interface(), result)
		}
	}
	return result
}

/*
   Converts a LinkedList to []Anything
*/
func ToSlice(list *LinkedList) []Anything {
	result := make([]Anything, list.Length())
	node := (*list)()
	for i := 0; node != nil; i++ {
		result[i] = node.Head
		if node.Tail != nil {
			node = (*node.Tail)()
		} else {
			node = nil
		}
	}
	return result
}

/*
   Returns a new LinkedList containing the first N elements.
*/
func (list *LinkedList) Take(n int) *LinkedList {
	var taken LinkedList
	taken = func() *Node {
		if n > 0 {
			node := (*list)()
			if node != nil {
				return &Node{node.Head, node.Tail.Take(n - 1)}
			}
		}
		return nil
	}
	return &taken
}

/*
   Returns a new LinkedList with the first n elements dropped.
*/
func (list *LinkedList) Drop(n int) *LinkedList {
	var remaining LinkedList
	remaining = func() *Node {
		node := (*list)()
		if node != nil {
			if n > 0 {
				n--
				list = node.Tail
				return remaining()
			}
			return node
		}
		return nil
	}
	return &remaining
}

/*
   Maps a function to each element of a list. This is a lazy operation.

   Example:
       list := List(1, 2, 3)
       squared := list.Map(func(x int) int { return x * x })
*/
func (list *LinkedList) Map(f Anything) *LinkedList {
	expr := reflect.ValueOf(f)
	var mapped LinkedList
	mapped = func() *Node {
		node := (*list)()
		if node != nil {
			args := []reflect.Value{reflect.ValueOf(node.Head)}
			head := expr.Call(args)[0].Interface()
			tail := Empty
			if node.Tail != nil {
				tail = node.Tail.Map(f)
			}
			return &Node{head, tail}
		}
		return nil
	}
	return &mapped
}

/*
   Reduces the elements of a list to a single value.

   Example:
       list := List(1, 2, 3)
       sum := list.Reduce(func(acc, x int) int { return acc + x }, 0) // => 6
*/
func (list *LinkedList) Reduce(f Anything, memo Anything) Anything {
	expr := reflect.ValueOf(f)
	node := (*list)()
	for node != nil {
		args := []reflect.Value{reflect.ValueOf(memo), reflect.ValueOf(node.Head)}
		memo = expr.Call(args)[0].Interface()
		if node.Tail != nil {
			node = (*node.Tail)()
		} else {
			node = nil
		}
	}
	return memo
}
