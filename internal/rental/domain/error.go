package domain

import "errors"

var (
	ErrInvalidEntity = errors.New("invalid entity")

	ErrInvalidMaintenance = errors.New("invalid maintenance")
	ErrInvalidTransit     = errors.New("invalid transit")
	ErrInvalidReserve     = errors.New("invalid reserve")
	ErrInvalidPark        = errors.New("invalid park")

	ErrInvalidReservedDate = errors.New("date reserved from is bigger than date reserved to")
	ErrCarUnavailable      = errors.New("car is not available to rent")
	ErrCategoryUnavailable = errors.New("car is not available in this category")
	ErrInvalidCarStation   = errors.New("car is not in this station")
	ErrIvalidConfirmDate   = errors.New("date is not betweeen reserved dates")
	ErrClose               = errors.New("rent order can not be closed")
	ErrIvalidCloseDate     = errors.New("date to is not bigger than date from")
	ErrIvalidCloseTax      = errors.New("close tax is invalid")
	ErrIvalidCloseDiscount = errors.New("close discount is invalid")
	ErrCancel              = errors.New("rent order can not be canceled")
)
