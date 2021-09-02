package main

import ( "fmt" )

// A function declaration: Types after variable names; Multiple return values
func multi(x int, y int) (string, int) {
     return "The answer is", x+y
}

// A class
type Named struct {
     name string // Member of a class
}

// Method on a class
func (nm *Named) PrintName(n int) {
    if n < 0 { panic(-1) } // an "exception"
    for i:=0; i<n; i++ {
        fmt.Print(nm.name + "\n")
    }
}


func main() {
    var i int // Explicit declaration of variable
    i = 21    // = is used for assignment
    j := 21   // := declares the variable based on the type of the value

    decr := func() int { // a closure
        j = j-7
        return j
    }

    str, m := multi(i, j)
    defer fmt.Println(str, m) // run after return or a panic

    fmt.Println(decr())
    fmt.Println(decr())

    nm1 := Named{ name: "Jesper" } // value
    nm2 := &Named{} // pointer

    nm1.PrintName(2) // Calling a method
    nm2.PrintName(2) // . does auto-dereference of pointer
    nm2.name = "Claudio" ; nm2.PrintName(2) // RETURN or ; as seperator
    var nm3 *Named = nm2; // how to explicitly declare a pointer
    nm2.name = "Ivan"
    nm3.PrintName(2)
    nm3.PrintName(-1) // let us cause a panic
    fmt.Println("Will we make it here?")
}
