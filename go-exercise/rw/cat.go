package main
import(
	"fmt"
	"reflect"
)

//接口
type Cat interface {
    Meow()
}
//实现类1
type Tabby struct{}
func (*Tabby) Meow() { fmt.Println("Tabby meow") }
func GetNilTabbyCat() Cat {
    var myTabby *Tabby = nil
    return myTabby
}
func GetTabbyCat() Cat {
    var myTabby *Tabby = &Tabby{}
    return myTabby
}

//实现类2
type Gafield struct{}
func (*Gafield) Meow() { fmt.Println("Gafield meow") }
func GetNilGafieldCat() Cat {
    var myGafield *Gafield = nil
    return myGafield
}
func GetGafieldCat() Cat {
    var myGafield *Gafield = &Gafield{}
    return myGafield
}

func main() {
	var (
		cat1 Cat = nil
		cat2     = GetNilTabbyCat()
		cat3     = GetTabbyCat()
		cat4     = GetNilGafieldCat()
	)
	fmt.Printf("cat1 information: nil?:%5v, type=%15v, value=%5v  \n", cat1 == nil, reflect.TypeOf(cat1), reflect.ValueOf(cat1))                                                //接口变量，type、value都是nil，所以cat1==nil
	fmt.Printf("cat2 information: nil?:%5v, type=%15v, type.kind=%5v, value=%5v  \n", cat2 == nil, reflect.TypeOf(cat2), reflect.TypeOf(cat2).Kind(), reflect.ValueOf(cat2))   //接口变量，type!=nil，所以cat2!==nil
	fmt.Printf("cat3 information: nil?:%5v, type=%15v, type.kind=%5v, value=%5v  \n", cat3 == nil, reflect.TypeOf(cat3), reflect.TypeOf(cat3).Kind(), reflect.ValueOf(cat3))   //接口变量，type!=nil, 所以cat3!=nil
	fmt.Printf("cat4 information: nil?:%5v, type=%15v, type.kind=%5v, value=%5v  \n", cat4 == nil, reflect.TypeOf(cat4), reflect.TypeOf(cat4).Kind(), reflect.ValueOf(cat4)) //接口变量，
	fmt.Printf("cat1==cat2?%5v , cat2==cat3?%5v， cat2==cat4?%5v \n", cat1 == cat2, cat2 == cat3, cat2 == cat4)

	//fmt.Printf("cat1 information: type=%v,kind=%v \n",reflect.TypeOf(cat2),reflect.TypeOf(cat2).Kind())
}
