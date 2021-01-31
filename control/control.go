package control

import (
	"github.com/Mahamed-Belkheir/sunduq"
	"github.com/Mahamed-Belkheir/sunduq/event"
	"github.com/Mahamed-Belkheir/sunduq/storage"
)

//Messages controls the message flow
type Messages struct {
	events   *event.Manager
	store    *storage.Storage
	rootUser string
}

//Authenticate attempts to authenticates the user
func (m Messages) Authenticate(msg sunduq.Message) *sunduq.Message {
	var result sunduq.Message
	username := msg.Key
	password := string(msg.Value)
	user, err := m.store.
		Query(sunduq.Get).
		User(m.rootUser).
		Table(storage.SystemUsersTable).
		Key(username).
		Exec()
	if err != nil {
		result = sunduq.NewResult(msg.ID, true, sunduq.String, []byte(err.Error()))
		return &result
	}
	if password != string(user.Data) {
		result = sunduq.NewResult(msg.ID, true, sunduq.String, []byte("incorrect credentials"))
		return &result
	}
	return nil
}

func (m Messages) recieveHandler(env sunduq.Envelope) {
	msg := env.Message
	var response sunduq.Message

	data, err := m.store.
		Query(msg.Type).
		Key(msg.Key).
		Table(msg.Table).
		Value(storage.Value{
			Type: msg.ValueType,
			Data: msg.Value,
		}).User(env.User).Exec()
	if err != nil {
		response = sunduq.NewResult(
			msg.ID,
			true,
			sunduq.String,
			[]byte(err.Error()),
		)
	} else {
		response = sunduq.NewResult(
			msg.ID,
			false,
			data.Type,
			data.Data,
		)
	}
	m.events.Send(sunduq.NewEnvelope(env.ID, response, env.User))
}
