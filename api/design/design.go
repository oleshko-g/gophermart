// Package design is the gophermart goa API design
package design

import . "goa.design/goa/v3/dsl" // revive:disable-line:dot-imports goa recommends to dot import the DSL

var _ = API("gophermart", func() {
	// INFO: Gophermart
	Version("1.0")
	License(func() {
		Name("Apache 2.0")
		URL("https://github.com/oleshko-g/gophermart/blob/dev/LICENSE")
	})
	HTTP(func() {
		Path("/api")
		Consumes("text/plain", "application/json")
	})
})

var _ = Service("docs", func() {
	// INFO:  Docs
	Files("/openapi.yaml", "./internal/gen/http/openapi.yaml", func() {
		Description("OpenAPI 2.0")
	})
	Files("/openapi3.yaml", "./internal/gen/http/openapi3.yaml", func() {
		Description("OpenAPI 3.0")
	})
})

var _ = Service("user", func() {
	// INFO:  User service
	Error("Invalid input parameter", errorType)
	Error("User is not authenticated", errorType)
	Error("Internal service error", errorType)

	// INFO:    Register
	Method("register", func() {
		Payload(loginPassword)
		Result(_JWTToken)
		Error("Login is taken already", errorType)
		HTTP(func() {
			POST("/user/register")
			Response(StatusOK, func() {
				Header("authToken:Authorization")
				Body(Empty)
			})
			Response("Invalid input parameter", StatusBadRequest, func() {
				Body(Empty)
			})
			Response("Login is taken already", StatusConflict, func() {
				Description("Login is taken already")
				Body(Empty)
			})
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
		})

	})
	// INFO:    Login
	Method("login", func() {
		Payload(loginPassword)
		Result(_JWTToken)
		HTTP(func() {
			POST("/user/login")
			Response(StatusOK, func() {
				Header("authToken:Authorization")
				Body(Empty)
			})
			Response("Invalid input parameter", StatusBadRequest, func() {
				Body(Empty)
			})
			Response("User is not authenticated", StatusUnauthorized, func() {
				Description("User is not authenticated")
				Body(Empty)
			})
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
		})
	})
})
var _ = Service("balance", func() {
	// INFO:  Balance Service
	Security(_JWTAuth)
	Error("Invalid input parameter", errorType)
	Error("User is not authenticated", errorType)
	Error("Internal service error", errorType)
	Error("Not implemented", errorType)
	Error("missing_field")
	HTTP(func() {
		Header("Authorization", func() {
			Example(func() {
				Value("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")
			})
		})

	})

	// INFO:    UploadUserOrder
	Method("UploadUserOrder", func() {
		Description("Upload user order")
		Result(func() {
			Attribute("accepted", func() {
				Meta("struct:tag:json", "-")
				Meta("openapi:generate", "false")
				Meta("openapi:example", "false")
			})
			Meta("openapi:example", "false")
		})
		Payload(func() {
			Token("Authorization", String, "A JWT token used to authenticate a request", func() {
				Example(func() {
					Value("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")
				})
			})
			Attribute("OrderNumber", String, func() {
				Description("Unique user order number")
				Pattern("[1-9][0-9]*")
			})
			Required("Authorization", "OrderNumber")
		})
		Error("The order belongs to another user", errorType)
		Error("Invalid order number", errorType)
		HTTP(func() {
			POST("/user/orders")
			Body("OrderNumber", func() {
				Example(func() {
					Value("12345678903")
				})
			})
			Response(StatusOK, func() {
				Description("The order has been accepted for processing before.")
				Body(Empty)
			})
			Response(StatusAccepted, func() {
				Tag("accepted", "yes")
				Description("The order has been accepted for processing.")
				Body(Empty)
			})
			Response("Invalid input parameter", StatusBadRequest, func() {
				Body(Empty)
			})
			Response("The order belongs to another user", StatusConflict, func() {
				Description("The order belongs to another user")
				Body(Empty)
			})
			Response("Invalid order number", StatusUnprocessableEntity, func() {
				Description("Invalid order number")
				Body(Empty)
			})
			Response("User is not authenticated", StatusUnauthorized, func() {
				Description("User is not authenticated")
				Body(Empty)
			})
			Response("missing_field", StatusUnauthorized, func() {
				Description("Missing or empty Authorization header")
				Body(Empty)
			})
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
		})
	})
	// INFO:    ListUserOrders
	Method("ListUserOrders", func() {
		Description("List user orders")
		Payload(func() {
			Token("Authorization", String, "A JWT token used to authenticate a request", func() {
				Example(func() {
					Value("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")
				})
			})
			Required("Authorization")
		})
		Result(func() {
			Attribute("orders", ArrayOf(order), func() {
				Example(func() {
					Value([]Val{
						{
							"number":      "9278923470",
							"status":      "PROCESSED",
							"accrual":     500,
							"uploaded_at": "2020-12-10T15:15:45+03:00",
						},
						{
							"number":      "12345678903",
							"status":      "PROCESSING",
							"uploaded_at": "2020-12-10T15:12:01+03:00",
						},
						{
							"number":      "346436439",
							"status":      "INVALID",
							"uploaded_at": "2020-12-09T16:09:53+03:00",
						}})
				})
			})
			Attribute("no orders", func() {
				Meta("struct:tag:json", "-")
				Meta("openapi:generate", "false")
				Meta("openapi:example", "false")
			})
		})
		HTTP(func() {
			GET("/user/orders")
			Response(StatusOK, func() {
				Body("orders")
			})
			Response(StatusNoContent, func() {
				Tag("no orders", "yes")
				Description("No orders available")
				Body(Empty)
			})
			Response("User is not authenticated", StatusUnauthorized, func() {
				Description("User is not authenticated")
				Body(Empty)
			})
			Response("missing_field", StatusUnauthorized, func() {
				Description("Missing or empty Authorization header")
				Body(Empty)
			})
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
		})
	})
	// INFO:    GetUserBalance
	Method("GetUserBalance", func() {
		Description("Get user balance")
		Payload(func() {
			Token("Authorization", String, "A JWT token used to authenticate a request", func() {
				Example(func() {
					Value("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")
				})
			})
			Required("Authorization")
		})
		Result(func() {
			Attribute("current", Float64, func() {
				Minimum(0)
			})
			Attribute("withdrawn", Float64, func() {
				Minimum(0)
			})
			Example(func() {
				Value(Val{
					"current":   500.5,
					"withdrawn": 42,
				})
			})
			Required("current", "withdrawn")
		})
		HTTP(func() {
			GET("/user/balance")
			Response(StatusOK)
			Response("User is not authenticated", StatusUnauthorized, func() {
				Description("User is not authenticated")
				Body(Empty)
			})
			Response("missing_field", StatusUnauthorized, func() {
				Description("Missing or empty Authorization header")
				Body(Empty)
			})
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
		})
	})
	// INFO:    WithdrawUserBalance
	Method("WithdrawUserBalance", func() {
		Payload(func() {
			Token("Authorization", String, "A JWT token used to authenticate a request", func() {
				Example(func() {
					Value("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")
				})
			})
			Attribute("order", String, func() {
				Pattern("[1-9][0-9]*")
			})
			Attribute("sum", Float64, func() {
				ExclusiveMinimum(0)
			})
			Required("Authorization", "order", "sum")
			Example(func() {
				Value(
					Val{
						"order": "2377225624",
						"sum":   751,
					})
			})
		})
		Error("Insufficient funds", errorType)
		Error("Invalid order number", errorType)
		HTTP(func() {
			POST("/user/balance/withdraw")
			Response(StatusOK, func() {
				Body(Empty)
			})
			Response("User is not authenticated", StatusUnauthorized, func() {
				Description("User is not authenticated")
			})
			Response("missing_field", StatusUnauthorized, func() {
				Body(Empty)
				Description("Missing or empty Authorization header")
			})
			Response("Insufficient funds", StatusPaymentRequired, func() {
				Body(Empty)
				Description("Insufficient funds")
			})
			Response("Invalid order number", StatusUnprocessableEntity, func() {
				Body(Empty)
				Description("Invalid order number")
			})
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
		})
	})
	// INFO:    GetWithdrawals
	Method("GetWithdrawals", func() {
		Payload(func() {
			Token("Authorization", String, "A JWT token used to authenticate a request", func() {
				Example(func() {
					Value("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")
				})
			})
			Required("Authorization")
		})
		Result(func() {
			Attribute("withdrawals", ArrayOf(withdrawal))
			Attribute("NoResult")
		})
		HTTP(func() {
			GET("/user/withdrawals")
			Response(StatusOK, func() {
				Body("withdrawals")
			})
			Response(StatusNoContent, func() {
				Tag("NoResult", "true")
				Description("User has no withdrawals")
				Body(Empty)
			})
			Response("Invalid input parameter", StatusBadRequest, func() {
				Body(Empty)
			})
			Response("User is not authenticated", StatusUnauthorized, func() {
				Description("User is not authenticated")
				Body(Empty)
			})
			Response("missing_field", StatusUnauthorized, func() {
				Description("Missing or empty Authorization header")
				Body(Empty)
			})
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
			Response("Not implemented", StatusNotImplemented, func() {
				Body(Empty)
			})
		})
	})
})

