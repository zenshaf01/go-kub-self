# go-kub-self
The project is an attempt to learn building production grade applications using go and kubernetes

# Design Philosophy, Guidelines:
programming mode is trying to get something to build and work (even though if, it doing dirty programming)
- We do not make things easy to do, we make things easier to understand. (Don't always stop at programming code, go into engineering mode (consider performance, logging, error handling))
    - Code should be (By anyone):
        - maintainable
        - manageable
        - debuggable
- Every encapsulation must define a new semantic. This means that if you introduce new variables or functions should introduce something new.
- Engineering should have clear layers of concern and purpose.
