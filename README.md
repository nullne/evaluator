#### Expression syntax
We use s-epxression syntax to parse and evaluate.
> In computing, s-expressions, sexprs or sexps (for "symbolic expression") are a notation for nested list (tree-structured) data, invented for and popularized by the programming language Lisp, which uses them for source code as well as data. 


#### How to
The simplest usage:

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

#### Functions
##### Implemented functions
##### How to use self-defined functions