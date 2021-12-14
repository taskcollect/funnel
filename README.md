HTTP SPEC
---
POST /v1/lessons

Request:
```jsonc
{
    "username": "some_user123",
    "password": "p@ssw0rd_in_plaintext",
    "start": 1609421400, // timestamp, for example this is jan 1st 2021
    "end": 1640957400 // timestamp, for example this is jan 1st 2022
}
```

Response: 200
```jsonc
[
    // array of lessons
    {
        "name": "10 English IGNITE1A 3EG02", // example name
        "id": 1231234, // internal daymap ID
        "start": 1614308100, // timestamp of when the lesson starts
        "finish": 1614311700, // timestamp of when the lesson finishes
        "attendance": "AbsentApproved", // daymap's attendance data
        "resources": false, // does this have any resources
        "links": [
            {
                // did the teacher post any links? these don't contain the link itself, they need to be loaded from daymap
                "label": " Out of Class Learning: Class Reading - The Road",
                "planId": 240373,
                "eventId": 2454307
            },
            {
                "label": "Out of Class Learning: Comparative Poetry Assignment",
                "planId": 244261,
                "eventId": 2454307
            }
        ]
    },
    { /* another lesson */ },
    { /* another lesson */ },
    { /* another lesson */ }
]
```
---
POST /v1/lessons/plans

Request:
```jsonc
{
    "username": "some_user123",
    "password": "p@ssw0rd_in_plaintext",
    "lesson_id": 1234567
}
```

Response: 200
```jsonc
{
    "notes": [
        {
            // fields can disappear from here if they're not set in daymap
            "title": "Note Title",
            "content": "Blah blah blah, we're doing learning today...",
            "links": {
                // links fished out from content
                "Video You Should Watch": "https://youtube.com/(something)"
            },
            "files": {
                // this is a daymap attachment id, a download url can be constructed using it
                "CoolDocument.docx": 123456
            }
        },
        { /* another note on same lesson */ },
        { /* another note on same lesson */ }
    ],
    "extra_files": {
        // files not belonging to any specific note, just placed out in the void
        "ExtraDocument.pdf": 111222
    }
}
```