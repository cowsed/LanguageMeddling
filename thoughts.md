We really dont know how to compute
https://www.youtube.com/watch?v=8Ab3ArE8W3s

```
1 + 1 2 3 4   =>    2 3 4 5
syntax pretty neat, maybe some brackets
1 + [1 2 3 4] //makes it easer to see that its an array
res = 1 + [1 2 3]
res = [1 2 3] + [1 2 3] => [2 4 6]
res = [1 2] + [1 2 3] => throws error, cant add two differently sized vectors
```

vector (slice under the hood) (benifits of static array for most of the time, can be extended with minimal effort)

res = [1 2]:[3 4] = [1 2 3 4]
res = [1 2 3][1] = 2
res = [1 2 3][0:2] = [1 2]

func __concat__(left, right vector) vector {
    
}
func __equals__(left, right vector) bool {
    //for each item, check if theyre equal, if different sizes, unequal
}

func __to_float__(left number){

}


----
int a = 2
float b = float(a) // calls a__to_float__()


a = 2 //last expression = Nothing
2 //last expression = 2