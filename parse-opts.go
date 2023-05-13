package failing

type parsedOpt struct {
	message          *Message
	args             Args
	additionalFields map[string]interface{}
	err              error
}

func (s *Service) parseOpts(opts []interface{}) *parsedOpt {
	var (
		opt    = &parsedOpt{}
		msgKey string
	)

	for _, it := range opts {
		switch val := it.(type) {

		case map[string]interface{}:
			if opt.additionalFields != nil {
				s.logger.Errorf("the field additionalFields is already filled with the value %+v", opt.additionalFields)
			}
			opt.additionalFields = val

		case error:
			opt.err = val

		case Args:
			opt.args = val

		case string:
			msgKey = val

		case *Message:
			opt.message = val

		case Message:
			if opt.message != nil {
				s.logger.Errorf("the field Message is already filled with the value %+v", opt.message)
			}
			opt.message = &val

		case nil:

		default:
			s.logger.Errorf("failing: use only allowed types. type = %T value = %s", val, val)
		}
	}

	if msgKey != "" {
		if opt.message == nil {
			var ok bool
			if opt.message, ok = s.messages[msgKey]; !ok {
				s.logger.Errorf("failing: message not found by key [%s]", msgKey)
			}
		} else {
			s.logger.Errorf("failing: simultaneous use of the Message and message Key parameters is not allowed")
		}
	}

	return opt
}
