// Package pagy is a package for collecting pagination data from clients and returning a consistent paginated structure.
/*
pagy is not limited to any one framework or library, it works directly with the stdlib http.Request interface
so any project/framework that supports stdlib it will support.

pagy works by collecting pagination based query params from the http.Request and formatting it into an expected
and consistent structure to use throughout your service. A set of tools is available for getting the pagination data
as well as structures that are generic-based for responding with a consistent data contract.
*/
package pagy
