package survey

import (
	"errors"

	"github.com/djgilcrease/survey/core"
)

// Validator is a function passed to a Question after a user has provided a response.
// If the function returns an error, then the user will be prompted again for another
// response.
type Validator func(ans interface{}) error

// Transformer is a function passed to a Question after a user has provided a response.
// The function can be used to implement a custom logic that will result to return
// a different representation of the given answer.
//
// Look `TransformString`, `ToLower` `Title` and `ComposeTransformers` for more.
type Transformer func(ans interface{}) (newAns interface{})

// Question is the core data structure for a survey questionnaire.
type Question struct {
	Name      string
	Prompt    Prompt
	Validate  Validator
	Transform Transformer
}

/*
AskOne performs the prompt for a single prompt and asks for validation if required.
Response types should be something that can be casted from the response type designated
in the documentation. For example:

	name := ""
	prompt := &survey.Input{
		Message: "name",
	}

	survey.AskOne(prompt, &name, nil)

*/
func AskOne(p Prompt, response interface{}, v Validator) error {
	err := Ask([]*Question{{Prompt: p, Validate: v}}, response)
	if err != nil {
		return err
	}

	return nil
}

/*
Ask performs the prompt loop, asking for validation when appropriate. The response
type can be one of two options. If a struct is passed, the answer will be written to
the field whose name matches the Name field on the corresponding question. Field types
should be something that can be casted from the response type designated in the
documentation. Note, a survey tag can also be used to identify a Otherwise, a
map[string]interface{} can be passed, responses will be written to the key with the
matching name. For example:

	qs := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "What is your name?"},
			Validate: survey.Required,
			Transform: survey.Title,
		},
	}

	answers := struct{ Name string }{}


	err := survey.Ask(qs, &answers)
*/
func Ask(qs []*Question, response interface{}) error {

	// if we weren't passed a place to record the answers
	if response == nil {
		// we can't go any further
		return errors.New("cannot call Ask() with a nil reference to record the answers")
	}

	// go over every question
	for _, q := range qs {
		// grab the user input and save it
		ans, err := q.Prompt.Prompt()
		// if there was a problem
		if err != nil {
			return err
		}

		// if there is a validate handler for this question
		if q.Validate != nil {
			// wait for a valid response
			for invalid := q.Validate(ans); invalid != nil; invalid = q.Validate(ans) {
				err := q.Prompt.Error(invalid)
				// if there was a problem
				if err != nil {
					return err
				}

				// ask for more input
				ans, err = q.Prompt.Prompt()
				// if there was a problem
				if err != nil {
					return err
				}
			}
		}

		if q.Transform != nil {
			// check if we have a transformer available, if so
			// then try to acquire the new representation of the
			// answer, if the resulting answer is not nil.
			if newAns := q.Transform(ans); newAns != nil {
				ans = newAns
			}
		}

		// tell the prompt to cleanup with the validated value
		q.Prompt.Cleanup(ans)

		// if something went wrong
		if err != nil {
			// stop listening
			return err
		}

		// add it to the map
		switch ans.(type) {
		case *Option:
			err = core.WriteAnswer(response, q.Name, ans.(*Option).Value)
		case Options:
			err = core.WriteAnswer(response, q.Name, OptionsValues(ans.(Options)))
		default:
			err = core.WriteAnswer(response, q.Name, ans)
		}

		// if something went wrong
		if err != nil {
			return err
		}

	}
	// return the response
	return nil
}
