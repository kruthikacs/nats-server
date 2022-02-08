package server

import "time"

type UserInput struct {
	Nkey        string                `json:"nkey"`
	Permissions *UserInputPermissions `json:"permissions,omitempty"`
}

type UserInputPermissions struct {
	Publish   *UserInputSubjectPermission  `json:"publish,omitempty"`
	Subscribe *UserInputSubjectPermission  `json:"subscribe,omitempty"`
	Response  *UserInputResponsePermission `json:"allow_responses,omitempty"`
}

type UserInputSubjectPermission struct {
	Allow []string `json:"allow,omitempty"`
	Deny  []string `json:"deny,omitempty"`
}

type UserInputResponsePermission struct {
	MaxMsgs int    `json:"max"`
	Expires string `json:"expires"`
}

func (s *Server) deserializeAddUserRequest(userInput *UserInput) (user *NkeyUser, err error) {
	defer s.HandlePanic(&err)

	var nkeyuser *NkeyUser
	nkeyuser = &NkeyUser{}

	// NKey
	nkeyuser.Nkey = userInput.Nkey

	// Permissions
	if userInput.Permissions != nil {
		nkeyuser.Permissions = &Permissions{}
		if userInput.Permissions.Publish != nil {
			nkeyuser.Permissions.Publish = &SubjectPermission{}
			if userInput.Permissions.Publish.Allow != nil {
				for _, element := range userInput.Permissions.Publish.Allow {
					nkeyuser.Permissions.Publish.Allow = append(nkeyuser.Permissions.Publish.Allow, element)
				}
			}

			if userInput.Permissions.Publish.Deny != nil {
				for _, element := range userInput.Permissions.Publish.Deny {
					nkeyuser.Permissions.Publish.Deny = append(nkeyuser.Permissions.Publish.Deny, element)
				}
			}
		}

		if userInput.Permissions.Subscribe != nil {
			nkeyuser.Permissions.Subscribe = &SubjectPermission{}
			if userInput.Permissions.Subscribe.Allow != nil {
				for _, element := range userInput.Permissions.Subscribe.Allow {
					nkeyuser.Permissions.Subscribe.Allow = append(nkeyuser.Permissions.Subscribe.Allow, element)
				}
			}

			if userInput.Permissions.Subscribe.Deny != nil {
				for _, element := range userInput.Permissions.Subscribe.Deny {
					nkeyuser.Permissions.Subscribe.Deny = append(nkeyuser.Permissions.Subscribe.Deny, element)
				}
			}
		}

		if userInput.Permissions.Response != nil {
			nkeyuser.Permissions.Response = &ResponsePermission{}
			nkeyuser.Permissions.Response.MaxMsgs = userInput.Permissions.Response.MaxMsgs
			nkeyuser.Permissions.Response.Expires, err = time.ParseDuration(userInput.Permissions.Response.Expires)
			if err != nil {
				return nil, err
			}
		}
	}

	return nkeyuser, err
}
