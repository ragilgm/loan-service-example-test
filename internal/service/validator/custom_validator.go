package validator

type CustomValidator interface {
	ValidateCreate(data interface{}) error
	ValidateUpdate(data interface{}) error
	ValidateTransitionStatus(from interface{}, to interface{}) bool
}
