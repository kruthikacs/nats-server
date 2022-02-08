package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
)

func (s *Server) AddUser(userInput []byte) (errorPlaceHolder error) {
	defer s.HandlePanic(&errorPlaceHolder)

	// Deserialize to struct and convert to internal NKeyUser
	s.mu.Lock()
	defer s.mu.Unlock()

	var userinput *UserInput
	errorPlaceHolder = json.Unmarshal(userInput, &userinput)
	if errorPlaceHolder != nil {
		return errorPlaceHolder
	}

	if userinput.Nkey == "" {
		errorPlaceHolder = errors.New("nkey is mandatory in the user input")
		return errorPlaceHolder
		//panic("nkey is mandatory in the user input")
	}

	// Deserialize
	var nkeyuser, e = s.deserializeAddUserRequest(userinput)
	if e != nil {
		errorPlaceHolder = e
		return errorPlaceHolder
	}

	var tempArray []*NkeyUser
	tempArray = append(tempArray, nkeyuser)
	nkeyusers, _ := s.buildNkeysAndUsersFromOptions(tempArray, []*User{})
	if s.nkeys == nil {
		s.nkeys = map[string]*NkeyUser{}
	}
	s.nkeys[userinput.Nkey] = nkeyusers[userinput.Nkey]

	value, err := json.Marshal(s.nkeys[userinput.Nkey])
	if err != nil {
		return err
	}
	fmt.Printf("Created user with %s. Size of server nkeys : %d\n", string(value), len(s.nkeys))

	return errorPlaceHolder
}

func (s *Server) DeleteUser(nkey string) (errorPlaceHolder error) {
	defer s.HandlePanic(&errorPlaceHolder)
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.nkeys, nkey)
	fmt.Printf("Deleted user with Nkey : %s. Size of server nkeys : %d\n", nkey, len(s.nkeys))
	return errorPlaceHolder
}

func (s *Server) HandlePanic(errPlaceHolder *error) {
	if panicInformation := recover(); panicInformation != nil {
		fmt.Printf("Recovering.. \n%v, %s", panicInformation, string(debug.Stack()))

		var strBuilder strings.Builder
		strBuilder.WriteString(fmt.Sprintf("%v\n", panicInformation))
		strBuilder.WriteString(string(debug.Stack()))
		*errPlaceHolder = errors.New(strBuilder.String())
	}
}
