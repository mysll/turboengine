package protocol

func PackArgs(args ...any) *Message {
	ar := NewAutoExtendArchive(128)
	for _, arg := range args {
		err := ar.Put(arg)
		if err != nil {
			return nil
		}
	}

	return ar.Message()
}

func UnPackArgs(msg *Message, args ...any) error {
	ar := NewLoadArchive(msg.Body)
	for _, arg := range args {
		if err := ar.Get(arg); err != nil {
			return err
		}
	}
	return nil
}
