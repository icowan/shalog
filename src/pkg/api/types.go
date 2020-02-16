package api

import "encoding/xml"

type memberValue struct {
	Text    string `xml:",chardata"`
	Boolean string `xml:"boolean,omitempty"`
	String  string `xml:"string,omitempty"`
	//DateTimeIso8601 string `xml:"dateTime.iso8601,omitempty"`
	Array interface{} `xml:"array,omitempty"`
	//Int             string `xml:"int,omitempty"`
}

type member struct {
	Text  string      `xml:",chardata"`
	Name  string      `xml:"name,omitempty"`
	Value memberValue `xml:"value,omitempty"`
	//Array array       `xml:"array,omitempty"`
}

type valStruct struct {
	Text   string   `xml:",chardata"`
	Member []member `xml:"member,omitempty"`
	//Array  array    `xml:"array,omitempty"`
}

type dataValue struct {
	Text   string    `xml:",chardata"`
	String string    `xml:"string,omitempty"`
	Struct valStruct `xml:"struct,omitempty"`
}

type data struct {
	Text  string      `xml:",chardata"`
	Value interface{} `xml:"value,omitempty"`
}

type array struct {
	Text string `xml:",chardata"`
	Data data   `xml:"data,omitempty"`
}
type value struct {
	Text  string `xml:",chardata"`
	Array array  `xml:"array,omitempty"`
}

type param struct {
	Text  string `xml:",chardata"`
	Value value  `xml:"value,omitempty"`

	//DateTimeIso8601 string `xml:"dateTime.iso8601"`
	//String          string `xml:"string"`
}

type params struct {
	Text  string `xml:",chardata"`
	Param param  `xml:"param,omitempty"`
}

type getUsersBlogsResponse struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  params   `xml:"params,omitempty"`
}

type getCategoriesResponse struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  struct {
		Text  string `xml:",chardata"`
		Param struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text  string `xml:",chardata"`
				Array struct {
					Text string `xml:",chardata"`
					Data struct {
						Text  string      `xml:",chardata"`
						Value []dataValue `xml:"value"`
					} `xml:"data"`
				} `xml:"array"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

type getPostResponse struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  struct {
		Text  string `xml:",chardata"`
		Param struct {
			Text  string    `xml:",chardata"`
			Value dataValue `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

type postRequest struct {
	XMLName    xml.Name `xml:"methodCall"`
	Text       string   `xml:",chardata"`
	MethodName string   `xml:"methodName"`
	Params     struct {
		Text  string `xml:",chardata"`
		Param []struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text   string `xml:",chardata"`
				String string `xml:"string,omitempty"`
				Struct struct {
					Text   string `xml:",chardata"`
					Member []struct {
						Text  string `xml:",chardata"`
						Name  string `xml:"name"`
						Value struct {
							Text    string `xml:",chardata"`
							String  string `xml:"string"`
							Boolean string `xml:"boolean"`
							Base64  string `xml:"base64"`
							Array   struct {
								Text string `xml:",chardata"`
								Data []struct {
									Text  string `xml:",chardata"`
									Value struct {
										Text   string `xml:",chardata"`
										String string `xml:"string"`
									} `xml:"value"`
								} `xml:"data"`
							} `xml:"array"`
							DateTimeIso8601 string `xml:"dateTime.iso8601"`
						} `xml:"value"`
					} `xml:"member"`
				} `xml:"struct"`
				Boolean string `xml:"boolean"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

type mediaResponse struct {
	XMLName    xml.Name `xml:"methodCall"`
	Text       string   `xml:",chardata"`
	MethodName string   `xml:"methodName"`
	Params     struct {
		Text  string `xml:",chardata"`
		Param []struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text   string `xml:",chardata"`
				String string `xml:"string"`
				Struct struct {
					Text   string `xml:",chardata"`
					Member []struct {
						Text  string `xml:",chardata"`
						Name  string `xml:"name"`
						Value struct {
							Text    string `xml:",chardata"`
							Boolean string `xml:"boolean"`
							Base64  string `xml:"base64"`
							String  string `xml:"string"`
						} `xml:"value"`
					} `xml:"member"`
				} `xml:"struct"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}

type newPostResponse struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Params  struct {
		Text  string `xml:",chardata"`
		Param struct {
			Text  string `xml:",chardata"`
			Value struct {
				Text   string `xml:",chardata"`
				String string `xml:"string"`
			} `xml:"value"`
		} `xml:"param"`
	} `xml:"params"`
}
