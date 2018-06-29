{
  "persistent_menu":[
    {
      "locale":"default",
      "composer_input_disabled": true,
      "call_to_actions":[
        {
              "title":"Previous matches",
              "type":"postback",
              "payload":"GET_PREVIOUS_MATCHES"
        },
        {
          "title":"Preferences",
          "type":"nested",
          "call_to_actions":[
             {
              "title":"Change location preference",
              "type":"postback",
              "payload":"CHANGE_LOCATION_PAYLOAD"
            },
            {
              "title":"Subscribe",
              "type":"postback",
              "payload":"SUBSCRIBE_PAYLOAD"
            },
            {
              "title":"Unsubscribe",
              "type":"postback",
              "payload":"UNSUBSCRIBE_PAYLOAD"
            }
          ]
        },
       {
          "title":"Cheat Codes",
          "type":"nested",
          "call_to_actions":[
            {
              "title":"⬇︎︎⬆︎",
              "type":"postback",
              "payload":"TRIGGER_AVAILABILITY"
            },
            {
              "title":"⬅︎⬆︎",
              "type":"postback",
              "payload":"TRIGGER_MATCH"
            }
          ]
        }
      ]
    }
  ]
}