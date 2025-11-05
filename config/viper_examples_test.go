package config

import (
	"fmt"
	"reflect"
)

func ExampleEnv_basic() {
	type User struct {
		Name string `env:"NAME"`
		Age  int    `env:"AGE"`
	}

	env, err := Env(&User{
		Name: "John",
		Age:  30,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("NAME=%s\n", env["NAME"])
	fmt.Printf("AGE=%s\n", env["AGE"])
	// Output:
	// NAME=John
	// AGE=30
}

func ExampleEnv_nested() {
	type User struct {
		Name string `env:"NAME"`
		Age  int    `env:"AGE"`
	}

	type Organization struct {
		Admin *User `env:"ADMIN"`
	}

	env, err := Env(&Organization{
		Admin: &User{
			Name: "John",
			Age:  30,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("ADMIN_NAME=%s\n", env["ADMIN_NAME"])
	fmt.Printf("ADMIN_AGE=%s\n", env["ADMIN_AGE"])
	// Output:
	// ADMIN_NAME=John
	// ADMIN_AGE=30
}

func ExampleEnv_nestedWithDefaultEnvNaming() {
	type User struct {
		Name string
		Age  int
	}

	type Organization struct {
		AdminA *User
		AdminB *User
	}

	env, err := Env(&Organization{
		AdminA: &User{
			Name: "John",
			Age:  30,
		},
		AdminB: &User{
			Name: "Jane",
			Age:  31,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("ADMINA_NAME=%s\n", env["ADMINA_NAME"])
	fmt.Printf("ADMINA_AGE=%s\n", env["ADMINA_AGE"])
	fmt.Printf("ADMINB_NAME=%s\n", env["ADMINB_NAME"])
	fmt.Printf("ADMINB_AGE=%s\n", env["ADMINB_AGE"])
	// Output:
	// ADMINA_NAME=John
	// ADMINA_AGE=30
	// ADMINB_NAME=Jane
	// ADMINB_AGE=31
}

func ExampleEnv_nestedWithSkipTag() {
	type User struct {
		Name string
		Age  int
	}

	type Organization struct {
		AdminA *User `env:"-"`
		AdminB *User
	}

	env, err := Env(&Organization{
		AdminA: &User{
			Name: "John",
			Age:  30,
		},
		AdminB: &User{
			Name: "Jane",
			Age:  31,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("NAME=%s\n", env["NAME"])
	fmt.Printf("AGE=%s\n", env["AGE"])
	fmt.Printf("ADMINB_NAME=%s\n", env["ADMINB_NAME"])
	fmt.Printf("ADMINB_AGE=%s\n", env["ADMINB_AGE"])
	// Output:
	// NAME=John
	// AGE=30
	// ADMINB_NAME=Jane
	// ADMINB_AGE=31
}

type Role int

func (r Role) String() string {
	if r == Role(0) {
		return "admin"
	}
	return "user"
}

func ExampleEnv_defaultEncodeHook() {
	type User struct {
		Name string
		Age  int
		Role Role
	}

	// Role implements fmt.Stringer, so it will be encoded as a string
	env, err := Env(&User{
		Name: "John",
		Age:  30,
		Role: Role(1),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("NAME=%s\n", env["NAME"])
	fmt.Printf("AGE=%s\n", env["AGE"])
	fmt.Printf("ROLE=%s\n", env["ROLE"])
	// Output:
	// NAME=John
	// AGE=30
	// ROLE=user
}

func ExampleEnv_customEncodeHook() {
	type MyCustomType int

	customEncodeHook := func(val reflect.Value) (reflect.Value, error) {
		if val.Type() == reflect.TypeOf(MyCustomType(0)) {
			return reflect.ValueOf(fmt.Sprintf("my great type: %d", val.Interface())), nil
		}
		return val, nil
	}

	type MyType struct {
		MyCustomField MyCustomType `env:"MY_CUSTOM_FIELD"`
	}

	env, err := Env(&MyType{
		MyCustomField: MyCustomType(1),
	}, customEncodeHook)
	if err != nil {
		panic(err)
	}

	fmt.Printf("MY_CUSTOM_FIELD=%s\n", env["MY_CUSTOM_FIELD"])
	// Output:
	// MY_CUSTOM_FIELD=my great type: 1
}
