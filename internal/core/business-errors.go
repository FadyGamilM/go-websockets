package core

type BUSINESS_ERROR string

const (
	Error_Internal_Logic        BUSINESS_ERROR = "internal logic error"
	Error_Non_Existing_Resource BUSINESS_ERROR = "resource doesn't exist"
	Error_WRONG_AUTH_CREDENTIAL BUSINESS_ERROR = "wrong crednetials"
)

type BusinessError struct {
	DbErr string
	Type  BUSINESS_ERROR
	Msg   string
}

func (be *BusinessError) Error() string {
	return be.Msg
}

func New_Business_InternalLogicError(dbErr string) *BusinessError {
	return &BusinessError{
		DbErr: dbErr,
		Type:  Error_Internal_Logic,
		Msg:   string(Error_Internal_Logic),
	}
}

func New_Business_NonExistingResourceError(dbErr string) *BusinessError {
	return &BusinessError{
		DbErr: dbErr,
		Type:  Error_Non_Existing_Resource,
		Msg:   string(Error_Non_Existing_Resource),
	}
}

func New_Business_WrongAuthCredentialsError(dbErr string) *BusinessError {
	return &BusinessError{
		DbErr: dbErr,
		Type:  Error_WRONG_AUTH_CREDENTIAL,
		Msg:   string(Error_WRONG_AUTH_CREDENTIAL),
	}
}
