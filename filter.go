package filter

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var (
	// ErrInvalidCondition error
	ErrInvalidCondition = errors.New("invalid condition")
	// ErrInvalidStructType error
	ErrInvalidStructType = errors.New("object is not a struct object")
	// ErrInvalidRangeFmt error
	ErrInvalidRangeFmt = errors.New("invalid range format tag")
	// ErrInvalidNumberFmt error
	ErrInvalidNumberFmt = errors.New("invalid number format tag")

	// default namespace
	DefaultNamespace = "fitler"
)

type (
	// Option struct
	Option struct {
		Namespace string
		Condition interface{}
	}
	// config struct
	config struct {
		namespace string
		level     int
		match     string
	}
)

// Go to filter struct
func Go(target interface{}, opts ...Option) (interface{}, error) {
	c := new(config)
	for _, opt := range opts {
		c.namespace = opt.Namespace
		switch o := opt.Condition.(type) {
		case string:
			c.match = o
		case int:
			c.level = o
		case int8:
			c.level = int(o)
		case int16:
			c.level = int(o)
		case int32:
			c.level = int(o)
		case int64:
			c.level = int(o)
		case uint:
			c.level = int(o)
		case uint8:
			c.level = int(o)
		case uint16:
			c.level = int(o)
		case uint32:
			c.level = int(o)
		case uint64:
			c.level = int(o)
		default:
			return nil, ErrInvalidCondition
		}
		break
	}
	if s := strings.TrimSpace(c.namespace); s == "" {
		c.namespace = DefaultNamespace
	}
	res, err := run(target, c)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// run copy value to new object depends on match mode
func run(o interface{}, conf *config) (interface{}, error) {
	v := reflect.ValueOf(o)
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, ErrInvalidStructType
	}
	c := reflect.New(reflect.TypeOf(o).Elem())
	if c.Kind() == reflect.Interface || c.Kind() == reflect.Ptr {
		c = c.Elem()
	}
	if c.Kind() != reflect.Struct {
		return nil, ErrInvalidStructType
	}
	for i := 0; i < v.NumField(); i++ {
		fval := v.Field(i)
		ftyp := v.Type().Field(i)
		ftag := strings.TrimSpace(ftyp.Tag.Get(conf.namespace))
		tval := c.Field(i)
		ttyp := c.Type().Field(i)
		ttag := strings.TrimSpace(ttyp.Tag.Get(conf.namespace))
		if !tval.IsValid() {
			continue
		}
		if !tval.CanSet() {
			continue
		}
		if ftag == "" || ftag == "-" || ttag == "" || ttag == "-" {
			continue
		}
		var o interface{}
		if fval.Kind() == reflect.Slice {
			d := reflect.MakeSlice(fval.Type(), fval.Len(), fval.Cap())
			for i := 0; i < fval.Len(); i++ {
				v, err := run(fval.Index(i).Interface(), conf)
				if err != nil {
					return nil, err
				}
				d.Index(i).Set(reflect.ValueOf(v))
			}
			o = d.Interface()
		} else if t := fval; t.Kind() == reflect.Interface || t.Kind() == reflect.Ptr {
			v, err := run(t.Interface(), conf)
			if err != nil {
				return nil, err
			}
			o = v
		} else {
			o = fval.Interface()
		}
		if o == nil {
			continue
		}
		conditions := strings.Split(ftag, ",")
		for _, condition := range conditions {
			switch {
			case strings.Contains(condition, "-"):
				ranges := strings.Split(condition, "-")
				if len(ranges) != 2 {
					return nil, ErrInvalidRangeFmt
				}
				var from, to int
				switch n, err := strconv.Atoi(ranges[0]); {
				case err != nil:
					return nil, ErrInvalidNumberFmt
				default:
					from = n
				}
				switch n, err := strconv.Atoi(ranges[1]); {
				case err != nil:
					return nil, ErrInvalidNumberFmt
				default:
					to = n
				}
				if conf.level >= from && conf.level <= to {
					tval.Set(reflect.ValueOf(o))
					continue
				}
			case strings.HasPrefix(condition, ">="):
				condition := strings.TrimPrefix(condition, ">=")
				switch n, err := strconv.Atoi(condition); {
				case err != nil:
					return nil, ErrInvalidNumberFmt
				default:
					if conf.level >= n {
						tval.Set(reflect.ValueOf(o))
						continue
					}
				}
			case strings.HasPrefix(condition, ">"):
				condition := strings.TrimPrefix(condition, ">")
				switch n, err := strconv.Atoi(condition); {
				case err != nil:
					return nil, ErrInvalidNumberFmt
				default:
					if conf.level > n {
						tval.Set(reflect.ValueOf(o))
						continue
					}
				}
			case strings.HasPrefix(condition, "<="):
				condition := strings.TrimPrefix(condition, "<=")
				switch n, err := strconv.Atoi(condition); {
				case err != nil:
					return nil, ErrInvalidNumberFmt
				default:
					if conf.level <= n {
						tval.Set(reflect.ValueOf(o))
						continue
					}
				}
			case strings.HasPrefix(condition, "<"):
				condition := strings.TrimPrefix(condition, "<")
				switch n, err := strconv.Atoi(condition); {
				case err != nil:
					return nil, ErrInvalidNumberFmt
				default:
					if conf.level < n {
						tval.Set(reflect.ValueOf(o))
						continue
					}
				}
			default:
				if n, err := strconv.Atoi(condition); err == nil && conf.level == n {
					tval.Set(reflect.ValueOf(o))
					continue
				} else if conf.match == condition {
					tval.Set(reflect.ValueOf(o))
					continue
				}
			}
		}
	}
	return c.Addr().Interface(), nil
}
