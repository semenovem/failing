package failing

type parsedOpt struct {
	message          *Message
	args             Args
	additionalFields map[string]interface{}
}

func (s *Service) parseOpts(opts []interface{}) *parsedOpt {
	var (
		opt    = &parsedOpt{}
		msgKey string
		err    error
	)

	for _, it := range opts {
		switch val := it.(type) {

		case map[string]interface{}:
			if opt.additionalFields != nil {
				s.logger.Errorf("the field additionalFields is already filled with the value %+v", opt.additionalFields)
			}
			opt.additionalFields = val

		case error:
			if err != nil {
				s.logger.Errorf("the err is already filled with the value %+v", err)
			}
			err = val

		case Args:
			if opt.args != nil {
				s.logger.Errorf("the field Args is already filled with the value %+v", opt.args)
			}
			opt.args = val

		case string:
			msgKey = val

		case *Message:
			if opt.message != nil {
				s.logger.Errorf("the field Message is already filled with the value %+v", opt.message)
			}
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

	if s.isDev && err != nil {
		if opt.additionalFields == nil {
			opt.additionalFields = map[string]interface{}{}
		}
		opt.additionalFields[errorFieldName] = err.Error()
	}

	return opt
}
