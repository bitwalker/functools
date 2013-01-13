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
