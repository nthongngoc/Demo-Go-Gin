package shared

import (
	"log"
	"reflect"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gopkg.in/mgo.v2/bson"
)

type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &DefaultValidator{}

func (v *DefaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

func (v *DefaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
		if err := v.validate.RegisterValidation("mongoid", func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			if value == "" {
				return true
			}
			ok := bson.IsObjectIdHex(value)
			if ok {
				objectID := bson.ObjectIdHex(value)
				ref := reflect.ValueOf(objectID)
				fl.Field().SetString(ref.String())
			}
			return ok
		}); err != nil {
			log.Fatal(err)
		}
	})
}
