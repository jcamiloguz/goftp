package model

type Payload struct {
	Channels []Channel `json:"channels"`
}

type Channel struct {
	Id          int16    `json:"id"`
	Subscribers []Client `json:"subscribers"`
	Files       []File   `json:"files"`
}

type File struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Size int    `json:"size"`
}

type Client struct {
	Id string `json:"id"`
}

func (p *Payload) AddSubscriber(channelId int16, clientId string) {
	for i, channel := range p.Channels {
		if channel.Id == channelId {
			p.Channels[i].Subscribers = append(channel.Subscribers, Client{Id: clientId})

		}
	}
}

func (p *Payload) AddFile(channelId int16, file File) {
	for i, channel := range p.Channels {
		if channel.Id == channelId {
			p.Channels[i].Files = append(channel.Files, file)
		}
	}
}

func (p *Payload) RemoveSubscriber(channelId int16, clientId string) {
	for channelIndex, channel := range p.Channels {
		if channel.Id == channelId {
			for subIndex, subscriber := range channel.Subscribers {
				if subscriber.Id == clientId {
					p.Channels[channelIndex].Subscribers = append(channel.Subscribers[:subIndex], channel.Subscribers[subIndex+1:]...)
				}
			}
		}
	}
}
