
# Language Specification
## statements
### var
```go
var a int //default initialized to 0
var b int = 5
var c string = "wow"

var v vec<int> = [a b]
var v2 vec<int> = [a b c] //illegal beacause c is a string
var v4 vec<int> = [a 2 b]
var z tuple = <a b c> //just fine since tuples can be anything

var f func(int, int)int = func(a int, b int)int{return a+b}
```
### const
```go
const the_answer int = 42
const arr vec<int> = [1 2 3 4 5] //compiler will yell at you if you try to assign to this or try to assign to a sub element
```
### print

```go
print a
print (a "on the left, on the right" b)
```

### return
```go
return a
return (a b c)
```

### assignment
```go
a = 2+3
```

## literals
### integer literal
```go
12
0
-13
```
### float literal
IEEE-754 float
```go
-120
12.2
0
-1e4
```
### vec literal
all the same type
```go
[1 2 3 4]
[a b] //assuming a and b are ints

```
### tuple literal
different types. for usage examples see print, multiple returns
```
<a, 4, "wow">
<3 1 "a">
```
## flow

### function call
```
f()
f(a b)
f(a, b)
f(t...)
```
### function definition
```go
func f(a int, b int)int{
    return a+b
}
func f2(a int, b int) tuple(int, float) {
    return (a+b float(a+b))
}
// return type assumed to be a tuple if the return type is bracketed
func f2(a int, b int) <int, float> {
    return <a+b float(a+b)>
}
```
### if/else statement
```go
if boolean_expression{

} elif boolean_expression 2{

} else {

}
```
### for
a for loop

```go
//  statement 1 must be variable declaration
// boolean expression must be a boolean expression using variable declared in statement 1
// statement 2 is a statment that has to do with variable declared in statement 1
for (statement 1; boolean_expression; statement 2){
    //do stuff
}
for (var i int = 0; i<12; i++){

}
```
### while
```
while (arbitrary_boolean_expression){

}
```
### solve
``` go
solve universe_of_discourse, name_of_solution_set{
    require universe_of_discourse[i] true

}
```
### option
If not executing in a solve block, chooses the first (first is default)
If executing in a solve block the solver takes both as options

## Type erasure

type erasure can be accomplished by wrapping in a tuple. 
operator overloads + wrapping means you could write 
```
func min(a, b <>) <> {
    if a<b{
        return a
    }
    return b
}
``` 
and it would work for any floats, ints or other things that have ```<``` defined on them
