package data

type Action struct {
	Twitter TwitterAction `json:"twitter" toml:"twitter" bson:"twitter" yaml:"twitter"`
	Slack   SlackAction   `json:"slack" toml:"slack" bson:"slack" yaml:"slack"`
}

func NewAction() Action {
	return Action{
		Twitter: NewTwitterAction(),
		Slack:   NewSlackAction(),
	}
}

func (a Action) Add(action Action) Action {
	result := a
	result.Twitter = a.Twitter.Add(action.Twitter)
	result.Slack = a.Slack.Add(action.Slack)
	return result
}

func (a Action) Sub(action Action) Action {
	result := a
	result.Twitter = a.Twitter.Sub(action.Twitter)
	result.Slack = a.Slack.Sub(action.Slack)
	return result
}

func (a Action) IsEmpty() bool {
	return a.Twitter.IsEmpty() && a.Slack.IsEmpty()
}