var uploadUserOrderResult = Type("PostOrderResult", func() {
	Attribute("accepted", func() {
		Meta("struct:tag:json", "-")
		Meta("openapi:generate", "false")
		Meta("openapi:example", "false")
	})
	Meta("openapi:example", "false")
})

var loginPassword = Type("LoginPassword", func() {
	Attribute("login", String)
	Attribute("password", String)
	Required("login", "password")
	Example(func() {
		Value(Val{
			"login":    "<login>",
			"password": "<password>",
		})
	})
})

var errorType = Type("GophermartError", func() {
	ErrorName("name", func() {
		Description("identifier to map an error to HTTP status codes")
		Meta("struct:tag:json", "-")
		Meta("openapi:generate", "false")
		Meta("openapi:example", "false")
	})
	Required("name")
	Meta("openapi:generate", "false")
	Meta("openapi:example", "false")
	Meta("struct:pkg:path", "service")
})

var _JWTAuth = JWTSecurity("jwt", func() {
	Description("Secures an endpoint by requiring a valid JWT token.")
})

var _JWTToken = Type("JWTToken", func() {
	Token("authToken", String, func() {
		Description("A JWT token used to authenticate a request")
		Example(func() {
			Value("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")
		})
	})
	Required("authToken")
	Meta("struct:pkg:path", "service")
})

var order = Type("Order", func() {
	Attribute("number", String, func() {
		Pattern("[1-9][0-9]*")
	})
	Attribute("status", String, func() {
		Enum("NEW", "PROCESSING", "INVALID", "PROCESSED")
	})
	Attribute("accrual", Float64, func() {
		ExclusiveMinimum(0)
	})
	Attribute("uploaded_at", String, func() {
		Format(FormatDateTime)
	})
	Required("number", "status", "uploaded_at")
})

var orderNumber = Type("OrderNumber", String, func() {
	Description("Unique user order number")
	Pattern("[1-9][0-9]*")
})

var withdrawal = Type("Withdrawal", func() {
	Attribute("order", orderNumber)
	Attribute("sum", Float64, func() {
		ExclusiveMinimum(0)
	})
	Attribute("processed_at", String, func() {
		Format(FormatDateTime)
	})
})
