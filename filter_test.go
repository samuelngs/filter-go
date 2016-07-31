package filter

import "testing"

type num struct {
	Field1 string `custom-filter:">=1,10,30-50"`
	Field2 string `custom-filter:"<=1,100"`
	Field3 string `custom-filter:">0,10-20,50-70,100"`
}

var custom = &num{
	Field1: "value-1",
	Field2: "value-2",
	Field3: "value-3",
}

type children struct {
	Field1 string `custom-filter:"read"`
	Field2 string `custom-filter:"write"`
	Field3 string `custom-filter:"read,write"`
}

type object struct {
	Field1 string      `custom-filter:"read,write"`
	Field2 string      `custom-filter:"read"`
	Field3 string      `custom-filter:"read,write"`
	Field4 string      `custom-filter:"read,write"`
	Field5 string      `custom-filter:""`
	Field6 string      `custom-filter:"write"`
	Field7 *children   `custom-filter:"read,write"`
	Field8 []*children `custom-filter:"read,write"`
}

var test = &object{
	Field1: "value-1",
	Field2: "value-2",
	Field3: "value-3",
	Field4: "value-4",
	Field5: "value-5",
	Field6: "value-6",
	Field7: &children{
		Field1: "value-1",
		Field2: "value-2",
		Field3: "value-3",
	},
	Field8: []*children{
		&children{
			Field1: "value-a-1",
			Field2: "value-a-2",
			Field3: "value-a-3",
		},
		&children{
			Field1: "value-b-1",
			Field2: "value-b-2",
			Field3: "value-b-3",
		},
	},
}

func TestReadable(t *testing.T) {
	r, err := Go(test, Option{
		Namespace: "custom-filter",
		Condition: "read",
	})
	if err != nil {
		t.Fatal(err)
	}
	o, ok := r.(*object)
	if !ok {
		t.Fatal("invalid object type")
	}
	if o.Field1 != test.Field1 {
		t.Errorf("Expected %v, got %v", test.Field1, o.Field1)
	}
	if o.Field2 != test.Field2 {
		t.Errorf("Expected %v, got %v", test.Field2, o.Field2)
	}
	if o.Field3 != test.Field3 {
		t.Errorf("Expected %v, got %v", test.Field3, o.Field3)
	}
	if o.Field4 != test.Field4 {
		t.Errorf("Expected %v, got %v", test.Field4, o.Field4)
	}
	if o.Field5 != "" {
		t.Errorf("Expected empty value, got %v", o.Field5)
	}
	if o.Field6 != "" {
		t.Errorf("Expected empty value, got %v", o.Field6)
	}
	if o.Field7.Field1 != test.Field7.Field1 {
		t.Errorf("Expected %v, got %v", test.Field7.Field1, o.Field7.Field1)
	}
	if o.Field7.Field2 != "" {
		t.Errorf("Expected empty value, got %v", o.Field7.Field2)
	}
	if o.Field7.Field3 != test.Field7.Field3 {
		t.Errorf("Expected %v, got %v", test.Field7.Field3, o.Field7.Field3)
	}
}

func TestWritable(t *testing.T) {
	r, err := Go(test, Option{
		Namespace: "custom-filter",
		Condition: "write",
	})
	if err != nil {
		t.Fatal(err)
	}
	o, ok := r.(*object)
	if !ok {
		t.Fatal(err)
	}
	if o.Field1 != test.Field1 {
		t.Errorf("Expected %v, got %v", test.Field1, o.Field1)
	}
	if o.Field2 != "" {
		t.Errorf("Expected empty value, got %v", o.Field2)
	}
	if o.Field3 != test.Field3 {
		t.Errorf("Expected %v, got %v", test.Field3, o.Field3)
	}
	if o.Field4 != test.Field4 {
		t.Errorf("Expected %v, got %v", test.Field4, o.Field4)
	}
	if o.Field5 != "" {
		t.Errorf("Expected empty value, got %v", o.Field5)
	}
	if o.Field6 != test.Field6 {
		t.Errorf("Expected %v, got %v", test.Field6, o.Field6)
	}
}

func TestFilter(t *testing.T) {
	if r, err := Go(custom, Option{
		Namespace: "custom-filter",
		Condition: 1,
	}); err != nil {
		t.Fatal(err)
	} else {
		if o, ok := r.(*num); !ok {
			t.Fatal("invalid object type")
		} else {
			if o.Field1 != custom.Field1 {
				t.Errorf("Expected %v, got %v", custom.Field1, o.Field1)
			}
			if o.Field2 != custom.Field2 {
				t.Errorf("Expected %v, got %v", custom.Field2, o.Field2)
			}
			if o.Field3 != custom.Field3 {
				t.Errorf("Expected %v, got %v", custom.Field3, o.Field3)
			}
		}
	}
	if r, err := Go(custom, Option{
		Namespace: "custom-filter",
		Condition: 50,
	}); err != nil {
		t.Fatal(err)
	} else {
		if o, ok := r.(*num); !ok {
			t.Fatal("invalid object type")
		} else {
			if o.Field1 != custom.Field1 {
				t.Errorf("Expected %v, got %v", custom.Field1, o.Field1)
			}
			if o.Field2 != "" {
				t.Errorf("Expected empty value, got %v", o.Field2)
			}
			if o.Field3 != custom.Field3 {
				t.Errorf("Expected %v, got %v", custom.Field3, o.Field3)
			}
		}
	}
}
