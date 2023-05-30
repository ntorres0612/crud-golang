package controllers

import "github.com/ntorres0612/ionix-crud/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
	s.Router.HandleFunc("/signup", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")

	//Store routes
	s.Router.HandleFunc("/store", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateStore))).Methods("POST")
	s.Router.HandleFunc("/stores", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetStores))).Methods("GET")
	s.Router.HandleFunc("/store/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetStore))).Methods("GET")
	s.Router.HandleFunc("/store/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateStore))).Methods("PUT")
	s.Router.HandleFunc("/store/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteStore)).Methods("DELETE")

	//Logistic type routes
	s.Router.HandleFunc("/logistic-types", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetLogisticTypes))).Methods("GET")
	s.Router.HandleFunc("/logistic-type/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetLogisticType))).Methods("GET")

	//Truck routes
	s.Router.HandleFunc("/truck", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateTruck))).Methods("POST")
	s.Router.HandleFunc("/trucks", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetTrucks))).Methods("GET")
	s.Router.HandleFunc("/truck/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetTruck))).Methods("GET")
	s.Router.HandleFunc("/truck/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateTruck))).Methods("PUT")
	s.Router.HandleFunc("/truck/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteTruck)).Methods("DELETE")

	//Product types
	s.Router.HandleFunc("/product_type", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateProductType))).Methods("POST")
	s.Router.HandleFunc("/product_types", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetProductTypes))).Methods("GET")
	s.Router.HandleFunc("/product_type/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetProductType))).Methods("GET")
	s.Router.HandleFunc("/product_type/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateProductType))).Methods("PUT")
	s.Router.HandleFunc("/product_type/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteProductType)).Methods("DELETE")

	//Deliveries
	s.Router.HandleFunc("/delivery", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateDelivery))).Methods("POST")
	s.Router.HandleFunc("/deliveries", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetDeliverys))).Methods("GET")
	s.Router.HandleFunc("/delivery/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetDelivery))).Methods("GET")
	s.Router.HandleFunc("/delivery/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateDelivery))).Methods("PUT")
	s.Router.HandleFunc("/delivery/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteDelivery)).Methods("DELETE")

	//Customers
	s.Router.HandleFunc("/customer", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateCustomer))).Methods("POST")
	s.Router.HandleFunc("/customers", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetCustomers))).Methods("GET")
	s.Router.HandleFunc("/customer/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetCustomer))).Methods("GET")
	s.Router.HandleFunc("/customer/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateCustomer))).Methods("PUT")
	s.Router.HandleFunc("/customer/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteCustomer)).Methods("DELETE")

}
