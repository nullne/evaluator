It's a common case to evaluate an expression in hot.
#### Expression syntax
We use s-epxression syntax to parse and evaluate.
> In computing, s-expressions, sexprs or sexps (for "symbolic expression") are a notation for nested list (tree-structured) data, invented for and popularized by the programming language Lisp, which uses them for source code as well as data. 


#### How to
You can evaluate directly:

    params := evaluator.MapParams{
        "gender": "female",
    }
    res, err := evaluator.EvalBool(`(in gender ("female" "male"))`, params)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(res)
    # true	
    
or you can reuse the `Expression` to evaluate multiple times:

    params := evaluator.MapParams{
        "gender": "female",
    }
    exp, err := evaluator.New(`(in gender ("female" "male"))`)
    if err != nil {
        log.Fatal(err)
    }
    res, err := exp.EvalBool(params)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(res)
    # true

#### And you can write expressions like this
- `(in gender ("male", "female"))`
- `(between now (type_default_time "2017-01-02 12:00:00") (type_default_time "2017-12-02 12:00:00"))`
- `(ne (mod (age birthdate) 7) 5)`
- or multiple-line for clarity

	```
	(and
		(ne os "ios")
		(eq gender "male")
		(beteen version (type_version "2.7.1") (type_version "2.9.1"))
	)
	```


#### Functions
##### Implemented functions
- in
- between
- and
- or
- not
- equal
- not equal
- greater than
- less than
- greater than or equal to
- less than or equal to
- mod
- plus
- minus
- multiply
- divide
- type version
- type time
- type default time
- type default date


##### How to use self-defined functions
Yes, you can write your own function.
	package main
	
	import (
		"errors"
		"log"
		"time"
	
		"github.com/nullne/evaluator"
		"github.com/nullne/evaluator/function"
	)
	
	type age struct{}
	// define your own function and don't forget to register
	func (f age) Eval(params ...interface{}) (interface{}, error) {
		if len(params) != 1 {
			return nil, errors.New("only one params accepted")
		}
		birth, ok := params[0].(string)
		if !ok {
			return nil, errors.New("birth format need to be string")
		}
		r, err := time.Parse("2006-01-02", birth)
		if err != nil {
			return nil, err
		}
		now := time.Now()
		a := r.Year() - now.Year()
		if r.Month() < now.Month() {
			a--
		} else if r.Month() == now.Month() {
			if r.Day() < now.Day() {
				a--
			}
		}
		return a, nil
	}
	func main() {
		if err := function.Regist("age", age{}); err != nil {
			log.Print(err)
		}
	
		exp := `(not (between (age birthdate) 18 20))`
		vvf := evaluator.MapParams{
			"birthdate": "2018-02-01",
		}
		e, err := evaluator.New(exp)
		if err != nil {
			log.Print(err)
		}
		r, err := e.Eval(vvf)
		if err != nil {
			log.Print(err)
		}
		log.Printf("expression: `%s`, wanna: %+v, got: %+v\r", exp, true, r)
	}
