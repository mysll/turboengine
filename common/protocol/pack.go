package protocol

func PackArgs(args ...interface{}) *Message {
	m := NewMessage(1024)
	ar := NewStoreArchiver(m.Body)
	for _, arg := range args {
		ar.Put(arg)
	}

	m.Body = m.Body[:ar.Len()]
	return m
}

func UnPackArgs(msg *Message, args ...interface{}) error {
	ar := NewLoadArchiver(msg.Body)
	for _, arg := range args {
		if err := ar.Get(arg); err != nil {
			return err
		}
	}
	return nil
}
