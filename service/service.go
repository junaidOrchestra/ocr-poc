package service

type ServiceContainer struct {
	SupplierService SupplierService
	// Add more services as needed
}

func NewServiceContainer(supplierService SupplierService) *ServiceContainer {
	return &ServiceContainer{
		SupplierService: supplierService,
		// Initialize more services as needed
	}
}
