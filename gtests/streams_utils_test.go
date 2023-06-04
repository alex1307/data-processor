package gtests

import (
	utils "data-processor/utils"
	"fmt"
	"testing"

	"strings"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name string
	Age  int
}

type Animal struct {
	Name  string
	Sound string
}

func uppercaseString(s string) string {
	return strings.ToUpper(s)
}

func lowercaseString(s string) string {
	return string([]rune(s)[0]-'A'+'a') + s[1:]
}

func TestStreamsUtils(t *testing.T) {
	people := []Person{
		{Name: "John", Age: 30},
		{Name: "Jane", Age: 25},
	}

	// Map people to uppercase names
	upperNameMap := func(item Person) Person {
		item.Name = uppercaseString(item.Name)
		return item
	}
	mappedPeople := utils.Map(people, upperNameMap)
	assert.Equal(t, "JOHN", mappedPeople[0].Name)
	// Print the mapped results
	for _, person := range mappedPeople {
		fmt.Println(person.Name)

	}

	animals := []Animal{
		{Name: "Cat", Sound: "Meow"},
		{Name: "Dog", Sound: "Woof"},
	}

	// Map animals to lowercase names
	lowerNameMap := func(item Animal) Animal {
		item.Name = lowercaseString(item.Name)
		return item
	}
	mappedAnimals := utils.Map(animals, lowerNameMap)

	// Print the mapped results
	for _, animal := range mappedAnimals {

		fmt.Println(animal.Name)

	}

}
