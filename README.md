# doapi
API Gateway for D05h15h4 University

[![CircleCI](https://circleci.com/gh/mikoim/doapi/tree/master.svg?style=svg)](https://circleci.com/gh/mikoim/doapi/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mikoim/doapi)](https://goreportcard.com/report/github.com/mikoim/doapi)

## Prerequisites

* Go : tested on `go` version 1.7.1
* Redis : `redis` version 1.0.0 or higher is required

## Install

```bash
go get github.com/mikoim/doapi
```

## Usage

### Retrieve cancelled classes on DUET

* Parameter
  * location
    * 1 : Imadegawa Campus
    * 2 : Kyotanabe Campus
    * 3 : Graduate School

```bash
doapi &
curl http://localhost:8080/v1/duet/cancelled_classes?location=1 | jq .
```

```json
[
  {
    "location": "今出川",
    "classes": [
      {
        "period": 1,
        "name": "（´・ω・） ｶﾜｲｿｽ 入門",
        "instructor": "河合 素数",
        "reason": "病気"
      },
      {
        "period": 1,
        "name": "近代山田（´・ω・） ｽ史",
        "instructor": "河合 一郎",
        "reason": "学会"
      }
    ],
    "date": "2016-10-10T00:00:00+09:00",
    "updatedAt": "2016-10-10T13:00:00+09:00"
  },
  {
    "location": "今出川",
    "classes": [
      {
        "period": 1,
        "name": "（´・ω・） ｶﾜｲｿｽ 応用論",
        "instructor": "河合 次郎",
        "reason": "学会"
      }
    ],
    "date": "2016-10-11T00:00:00+09:00",
    "updatedAt": "2016-10-10T13:00:00+09:00"
  },
  {
    "location": "今出川",
    "classes": [
      {
        "period": 4,
        "name": "（´・ω・） ｽﾓﾁ 概論",
        "instructor": "河合 素数子",
        "reason": ""
      }
    ],
    "date": "2016-10-12T00:00:00+09:00",
    "updatedAt": "2016-10-10T13:00:00+09:00"
  },
  {
    "location": "今出川",
    "classes": [
      {
        "period": 1,
        "name": "Logical （´・ω・） suning",
        "instructor": "Sosu Kawai",
        "reason": "学会"
      },
      {
        "period": 6,
        "name": "山田ウイルスと（´・ω・） ｽ",
        "instructor": "河合 三郎",
        "reason": "（´-ω-） ｽﾔｧ…"
      }
    ],
    "date": "2016-10-13T00:00:00+09:00",
    "updatedAt": "2016-10-10T13:00:00+09:00"
  },
  {
    "location": "今出川",
    "classes": [
      {
        "period": 2,
        "name": "（´・ω・） ｶﾜｲｿｽ 権-I",
        "instructor": "河合 四郎",
        "reason": "学会"
      },
      {
        "period": 4,
        "name": "アカデミック・（´・ω・） ｽｷﾙ-V",
        "instructor": "河合 五郎",
        "reason": "ポンポンペイン（；´・ω・） ｽ"
      }
    ],
    "date": "2016-10-14T00:00:00+09:00",
    "updatedAt": "2016-10-10T13:00:00+09:00"
  }
]
```
