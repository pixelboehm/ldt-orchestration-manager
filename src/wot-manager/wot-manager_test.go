package wotmanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const input = `{
    "@context": [
        "http://w3.org/ns/td",
        { "saref": "https://w3id.org/saref#" },
        { "htv": "http://www.w3.org/2011/http#" }
    ],
    "@type": [ "Thing", "saref#LightingDevice" ],
    "id": "urn:dev:wot:com:example:servient:lamp",
    "name": "Awesome Lightbulb",
    "securityDefinitions": {
        "basic_sc": {"scheme": "basic", "in":"header"}
    },
    "security": ["basic_sc"],
    "properties": {
        "status" : {
            "@type": "saref#OnOffState",
            "readOnly": true,
            "writeOnly": false,
            "observable": false,
            "type": "string",
            "forms": [{
                "href": "https://mylamp.example.com/status",
                "contentType": "application/json",
                "htv:methodName": "GET",
                "op": "readproperty"
            }]
        }
    },
    "actions": {
        "toggle" : {
            "@type": "saref#ToggleCommand",
            "idempotent": false,
            "safe": false,
            "forms": [{
                "href": "https://mylamp.example.com/toggle",
                "contentType": "application/json",
                "htv:methodName": "POST",
                "op": "invokeaction"
            }]
        }
    },
    "events":{
        "overheating":{
            "data": {"type": "string"},
            "forms": [{
                "href": "https://mylamp.example.com/oh",
                "contentType": "application/json",
                "subprotocol": "longpoll",
                "op": "subscribeevent"
            }]
        }
    }
}`

func Test_fetchWoTDescription(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	wotm := &WoTManager{description_raw: input}
	got, err := wotm.fetchWoTDescription()
	require.NoError(err)

	assert.Equal(got.Context[0], "http://w3.org/ns/td")
}
