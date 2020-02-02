# Go error handling with gRPC and HTTP error codes
This is a fork from the error handling example https://github.com/henrmota/errors-handling-example.

Additional to the fork I implemented gRPC and HTTP error codes for simplify error handling in gRPC or HTTP communications. The descriptions of the error codes can be found here: https://cloud.google.com/apis/design/errors#error_model.

This library is wrapped from https://github.com/pkg/errors.

Creating a typed errors is as easy:

```GO
errors.InvalidArgument.New("error parsing the input information")
```

You can create an untyped error as easy as:

```GO
errors.New("an untyped error")
```

Adding a new context to an existing error:

```GO
errors.AddContext(err, "field", "message")
```

In the top layer when you decide to log or return a web response:

```GO
errors.GetType(err) == errors.InvalidArgument // true
errors.GetContext(err) // map[string]string{"field": "field", "message": "message"}
```

Converting error type in gRPC error code:
```GO
errors.OutOfRange.New("an untyped error").Code()
```

or in HTTP code:

```GO
errors.OutOfRange.New("an untyped error").HTTP()
```

To add new error type is just adding a new constant to errors

```GO
const (
  // InvalidArgument error
  InvalidArgument ErrorType = ErrorType(codes.InvalidArgument)
  // FailedPrecondition error
  FailedPrecondition ErrorType = ErrorType(codes.FailedPrecondition)
  // OutOfRange error
  OutOfRange ErrorType = ErrorType(codes.OutOfRange)
  ...
  //ADD here
)
```

and to httpMap for converting to http code

```GO
var httpMap = map[ErrorType]uint32{
  InvalidArgument:    400,
  FailedPrecondition: 400,
  OutOfRange:         400,
  ...
  //ADD here
}
```
