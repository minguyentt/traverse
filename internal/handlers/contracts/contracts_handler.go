package contracts

type APIContractsHandler interface {
	GetContracts() error
    // GetContractsByID() error
    // PostNewContract() error
    // DeleteContractByID() error
}

type ContractsHandler struct {}

/*
- The handlers will be the exposed public api layer
- performs the URL parsing from the request 
- calls the appropriate service => repository => db
- services will return the DTO json model
*/

func (c *ContractsHandler) GetContracts() error {

    return nil
}




